package service

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
	"github.com/ncw/swift"
	"io"
	"strconv"
)

// ------------------------------------------------------------
// Service implementation

func NewSwiftService(endpoint string) Service {
	if endpoint == "" {
		endpoint = envStAuth.Get()
	}
	user := envStUser.Get()
	key := envStKey.Get()
	return &swiftService{user: user, key: key, authUrl: endpoint}
}

func (s *swiftService) Type() ServiceType {
	return Swift
}

func (s *swiftService) Get(container string, key string) (int64, io.ReadCloser, error) {
	cnx := s.connection()
	file, headers, err := cnx.ObjectOpen(container, key, false, nil)
	if err != nil {
		defer quietly.Close(file)
		return -1, nil, err
	}

	var length int64 = -1
	if contentLength := headers["Content-Length"]; contentLength != "" {
		length, err = strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return -1, nil, err
		}
	}

	return length, file, nil
}

func (s *swiftService) Each(container string, prefix string, do HandleMetadata) (int, error) {
	return s.objectsIn(container, prefix).forEach(func(o swift.Object) error {
		return do(o.Name, o.Bytes)
	})
}

func (s *swiftService) GetEach(container string, prefix string, do HandleObject) (int, error) {
	return s.objectsIn(container, prefix).forEach(func(o swift.Object) error {
		size, body, err := s.Get(container, o.Name)
		return do(o.Name, size, body, err)
	})
}

// ------------------------------------------------------------
// Unexported implementation

const (
	defaultRetries = 3
)

type swiftService struct {
	user    string
	key     string
	authUrl string

	cnx *swift.Connection
}

func (s *swiftService) String() string {
	return fmt.Sprintf("%v (%#v)", Swift, s.authUrl)
}

func (s *swiftService) connection() *swift.Connection {
	if s.cnx == nil {
		s.cnx = &swift.Connection{
			UserName: s.user,
			ApiKey:   s.key,
			AuthUrl:  s.authUrl,
			Retries:  defaultRetries,
		}
	}
	return s.cnx
}

func (s *swiftService) objectsIn(container, prefix string) *swiftObjectIterator {
	return &swiftObjectIterator{cnx: s.connection(), container: container, prefix: prefix}
}

// ------------------------------------------------------------
// Helper types

type swiftObjectIterator struct {
	cnx *swift.Connection
	container string
	prefix    string
}

func (it *swiftObjectIterator) forEach(do func(o swift.Object) error) (int, error) {
	var opts *swift.ObjectsOpts
	if it.prefix != "" {
		opts = &swift.ObjectsOpts{Prefix: it.prefix}
	}

	objects, err := it.cnx.Objects(it.container, opts)
	if err != nil {
		return -1, err
	}
	for i, o := range objects {
		err = do(o)
		if err != nil {
			return i, err
		}
	}
	return len(objects), nil
}
package service

import (
	"fmt"
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
	cnx := s.Connection()
	file, headers, err := cnx.ObjectOpen(container, key, false, nil)
	if err != nil {
		defer CloseQuietly(file)
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

func (s *swiftService) ContentLength(container string, key string) (int64, error) {
	cnx := s.Connection()
	info, _, err := cnx.Object(container, key)
	if err != nil {
		return -1, err
	}
	return info.Bytes, nil
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

	connection *swift.Connection
}

func (s *swiftService) String() string {
	return fmt.Sprintf("%v (%#v)", Swift, s.authUrl)
}

func (s *swiftService) Connection() *swift.Connection {
	if s.connection == nil {
		s.connection = &swift.Connection{
			UserName: s.user,
			ApiKey:   s.key,
			AuthUrl:  s.authUrl,
			Retries:  defaultRetries,
		}
	}
	return s.connection
}

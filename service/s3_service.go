package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dmolesUC3/mrt-bits/internal/pointers"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
	"github.com/dmolesUC3/mrt-bits/internal/streams"
	"io"
	"regexp"
)

// ------------------------------------------------------------
// Service implementation

func NewS3Service(region, endpoint string) Service {
	if endpoint == "" {
		endpoint = defaultAwsEndpoint
	}
	if region == "" {
		region = envAwsRegion.Get()
	}
	if region == "" {
		region = regionFromEndpoint(endpoint)
	}
	key := envAwsAccessKeyId.Get()
	secret := envAwsSecretAccessKey.Get()
	sessionToken := envAwsSessionToken.Get()
	if key != "" {
		return &s3Service{
			region:      region,
			endpoint:    endpoint,
			credentials: credentials.NewStaticCredentials(key, secret, sessionToken),
		}
	}
	return &s3Service{region: region, endpoint: endpoint}
}

func (s *s3Service) Type() ServiceType {
	return S3
}

func (s *s3Service) Get(container string, key string) (int64, io.ReadCloser, error) {
	s3svc, err := s.s3()
	if err != nil {
		return -1, nil, err
	}
	input := &s3.GetObjectInput{Bucket: &container, Key: &key}
	output, err := s3svc.GetObject(input)
	if err != nil {
		defer quietly.Close(output.Body)
		return -1, nil, err
	}
	return *output.ContentLength, output.Body, nil
}

func (s *s3Service) EachMetadata(container string, prefix string, do HandleMetadata) (int, error) {
	objects, err := s.objectsIn(container, prefix)
	if err != nil {
		return -1, err
	}
	return objects.forEach(func(o *s3.Object) error {
		key := pointers.ToString(o.Key)
		contentLength := pointers.ToInt64(o.Size)
		return do(key, contentLength)
	})
}

func (s *s3Service) EachObject(container string, prefix string, do HandleObject) (int, error) {
	objects, err := s.objectsIn(container, prefix)
	if err != nil {
		return -1, err
	}
	return objects.forEach(func(o *s3.Object) error {
		key := pointers.ToString(o.Key)
		size := pointers.ToInt64(o.Size)
		if size == 0 {
			return do(key, size, streams.EmptyReader(), nil)
		}
		size, body, err := s.Get(container, key)
		return do(key, size, body, err)
	})
}

// ------------------------------------------------------------
// Unexported implementation

type s3Service struct {
	region      string
	endpoint    string
	credentials *credentials.Credentials

	awsSession   *session.Session
	s3Svc        *s3.S3
	s3Downloader *s3manager.Downloader
}

func (s *s3Service) String() string {
	return fmt.Sprintf("%v (%#v, %#v)", S3, s.region, s.endpoint)
}

func (s *s3Service) session() (*session.Session, error) {
	if s.awsSession == nil {
		awsSession, err := s.newSession()
		if err != nil {
			return nil, err
		}
		s.awsSession = awsSession
	}
	return s.awsSession, nil
}

func (s *s3Service) downloader() (*s3manager.Downloader, error) {
	if s.s3Downloader == nil {
		awsSession, err := s.session()
		if err != nil {
			return nil, err
		}
		s.s3Downloader = s3manager.NewDownloader(awsSession)
	}
	return s.s3Downloader, nil
}

func (s *s3Service) s3() (*s3.S3, error) {
	if s.s3Svc == nil {
		awsSession, err := s.session()
		if err != nil {
			return nil, err
		}
		s.s3Svc = s3.New(awsSession)
	}
	return s.s3Svc, nil
}

func (s *s3Service) newSession() (awsSession *session.Session, err error) {
	forcePathStyle := true
	credentialsChainVerboseErrors := false
	s3Config := aws.Config{
		Endpoint:                      &s.endpoint,
		Region:                        &s.region,
		S3ForcePathStyle:              &forcePathStyle,
		Credentials:                   s.credentials,
		CredentialsChainVerboseErrors: &credentialsChainVerboseErrors,
	}
	s3Opts := session.Options{
		Config:            s3Config,
		SharedConfigState: session.SharedConfigEnable,
	}
	return session.NewSessionWithOptions(s3Opts)
}

func (s *s3Service) objectsIn(container, prefix string) (*s3ObjectIterator, error) {
	s3svc, err := s.s3()
	if err != nil {
		return nil, err
	}
	return &s3ObjectIterator{s3svc: s3svc, container: container, prefix: prefix}, nil
}

// ------------------------------------------------------------
// Helper functions

const (
	defaultAwsEndpoint = "" // SDK will figure it out from region
	defaultAwsRegion   = "us-west-2"
	awsRegionRegexpStr = "https?://s3-([^.]+)\\.amazonaws\\.com"
)

var awsRegionRegexp = regexp.MustCompile(awsRegionRegexpStr)

func regionFromEndpoint(endpoint string) string {
	matches := awsRegionRegexp.FindStringSubmatch(endpoint)
	if len(matches) == 2 {
		regionStr := matches[1]
		return regionStr
	}
	return defaultAwsRegion
}

// ------------------------------------------------------------
// Helper types

type s3ObjectIterator struct {
	s3svc     *s3.S3
	container string
	prefix    string
}

func (it *s3ObjectIterator) forEach(do func(o *s3.Object) error) (int, error) {
	var errInner error
	var count int

	pageHandler := func(output *s3.ListObjectsV2Output, b bool) bool {
		for _, o := range output.Contents {
			errInner = do(o)
			if errInner != nil {
				return false
			}
			count += 1
		}
		return true
	}

	input := &s3.ListObjectsV2Input{Bucket: &it.container, Prefix: pointers.FromString(it.prefix)}
	errOuter := it.s3svc.ListObjectsV2Pages(input, pageHandler)
	if errInner != nil {
		return count, errInner
	}
	return count, errOuter
}


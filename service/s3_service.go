package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
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

func (s *s3Service) ContentLength(container string, key string) (int64, error) {
	s3svc, err := s.s3()
	if err != nil {
		return -1, err
	}
	input := &s3.HeadObjectInput{Bucket: &container, Key: &key}
	output, err := s3svc.HeadObject(input)
	if err != nil {
		return -1, err
	}
	return *output.ContentLength, nil
}

func (s *s3Service) Each(container string, prefix string, do func(string) error) (int, error) {
	panic("implement me")
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

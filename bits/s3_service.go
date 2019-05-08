package bits

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"regexp"
)

// ------------------------------------------------------------
// Service implementation

func NewS3Service(region string, endpoint string) Service {
	return &s3Service{region: region, endpoint: endpoint}
}

func NewS3ServiceForRegion(region string) Service {
	return NewS3Service(region, defaultAwsEndpoint)
}

func NewS3ServiceForEndpoint(endpoint string) Service {
	region := regionFromEndpoint(endpoint)
	return NewS3Service(region, endpoint)
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

// ------------------------------------------------------------
// Unexported implementation

type s3Service struct {
	region   string
	endpoint string

	awsSession   *session.Session
	s3Svc        *s3.S3
	s3Downloader *s3manager.Downloader
}

func (s *s3Service) session() (*session.Session, error) {
	if s.awsSession == nil {
		awsSession, err := validS3Session(s.endpoint, s.region)
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

// ------------------------------------------------------------
// Helper functions

const (
	defaultAwsEndpoint = ""
	defaultAwsRegion   = "us-west-2"
	awsRegionRegexpStr = "https?://s3-([^.]+)\\.amazonaws\\.com"
)

var awsRegionRegexp = regexp.MustCompile(awsRegionRegexpStr)

func validS3Session(endpoint string, region string) (awsSession *session.Session, err error) {
	forcePathStyle := true
	credentialsChainVerboseErrors := false
	s3Config := aws.Config{
		Endpoint:                      &endpoint,
		Region:                        &region,
		S3ForcePathStyle:              &forcePathStyle,
		CredentialsChainVerboseErrors: &credentialsChainVerboseErrors,
	}
	s3Opts := session.Options{
		Config:            s3Config,
		SharedConfigState: session.SharedConfigEnable,
	}
	return session.NewSessionWithOptions(s3Opts)
}

func regionFromEndpoint(endpoint string) string {
	matches := awsRegionRegexp.FindStringSubmatch(endpoint)
	if len(matches) == 2 {
		regionStr := matches[1]
		return regionStr
	}
	return defaultAwsRegion
}

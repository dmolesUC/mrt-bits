package legacy

import (
	"net/url"
)

// ------------------------------------------------------------
// S3Service

func S3Service(name string, endpoint *url.URL, bucket string) Service {
	service := service{name: name, endpoint: endpoint}
	return &s3Service{service: service, bucket: bucket}
}

// ------------------------------------------------------------
// Unexported symbols

type s3Service struct {
	service
	bucket string
}

func (s *s3Service) Type() ServiceType {
	return s3
}

func (s *s3Service) ContainerFor(ark string) string {
	return s.bucket
}


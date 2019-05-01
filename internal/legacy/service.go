package legacy

import (
	"net/url"
)

// ------------------------------------------------------------
// Service

type Service interface {
	Name() string
	Endpoint() *url.URL
	Type() ServiceType
	ContainerFor(ark string) string
}

// ------------------------------------------------------------
// Unexported symbols

// ------------------------------
// service

type service struct {
	name        string
	endpoint    *url.URL
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Endpoint() *url.URL {
	return s.endpoint
}

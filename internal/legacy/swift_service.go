package legacy

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ------------------------------------------------------------
// SwiftService

const (
	ST_AUTH = "ST_AUTH"
	ST_USER = "ST_USER"
	ST_KEY  = "ST_KEY"
)

func SwiftService(name string, endpoint *url.URL, user string, key string, containerBase string) Service {
	return &swiftService{
		service: service{
			name:     name,
			endpoint: endpoint,
		},
		user:          user,
		key:           key,
		containerBase: containerBase,
	}
}

func SwiftServiceFromEnv(name string, containerBase string) (Service, error) {
	authUrl := os.Getenv(ST_AUTH)
	if authUrl == "" {
		return nil, fmt.Errorf("$%v not set", ST_AUTH)
	}
	endpoint, err := url.Parse(authUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing $%v %#v: %v", ST_AUTH, authUrl, err.Error())
	}
	user := os.Getenv(ST_USER)
	if user == "" {
		return nil, fmt.Errorf("$%v not set", ST_USER)
	}
	key := os.Getenv(ST_KEY)
	if key == "" {
		return nil, fmt.Errorf("$%v not set", ST_KEY)
	}
	return SwiftService(name, endpoint, user, key, containerBase), nil
}

// ------------------------------------------------------------
// Unexported symbols

type swiftService struct {
	service
	containerBase string
	user          string
	key           string
}

func (s *swiftService) Type() ServiceType {
	return swift
}

func (s *swiftService) ContainerFor(ark string) string {
	if !strings.HasSuffix(s.containerBase, "__") {
		return s.containerBase
	}
	hash := md5.New()
	hash.Write(([]byte)(ark))
	resultStr := fmt.Sprintf("%x", hash.Sum(nil))
	return s.containerBase + resultStr[0:3]
}

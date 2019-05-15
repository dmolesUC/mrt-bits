package service

import (
	"fmt"
	"io"
)

type Service interface {
	fmt.Stringer
	Type() ServiceType
	Get(container string, key string) (contentLength int64, body io.ReadCloser, err error)
	GetEach(container string, prefix string, do HandleObject) (int, error)
	Each(container string, prefix string, do HandleMetadata) (int, error)
}

type HandleMetadata func(key string, contentLength int64) error

type HandleObject func(key string, contentLength int64, body io.ReadCloser, err error) error
package service

import (
	"fmt"
	"io"
)

type Service interface {
	fmt.Stringer
	Type() ServiceType
	ContentLength(container string, key string) (int64, error)
	Get(container string, key string) (contentLength int64, body io.ReadCloser, err error)
	GetEach(container string, prefix string, do func(int64, io.ReadCloser, error) error) (int, error)
	Each(container string, prefix string, do func(string) error) (int, error)
}


package service

import (
	"io"
)

type Service interface {
	Type() ServiceType
	Get(container string, key string) (int64, io.ReadCloser, error)
	ContentLength(container string, key string) (int64, error)
}


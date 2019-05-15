package service

import (
	"fmt"
	"io"
)

type Service interface {
	fmt.Stringer
	ObjectIterator
	Type() ServiceType
	Get(container string, key string) (contentLength int64, body io.ReadCloser, err error)
}

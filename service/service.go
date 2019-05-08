package service

import (
	"fmt"
	"io"
)

type Service interface {
	fmt.Stringer
	Type() ServiceType
	Get(container string, key string) (contentLength int64, body io.ReadCloser, err error)
	ContentLength(container string, key string) (int64, error)
}

func CloseQuietly(body io.Closer) func() {
	return func() {
		if body != nil {
			_ = body.Close()
		}
	}
}
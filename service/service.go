package service

import (
	"fmt"
	"io"
	"time"
)

type Service interface {
	fmt.Stringer
	ObjectIterator
	Type() ServiceType
	Get(container string, key string) (contentLength int64, body io.ReadCloser, err error)
	GetSize(container string, key string) (int64, error)
}

type Info interface {
	ContentLength() int64
	LastModified() time.Time
}

type info struct {
	contentLength int64
	lastModified time.Time
}

func (i *info) LastModified() time.Time {
	panic("implement me")
}

func (i *info) ContentLength() int64 {
	return i.contentLength
}


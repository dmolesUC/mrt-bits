package service

import (
	"io"
)

// TODO: introduce explicit Metadata type
type HandleMetadata func(key string, contentLength int64) error

// TODO: introduce explict Object type
type HandleObject func(key string, contentLength int64, body io.ReadCloser, err error) error

type ObjectIterator interface {
	EachMetadata(container string, prefix string, do HandleMetadata) (int, error)
	EachObject(container string, prefix string, do HandleObject) (int, error)
}

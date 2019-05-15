package operations

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/service"
	"io"
)

type ListObjects interface {
	To(out io.Writer) (int, error)
}

func NewListObjects(svc service.Service, container, prefix string) ListObjects {
	return &list{svc: svc, container: container, prefix: prefix}
}

// ------------------------------------------------
// Unexported symbols

type list struct {
	svc       service.Service
	container string
	prefix    string
}

func (l *list) To(out io.Writer) (int, error) {
	// TODO: optionally print size
	return l.svc.EachMetadata(l.container, l.prefix, func(key string, contentLength int64) error {
		_, err := fmt.Fprintln(out, key)
		return err
	})
}



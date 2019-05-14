package operations

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/service"
	"io"
)

type List interface {
	To(out io.Writer) (int, error)
}

func NewList(svc service.Service, container, prefix string) List {
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
	return l.svc.Each(l.container, l.prefix, func(s string) error {
		_, err := fmt.Fprintln(out, s)
		return err
	})
}



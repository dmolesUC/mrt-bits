package operations

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/dmolesUC3/mrt-bits/service"
	"hash"
)

type Info interface {
	Size() (int64, error)
	SHA256() ([]byte, error)
	MD5() ([]byte, error)
}

func NewInfo(svc service.Service, container, key string) Info {
	return &info{svc: svc, container: container, key: key}
}

// ------------------------------------------------
// Unexported symbols

type info struct {
	svc       service.Service
	container string
	key       string
}

func (i *info) Size() (int64, error) {
	return i.svc.GetSize(i.container, i.key)
}

func (i *info) SHA256() ([]byte, error) {
	return i.Hash(sha256.New())
}

func (i *info) MD5() ([]byte, error) {
	return i.Hash(md5.New())
}

func (i *info) Hash(h hash.Hash) ([]byte, error) {
	expectedSize, err := i.Size()
	if err != nil {
		return nil, err
	}

	download := NewDownload(i.svc, i.container, i.key)
	n, err := download.To(h)
	if err != nil {
		return nil, err
	}
	if int64(n) != expectedSize {
		err = fmt.Errorf("expected %d bytes, got %d", expectedSize, n)
	}
	digest := h.Sum(nil)
	return digest, err
}

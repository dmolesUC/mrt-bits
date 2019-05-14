package operations

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
	"github.com/dmolesUC3/mrt-bits/service"
	"io"
	"os"
	"path"
)

type Download interface {
	To(out io.Writer) (int, error)
	ToFile(filename string) (int, error)
	ToRemoteFile() (int, error)
}

func NewDownload(svc service.Service, container, key string) Download {
	return &download{svc: svc, container: container, key: key}
}

// ------------------------------------------------
// Unexported symbols

type download struct {
	svc       service.Service
	container string
	key       string
}

func (d *download) ToFile(filename string) (int, error) {
	// TODO: don't create file till we know we've got something
	//       (and/or quietly delete file)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		defer quietly.Close(file)
		if err != nil {
			return 0, err
		}
		return d.To(file)
	}
	return 0, fmt.Errorf("file %#v already exists", filename)
}

func (d *download) ToRemoteFile() (int, error) {
	return d.ToFile(path.Base(d.key))
}

func (d *download) To(out io.Writer) (int, error) {
	_, body, err := d.svc.Get(d.container, d.key)
	defer quietly.Close(body)
	if err != nil {
		return 0, err
	}
	total := 0
	buffer := make([]byte, bufsize)
	for {
		n, err := body.Read(buffer)
		if n > 0 {
			total += n
			_, err2 := out.Write(buffer[:n])
			if err2 != nil {
				return total, err2
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

package streams

import (
	"io"
)

// ------------------------------------------------------------
// EmptyReader

func EmptyReader() io.ReadCloser {
	return &r
}

type emptyReader int

var r emptyReader

func (*emptyReader) Close() error {
	return nil
}

func (*emptyReader) Read(p []byte) (n int, err error) {
	return 0, nil
}





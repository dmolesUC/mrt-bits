package operations

import (
	"io"
)

type Archive interface {
	To(out io.Writer) (int, error)
}

// ------------------------------------------------
// Unexported symbols

type archive struct {

}
package quietly

import (
	"fmt"
	"io"
)

func Close(body io.Closer) func() {
	return func() {
		if body != nil {
			_ = body.Close()
		}
	}
}

func Fprintf(out io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(out, format, a...)
}
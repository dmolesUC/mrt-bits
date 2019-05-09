package quietly

import (
	"io"
)

func Close(body io.Closer) func() {
	return func() {
		if body != nil {
			_ = body.Close()
		}
	}
}
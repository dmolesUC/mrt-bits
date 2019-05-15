package operations

import (
	"archive/zip"
	"fmt"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
	"github.com/dmolesUC3/mrt-bits/service"
	"io"
	"strings"
)

type Archive interface {
	To(out io.Writer) (int, error)
}

// ------------------------------------------------
// Unexported symbols

const (
	fileHeaderLen       = 30 // + filename + extra
	directoryHeaderLen  = 46 // + filename + extra + comment
	directoryEndLen     = 22 // + comment
	dataDescriptorLen   = 16 // four uint32: descriptor signature, crc32, compressed size, size
	dataDescriptor64Len = 24 // descriptor with 8 byte sizes
	directory64LocLen   = 20 //
	directory64EndLen   = 56 // + extra

	// Limits for non zip64 files.
	uint16max = (1 << 16) - 1
	uint32max = (1 << 32) - 1
)

type zipArchive struct {
	service   service.Service
	container string
	prefix    string
}

func (a *zipArchive) Size() (int64, error) {
	var count int64
	var size int64
	var cdSize int64
	_, err := a.service.Each(a.container, a.prefix, func(key string, contentLength int64) error {
		entryName := key
		nameLen := int64(len(entryName))

		size += fileHeaderLen
		size += nameLen

		// directories don't get data descriptors
		if strings.HasSuffix(entryName, "/") {
			if contentLength > 0 {
				return fmt.Errorf("bad entry %#v: can't zip a plain file that looks like a directory", entryName)
			}
		} else {
			size += contentLength
			isZip64 := contentLength >= uint32max
			if isZip64 {
				size += dataDescriptor64Len
			} else {
				size += dataDescriptorLen
			}
		}

		cdSize += directoryHeaderLen + nameLen

		count++
		return nil
	})

	cdOffset := size

	if count >= uint16max || cdSize >= uint32max || cdOffset >= uint32max {
		size += directory64EndLen + directory64LocLen
	} else {
		size += directoryEndLen
	}

	return size, err
}

func (a *zipArchive) To(out io.Writer) (int, error) {
	w := zip.NewWriter(out)
	defer quietly.Close(w)
	return a.service.GetEach(a.container, a.prefix, func(key string, contentLength int64, body io.ReadCloser, err error) error {
		defer quietly.Close(body)
		entryName := key
		header := &zip.FileHeader{Name: entryName, Method: zip.Store}
		out, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		bytesWritten, err := io.Copy(out, body)
		if err != nil {
			return err
		}
		if bytesWritten != contentLength {
			return fmt.Errorf("error writing %#v; expected to write %d bytes, wrote %d", entryName, contentLength, bytesWritten)
		}
		return nil
	})
}

package operations

import (
	"archive/zip"
	"fmt"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
	"github.com/dmolesUC3/mrt-bits/service"
	"io"
	"math"
	"strings"
)

type Archive interface {
	Size() (size int64, count int, err error)
	To(out io.Writer) (int, error)
}

func NewZipArchive(service service.ObjectIterator, container string, prefix string) Archive {
	return &zipArchive{service: service, container: container, prefix: prefix}
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
)

type zipArchive struct {
	service   service.ObjectIterator
	container string
	prefix    string
}

func (a *zipArchive) Size() (int64, int, error) {
	var count int
	var size int64
	var cdSize int64
	_, err := a.service.EachMetadata(a.container, a.prefix, func(key string, contentLength int64) error {
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
			isZip64 := contentLength >= math.MaxUint32
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

	if count >= math.MaxUint16 || cdSize >= math.MaxUint32 || cdOffset >= math.MaxUint32 {
		size += directory64EndLen + directory64LocLen
	} else {
		size += directoryEndLen
	}

	size += cdSize

	return size, count, err
}

func (a *zipArchive) To(out io.Writer) (int, error) {
	w := zip.NewWriter(out)
	count, err := a.service.EachObject(a.container, a.prefix, func(key string, contentLength int64, body io.ReadCloser, err error) error {
		defer quietly.Close(body)
		entryName := key
		// TODO: more metadata -- at least modification time
		header := &zip.FileHeader{Name: entryName, Method: zip.Store}
		entryWriter, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		bytesWritten, err := io.Copy(entryWriter, body)
		if err != nil {
			return err
		}
		if bytesWritten != contentLength {
			return fmt.Errorf("error writing %#v; expected to write %d bytes, wrote %d", entryName, contentLength, bytesWritten)
		}
		return nil
	})
	err2 := w.Close()
	if err2 == nil {
		return count, err
	}
	if err == nil {
		return count, err2
	}
	return count, fmt.Errorf("error creating archive: %v. In addition, an error occurred closing the zip stream: %v", err.Error(), err2.Error())
}

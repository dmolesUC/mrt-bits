package operations

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
	"github.com/dmolesUC3/mrt-bits/service"
	. "gopkg.in/check.v1"
	"io"
	"io/ioutil"
	"math"
	"math/bits"
	"math/rand"
)

// ------------------------------------------------
// Fixture

type ArchiveSuite struct{}

var _ = Suite(&ArchiveSuite{})

type objMock struct {
	index         int
	key           string
	contentLength int64
	rnd           *rand.Rand
}

func (o *objMock) body() io.ReadCloser {
	if o.rnd == nil {
		o.rnd = rand.New(rand.NewSource(int64(o.index)))
	}
	lr := &io.LimitedReader{R: o.rnd, N: o.contentLength}
	return ioutil.NopCloser(lr)
}


func md5sum(in io.ReadCloser, expectedSize int64) string {
	defer quietly.Close(in)
	hash := md5.New()
	n, err := io.Copy(hash, in)
	if err != nil {
		panic(err)
	}
	if n != expectedSize {
		panic(fmt.Errorf("error in hash calculation: expected to hash %d bytes, got %d", expectedSize, n))
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

type objIterMock struct {
	size          int
	key           func(index int) string
	contentLength func(index int) int64
}

func (it *objIterMock) objectAt(index int) *objMock {
	return &objMock{index: index, key: it.key(index), contentLength: it.contentLength(index)}
}

func (it *objIterMock) EachMetadata(container string, prefix string, do service.HandleMetadata) (int, error) {
	for i := 0; i < it.size; i++ {
		o := it.objectAt(i)
		err := do(o.key, o.contentLength)
		if err != nil {
			return i, err
		}
	}
	return it.size, nil
}

func (it *objIterMock) EachObject(container string, prefix string, do service.HandleObject) (int, error) {
	for i := 0; i < it.size; i++ {
		o := it.objectAt(i)
		err := do(o.key, o.contentLength, o.body(), nil)
		if err != nil {
			return i, err
		}
	}
	return it.size, nil
}

func defaultKey(index int) string {
	return fmt.Sprintf("file-%d.bin", index)
}

func newFixedIterator(size int, contentLength int64) *objIterMock {
	return &objIterMock{
		size: size,
		key:  defaultKey,
		contentLength: func(index int) int64 {
			return contentLength
		},
	}
}

func newExponentialIterator(minContentLengthMax uint64) *objIterMock {
	if minContentLengthMax > math.MaxInt64 {
		panic(fmt.Errorf("invalid content-length: %d", minContentLengthMax))
	}
	return &objIterMock{
		size: bits.Len64(minContentLengthMax),
		key:  defaultKey,
		contentLength: func(index int) int64 {
			return 1 << uint64(index)
		},
	}
}

// ------------------------------------------------
// Test

func (s *ArchiveSuite) TestArchiveTo(c *C) {
	it := newFixedIterator(10, 1024)
	archive := NewZipArchive(it, "", "") // TODO: validate container and prefix

	out := new(bytes.Buffer)
	count, err := archive.To(out)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, it.size)

	size := out.Len()
	zipdata := make([]byte, size)
	copy(zipdata, out.Bytes())

	in, err := zip.NewReader(bytes.NewReader(zipdata), int64(size))
	entries := in.File
	c.Assert(len(entries), Equals, count)

	for i, f := range entries {
		o := it.objectAt(i)
		c.Assert(f.Name, Equals, o.key)
		expectedSize := uint64(o.contentLength)
		c.Assert(f.CompressedSize64, Equals, expectedSize)
		c.Assert(f.UncompressedSize64, Equals, expectedSize)

		content, err := f.Open()
		c.Assert(err, IsNil)
		md5Expected := md5sum(o.body(), o.contentLength)
		md5Actual := md5sum(content, o.contentLength)
		c.Assert(md5Actual, Equals, md5Expected)
	}
}

func (s *ArchiveSuite) TestSize(c *C) {
	it := newFixedIterator(10, 1024)
	archive := NewZipArchive(it, "", "")
	size, err := archive.Size()
	c.Assert(err, IsNil)
	var expectedSize int64 = 11382 // based on previous test
	c.Assert(size, Equals, expectedSize)
}

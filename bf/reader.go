package bf

import (
	"errors"
	"io"
)

var (
	maxEmptyReads    = 100
	errReaderIsNil   = errors.New("BufferReader: reader is nil")
	errNegativeCount = errors.New("BufferReader: read return negative count")
	errNoProgress    = errors.New("BufferReader: multiple Read calls return no data or error")
	errTooLarge      = errors.New("BufferReader: make byte slice too large")
)

// BufferReader ...
type BufferReader struct {
	buf    []byte
	reader io.Reader
	size   int
	r, w   int
	err    error
}

// NewBufferReader ...
func NewBufferReader(reader io.Reader, size int) *BufferReader {
	return &BufferReader{
		reader: reader,
		size:   size,
		buf:    make([]byte, size),
	}
}

// Reset ...
func (br *BufferReader) Reset() {
	if br.w > br.r {
		copy(br.buf, br.buf[br.r:br.w])
	}

	br.w = br.w - br.r
	br.r = 0
}

// ReadFull ...
func (br *BufferReader) ReadFull(min int) (data []byte, err error) {
	if br.reader == nil {
		return nil, errReaderIsNil
	}

	if min == 0 {
		err = br.err
		br.err = nil
		return make([]byte, 0, 0), err
	}

	if min > (cap(br.buf) - br.r) {
		br.Grow(min)
	}

	for (br.w-br.r) < min && err == nil {
		br.fill()
		err = br.err
	}

	if (br.w - br.r) >= min {
		data = br.buf[br.r : br.r+min]
		br.r = br.r + min
		err = nil
	} else {
		data = br.buf[br.r:br.w]
		br.r = br.w
		err = br.err
		br.err = nil
	}
	return
}

func (br *BufferReader) fill() {
	if br.w >= cap(br.buf) {
		br.Grow(br.w - br.r)
	}

	for i := maxEmptyReads; i > 0; i-- {
		n, err := br.reader.Read(br.buf[br.w:])
		if n < 0 {
			panic(errNegativeCount)
		}
		br.w = br.w + n
		if err != nil {
			br.err = err
			return
		}
		if n > 0 {
			return
		}
	}

	br.err = errNoProgress
}

// Grow ...
func (br *BufferReader) Grow(n int) {
	defer func() {
		if recover() != nil {
			panic(errTooLarge)
		}
	}()

	var buf []byte = nil
	if n > br.size {
		buf = make([]byte, n)
	} else {
		buf = make([]byte, br.size)
	}

	if br.w > br.r {
		copy(buf, br.buf[br.r:br.w])
	}

	br.w = br.w - br.r
	br.r = 0
	br.buf = buf
}

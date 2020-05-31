package bf

import (
	"bufio"
	"io"
)

// BufferWriter ...
type BufferWriter struct {
	*bufio.Writer
}

// NewBufferWriter ...
func NewBufferWriter(wr io.Writer, size int) *BufferWriter {
	return &BufferWriter{
		Writer: bufio.NewWriterSize(wr, size),
	}
}

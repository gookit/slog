package bufwrite

import (
	"bufio"
	"io"
)

// BufIOWriter struct
type BufIOWriter struct {
	*bufio.Writer
}

// Close implements the io.Closer
func (w *BufIOWriter) Close() error {
	return w.Flush()
}

// Sync implements the Syncer
func (w *BufIOWriter) Sync() error {
	return w.Flush()
}

// NewBufIOWriterSize instance with size
func NewBufIOWriterSize(w io.Writer, size int) *BufIOWriter {
	return &BufIOWriter{
		Writer: bufio.NewWriterSize(w, size),
	}
}

// NewBufIOWriter instance
func NewBufIOWriter(w io.Writer) *BufIOWriter {
	return &BufIOWriter{
		Writer: bufio.NewWriterSize(w, defaultBufSize),
	}
}

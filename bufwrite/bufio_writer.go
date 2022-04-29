package bufwrite

import (
	"bufio"
	"io"
)

// BufIOWriter wrap the bufio.Writer, implements the Sync() Close() methods
type BufIOWriter struct {
	bufio.Writer
	// backup the bufio.Writer.wr
	writer io.Writer
}

// NewBufIOWriterSize instance with size
func NewBufIOWriterSize(w io.Writer, size int) *BufIOWriter {
	return &BufIOWriter{
		writer: w,
		Writer: *bufio.NewWriterSize(w, size),
	}
}

// NewBufIOWriter instance
func NewBufIOWriter(w io.Writer) *BufIOWriter {
	return NewBufIOWriterSize(w, defaultBufSize)
}

// Close implements the io.Closer
func (w *BufIOWriter) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}

	// is closer
	if c, ok := w.writer.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// Sync implements the Syncer
func (w *BufIOWriter) Sync() error {
	return w.Flush()
}

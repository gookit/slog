package bufwrite

import (
	"io"
)

const (
	defaultBufSize = 1024 * 8
)

// LineWriter implements buffering for an io.Writer object.
// If an error occurs writing to a LineWriter, no more data will be
// accepted and all subsequent writes, and Flush, will return the error.
// After all data has been written, the client should call the
// Flush method to guarantee all data has been forwarded to
// the underlying io.Writer.
//
// from bufio.Writer.
//
// Change:
//
// always keep write full line. more difference please see Write
type LineWriter struct {
	err error
	buf []byte
	n   int
	wr  io.Writer
}

// NewLineWriterSize returns a new LineWriter whose buffer has at least the specified
// size. If the argument io.Writer is already a LineWriter with large enough
// size, it returns the underlying LineWriter.
func NewLineWriterSize(w io.Writer, size int) *LineWriter {
	// Is it already a LineWriter?
	b, ok := w.(*LineWriter)
	if ok && len(b.buf) >= size {
		return b
	}
	if size <= 0 {
		size = defaultBufSize
	}

	return &LineWriter{
		buf: make([]byte, size),
		wr:  w,
	}
}

// NewLineWriter returns a new LineWriter whose buffer has the default size.
func NewLineWriter(w io.Writer) *LineWriter {
	return NewLineWriterSize(w, defaultBufSize)
}

// Size returns the size of the underlying buffer in bytes.
func (b *LineWriter) Size() int { return len(b.buf) }

// Reset discards any un-flushed buffered data, clears any error, and
// resets b to write its output to w.
func (b *LineWriter) Reset(w io.Writer) {
	b.err = nil
	b.n = 0
	b.wr = w
}

// Close implements the io.Closer
func (b *LineWriter) Close() error {
	if err := b.Flush(); err != nil {
		return err
	}

	// is closer
	if c, ok := b.wr.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// Sync implements the Syncer
func (b *LineWriter) Sync() error {
	return b.Flush()
}

// Flush writes any buffered data to the underlying io.Writer.
//
// TIP: please add lock before call the method.
func (b *LineWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}

// Available returns how many bytes are unused in the buffer.
func (b *LineWriter) Available() int { return len(b.buf) - b.n }

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *LineWriter) Buffered() int { return b.n }

// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *LineWriter) Write(p []byte) (nn int, err error) {
	// 原来的会造成 p 写了一部分到 b.wr, 还有一部分在 b.buf，
	// 如果现在外部工具从 b.wr 收集数据，会收集到一行无法解析的数据(例如每个p是一行json日志)
	// for len(p) > b.Available() && b.err == nil {
	// 	var n int
	// 	if b.Buffered() == 0 {
	// 		// Large write, empty buffer.
	// 		// Write directly from p to avoid copy.
	// 		n, b.err = b.wr.Write(p)
	// 	} else {
	// 		n = copy(b.buf[b.n:], p)
	// 		b.n += n
	// 		b.Flush()
	// 	}
	// 	nn += n
	// 	p = p[n:]
	// }

	// UP: 改造一下逻辑，如果 len(p) > b.Available() 就将buf 和 p 都写入 b.wr
	if len(p) > b.Available() && b.err == nil {
		nn = b.Buffered()
		if nn > 0 {
			_ = b.Flush()
			if b.err != nil {
				return nn, b.err
			}
		}

		var n int
		n, b.err = b.wr.Write(p)
		if b.err != nil {
			return nn, b.err
		}

		nn += n
		return nn, nil
	}

	if b.err != nil {
		return nn, b.err
	}

	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}

// WriteString to writer
func (b *LineWriter) WriteString(s string) (int, error) {
	return b.Write([]byte(s))
}

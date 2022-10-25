package bufwrite_test

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog/bufwrite"
)

func TestNewBufIOWriter_WriteString(t *testing.T) {
	w := new(bytes.Buffer)
	bw := bufwrite.NewBufIOWriterSize(w, 12)

	_, err := bw.WriteString("hello, ")
	assert.NoErr(t, err)
	assert.Eq(t, 0, w.Len())

	_, err = bw.WriteString("worlds. oh")
	assert.NoErr(t, err)
	assert.Eq(t, "hello, world", w.String()) // different the LineWriter

	assert.NoErr(t, bw.Close())
	assert.Eq(t, "hello, worlds. oh", w.String())
}

type closeWriter struct {
	errOnWrite bool
	errOnClose bool
}

func (w *closeWriter) Close() error {
	if w.errOnClose {
		return errorx.Raw("close error")
	}
	return nil
}

func (w *closeWriter) Write(p []byte) (n int, err error) {
	if w.errOnWrite {
		return 0, errorx.Raw("write error")
	}
	return len(p), nil
}

func TestBufIOWriter_Close_error(t *testing.T) {
	bw := bufwrite.NewBufIOWriterSize(&closeWriter{errOnWrite: true}, 24)
	_, err := bw.WriteString("hi")
	assert.NoErr(t, err)

	// flush write error
	err = bw.Close()
	assert.Err(t, err)
	assert.Eq(t, "write error", err.Error())

	bw = bufwrite.NewBufIOWriterSize(&closeWriter{errOnClose: true}, 24)

	// close error
	err = bw.Close()
	assert.Err(t, err)
	assert.Eq(t, "close error", err.Error())
}

func TestBufIOWriter_Sync(t *testing.T) {
	w := new(bytes.Buffer)
	bw := bufwrite.NewBufIOWriter(w)

	_, err := bw.WriteString("hello")
	assert.NoErr(t, err)
	assert.Eq(t, 0, w.Len())
	assert.Eq(t, "", w.String())

	assert.NoErr(t, bw.Sync())
	assert.Eq(t, "hello", w.String())
}

func TestNewLineWriter(t *testing.T) {
	w := new(bytes.Buffer)
	bw := bufwrite.NewLineWriter(w)

	assert.True(t, bw.Size() > 0)
	assert.NoErr(t, bw.Flush())

	_, err := bw.WriteString("hello")
	assert.NoErr(t, err)
	assert.Eq(t, "", w.String())

	assert.NoErr(t, bw.Sync())
	assert.Eq(t, "hello", w.String())

	bw.Reset(w)
}

func TestLineWriter_Write_error(t *testing.T) {
	w := &closeWriter{errOnWrite: true}
	bw := bufwrite.NewLineWriterSize(w, 6)

	_, err := bw.WriteString("hello, tom")
	assert.Err(t, err)
	assert.Eq(t, "write error", err.Error())

	// get old error
	w.errOnWrite = false

	_, err = bw.WriteString("hello, wo")
	assert.Err(t, err)
	assert.Eq(t, "write error", err.Error())

	bw.Reset(w)

	_, err = bw.WriteString("hello")
	assert.NoErr(t, err)
}

func TestLineWriter_Flush_error(t *testing.T) {
	w := &closeWriter{}
	bw := bufwrite.NewLineWriterSize(w, 6)

	_, err := bw.WriteString("hello")
	assert.NoErr(t, err)

	// write error on flush
	w.errOnWrite = true
	err = bw.Flush()
	assert.Err(t, err)
	assert.Eq(t, "write error", err.Error())

	// get old error
	w.errOnWrite = false

	err = bw.Flush()
	assert.Err(t, err)
	assert.Eq(t, "write error", err.Error())

	bw.Reset(w)

	_, err = bw.WriteString("hello")
	assert.NoErr(t, err)
}

func TestLineWriter_Close_error(t *testing.T) {
	w := &closeWriter{}
	bw := bufwrite.NewLineWriterSize(w, 8)

	_, err := bw.WriteString("hello")
	assert.NoErr(t, err)

	// error on flush
	w.errOnWrite = true
	err = bw.Close()
	assert.Err(t, err)
	assert.Eq(t, "write error", err.Error())

	w = &closeWriter{errOnClose: true}
	bw = bufwrite.NewLineWriterSize(w, 8)

	err = bw.Close()
	assert.Err(t, err)
	assert.Eq(t, "close error", err.Error())
}

func TestNewLineWriterSize(t *testing.T) {
	w := new(bytes.Buffer)
	bw := bufwrite.NewLineWriterSize(w, 12)

	_, err := bw.WriteString("hello, ")
	assert.NoErr(t, err)
	assert.Eq(t, 0, w.Len())
	assert.True(t, bw.Size() > 0)

	_, err = bw.WriteString("worlds. oh")
	assert.NoErr(t, err)
	assert.Eq(t, "hello, worlds. oh", w.String()) // different the BufIOWriter

	_, err = bw.WriteString("...")
	assert.NoErr(t, err)
	assert.NoErr(t, bw.Close())
	assert.Eq(t, "hello, worlds. oh...", w.String())
	w.Reset()

	bw = bufwrite.NewLineWriterSize(bw, 8)
	assert.Eq(t, 12, bw.Size())

	bw = bufwrite.NewLineWriterSize(w, -12)
	assert.True(t, bw.Size() > 12)
}

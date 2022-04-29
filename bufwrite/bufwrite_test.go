package bufwrite_test

import (
	"bytes"
	"testing"

	"github.com/gookit/slog/bufwrite"
	"github.com/stretchr/testify/assert"
)

func TestNewBufIOWriter_WriteString(t *testing.T) {
	w := new(bytes.Buffer)
	bw := bufwrite.NewBufIOWriterSize(w, 12)

	_, err := bw.WriteString("hello, ")
	assert.NoError(t, err)
	assert.Equal(t, 0, w.Len())

	_, err = bw.WriteString("worlds. oh")
	assert.NoError(t, err)
	assert.Equal(t, "hello, world", w.String()) // different the LineWriter

	assert.NoError(t, bw.Close())
	assert.Equal(t, "hello, worlds. oh", w.String())
}

func TestNewBufIOWriter_Sync(t *testing.T) {
	w := new(bytes.Buffer)
	bw := bufwrite.NewBufIOWriter(w)

	_, err := bw.WriteString("hello")
	assert.NoError(t, err)
	assert.Equal(t, 0, w.Len())
	assert.Equal(t, "", w.String())

	assert.NoError(t, bw.Sync())
	assert.Equal(t, "hello", w.String())
}

func TestNewLineWriterSize(t *testing.T) {
	w := new(bytes.Buffer)
	bw := bufwrite.NewLineWriterSize(w, 12)

	_, err := bw.WriteString("hello, ")
	assert.NoError(t, err)
	assert.Equal(t, 0, w.Len())

	_, err = bw.WriteString("worlds. oh")
	assert.NoError(t, err)
	assert.Equal(t, "hello, worlds. oh", w.String()) // different the BufIOWriter

	_, err = bw.WriteString("...")
	assert.NoError(t, err)
	assert.NoError(t, bw.Close())
	assert.Equal(t, "hello, worlds. oh...", w.String())
}

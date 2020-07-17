package slog_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
)

func TestAddProcessor(t *testing.T) {
	buf := new(bytes.Buffer)

	l := slog.JSONSugaredLogger(buf, slog.ErrorLevel)
	l.AddProcessor(slog.AddHostname())
	l.Info("message")

	hostname,_ := os.Hostname()

	// {"channel":"application","data":{},"datetime":"2020/07/17 12:01:35","extra":{},"hostname":"InhereMac","level":"INFO","message":"message"}
	str := buf.String()
	buf.Reset()
	assert.Contains(t, str, `"level":"INFO"`)
	assert.Contains(t, str, `"message":"message"`)
	assert.Contains(t, str, fmt.Sprintf(`"hostname":"%s"`, hostname))

	l.ResetProcessors()
	l.PushProcessor(slog.MemoryUsage)
	l.Info("message2")

	// {"channel":"application","data":{},"datetime":"2020/07/16 16:40:18","extra":{"memoryUsage":326072},"level":"INFO","message":"message2"}
	str = buf.String()
	buf.Reset()
	assert.Contains(t, str, `"message":"message2"`)
	assert.Contains(t, str, `"memoryUsage":`)
}

package slog

import (
	"sync"
	"testing"
	"time"
)

// holdReadWriter enlarges the buffer-reuse race window: it keeps reading the
// passed-in bytes for a short while, during which a concurrent Format() from
// another logger could have grabbed the same pooled buffer and overwritten it.
type holdReadWriter struct{ acc byte }

func (w *holdReadWriter) Write(p []byte) (int, error) {
	time.Sleep(30 * time.Microsecond)
	for _, b := range p { // actively read every byte of the returned slice
		w.acc += b
	}
	return len(p), nil
}

// TestFormatter_ConcurrentLoggers_NoRace guards against returning a pooled
// buffer that is shared across loggers. Run with `go test -race`.
func TestFormatter_ConcurrentLoggers_NoRace(t *testing.T) {
	mk := func(useJSON bool) *SugaredLogger {
		l := NewSugaredLogger(&holdReadWriter{}, DebugLevel)
		l.ReportCaller = false
		if useJSON {
			l.Formatter = NewJSONFormatter()
		}
		return l
	}
	la, lb := mk(false), mk(true)

	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			lg := la
			if id%2 == 0 {
				lg = lb
			}
			for i := 0; i < 200; i++ {
				lg.Info("concurrent log message", id, i, "abcdefghijklmnop")
			}
		}(g)
	}
	wg.Wait()
}

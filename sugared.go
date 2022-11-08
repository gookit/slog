package slog

import (
	"io"
	"os"

	"github.com/gookit/color"
)

// SugaredLogger definition.
// Is a fast and usable Logger, which already contains
// the default formatting and handling capabilities
type SugaredLogger struct {
	*Logger
	// Formatter log message formatter. default use TextFormatter
	Formatter Formatter
	// Output writer
	Output io.Writer
	// Level for log handling. if log record level <= Level, it will be record.
	Level Level
}

// NewStdLogger instance
func NewStdLogger() *SugaredLogger {
	return NewSugaredLogger(os.Stdout, DebugLevel).Config(func(sl *SugaredLogger) {
		sl.SetName("stdLogger")
		// sl.CallerSkip += 1
		sl.ReportCaller = true
		// auto enable console color
		sl.Formatter.(*TextFormatter).EnableColor = color.SupportColor()
	})
}

// NewSugaredLogger create new SugaredLogger
func NewSugaredLogger(output io.Writer, level Level) *SugaredLogger {
	sl := &SugaredLogger{
		Level:  level,
		Output: output,
		Logger: New(),
		// default value
		Formatter: NewTextFormatter(),
	}

	// NOTICE: use self as an log handler
	sl.SetDefault(sl)

	return sl
}

// NewJSONSugared create new SugaredLogger with JSONFormatter
func NewJSONSugared(out io.Writer, level Level) *SugaredLogger {
	sl := NewSugaredLogger(out, level)
	sl.Formatter = NewJSONFormatter()

	return sl
}

// Config current logger
func (sl *SugaredLogger) Config(fn func(sl *SugaredLogger)) *SugaredLogger {
	fn(sl)
	return sl
}

// Configure current logger
func (sl *SugaredLogger) Configure(fn func(sl *SugaredLogger)) *SugaredLogger {
	fn(sl)
	return sl
}

// Reset the logger
func (sl *SugaredLogger) Reset() {
	*sl = *NewSugaredLogger(os.Stdout, DebugLevel)
}

// IsHandling Check if the current level can be handling
func (sl *SugaredLogger) IsHandling(level Level) bool {
	return sl.Level.ShouldHandling(level)
}

// Handle log record
func (sl *SugaredLogger) Handle(record *Record) error {
	bts, err := sl.Formatter.Format(record)
	if err != nil {
		return err
	}

	_, err = sl.Output.Write(bts)
	return err
}

// Close all log handlers
func (sl *SugaredLogger) Close() error {
	_ = sl.Logger.VisitAll(func(handler Handler) error {
		// TIP: must exclude self
		if _, ok := handler.(*SugaredLogger); !ok {
			err := handler.Close()
			if err != nil {
				sl.err = err
			}
		}
		return nil
	})

	return sl.err
}

// Flush all logs. alias of the FlushAll()
func (sl *SugaredLogger) Flush() error {
	return sl.FlushAll()
}

// FlushAll all logs
func (sl *SugaredLogger) FlushAll() error {
	return sl.Logger.VisitAll(func(handler Handler) error {
		if _, ok := handler.(*SugaredLogger); !ok {
			_ = handler.Flush()
		}
		return nil
	})
}

package formatter

import "github.com/gookit/slog"

// Formattable definition
type Formattable struct {
	formatter slog.Formatter
}

// DefaultFormatter setting
var DefaultFormatter = NewLineFormatter()

// Formatter get formatter
func (f *Formattable) Formatter() slog.Formatter {
	if f.formatter == nil {
		f.formatter = DefaultFormatter
	}
	return f.formatter
}

// SetFormatter to handler
func (f *Formattable) SetFormatter(formatter slog.Formatter) {
	f.formatter = formatter
}

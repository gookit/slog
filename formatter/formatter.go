package formatter

import "github.com/gookit/slog"

// Formattable definition
type Formattable struct {
	formatter slog.Formatter
}

var defaultFormatter = NewLineFormatter()

func (f *Formattable) Formatter() slog.Formatter {
	if f.formatter == nil {
		f.formatter = defaultFormatter
	}
	return f.formatter
}

func (f *Formattable) SetFormatter(formatter slog.Formatter) {
	f.formatter = formatter
}

package slog

// Formattable definition
type Formattable struct {
	formatter Formatter
}

// DefaultFormatter setting
var DefaultFormatter = NewLineFormatter()

// Formatter get formatter
func (f *Formattable) Formatter() Formatter {
	if f.formatter == nil {
		f.formatter = DefaultFormatter
	}
	return f.formatter
}

// SetFormatter to handler
func (f *Formattable) SetFormatter(formatter Formatter) {
	f.formatter = formatter
}

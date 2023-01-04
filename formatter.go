package slog

//
// Formatter interface
//

// Formatter interface
type Formatter interface {
	// Format you can format record and write result to record.Buffer
	Format(record *Record) ([]byte, error)
}

// FormatterFunc wrapper definition
type FormatterFunc func(r *Record) ([]byte, error)

// Format a log record
func (fn FormatterFunc) Format(r *Record) ([]byte, error) {
	return fn(r)
}

// Formattable interface
type Formattable interface {
	// Formatter get the log formatter
	Formatter() Formatter
	// SetFormatter set the log formatter
	SetFormatter(Formatter)
}

// FormattableTrait definition
type FormattableTrait struct {
	formatter Formatter
}

// Formatter get formatter. if not set, will return TextFormatter
func (f *FormattableTrait) Formatter() Formatter {
	if f.formatter == nil {
		f.formatter = NewTextFormatter()
	}
	return f.formatter
}

// SetFormatter to handler
func (f *FormattableTrait) SetFormatter(formatter Formatter) {
	f.formatter = formatter
}

// Format log record to bytes
func (f *FormattableTrait) Format(record *Record) ([]byte, error) {
	return f.Formatter().Format(record)
}

// AsTextFormatter util func
func AsTextFormatter(f Formatter) *TextFormatter {
	if tf, ok := f.(*TextFormatter); ok {
		return tf
	}
	panic("slog: cannot cast input as *TextFormatter")
}

// AsJSONFormatter util func
func AsJSONFormatter(f Formatter) *JSONFormatter {
	if jf, ok := f.(*JSONFormatter); ok {
		return jf
	}
	panic("slog: cannot cast input as *JSONFormatter")
}

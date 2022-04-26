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
type FormatterFunc func(r *Record) error

// Format a record
func (fn FormatterFunc) Format(r *Record) error {
	return fn(r)
}

// FormattableHandler interface
type FormattableHandler interface {
	// Formatter get the log formatter
	Formatter() Formatter
	// SetFormatter set the log formatter
	SetFormatter(Formatter)
}

// Formattable definition
type Formattable struct {
	formatter Formatter
}

// Formatter get formatter. if not set, will return TextFormatter
func (f *Formattable) Formatter() Formatter {
	if f.formatter == nil {
		f.formatter = NewTextFormatter()
	}
	return f.formatter
}

// SetFormatter to handler
func (f *Formattable) SetFormatter(formatter Formatter) {
	f.formatter = formatter
}

// FormatRecord to bytes
func (f *Formattable) FormatRecord(record *Record) ([]byte, error) {
	return f.Formatter().Format(record)
}

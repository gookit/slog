package formatter

import (
	"github.com/gookit/slog"
)

// JSONFormatter definition
type JSONFormatter struct {
	fieldMap slog.FieldMap
}

// Format an log record
func (f *JSONFormatter) Format(r *slog.Record) ([]byte,error) {
	panic("implement me")
}


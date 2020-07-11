package formatter

import (
	"github.com/gookit/slog"
)

// JSONFormatter definition
type JSONFormatter struct {

}

func (f *JSONFormatter) Format(r *slog.Record) ([]byte, error) {
	panic("implement me")
}


package formatter

import (
	"encoding/json"
	"time"

	"github.com/gookit/slog"
)

// JSONFormatter definition
type JSONFormatter struct {
	// FieldMap for output
	// item: `"field" : "output name"`
	// eg: {"message": "message"}
	FieldMap slog.FieldMap
}

// NewJSONFormatter create new JSONFormatter
func NewJSONFormatter(fieldMap slog.FieldMap) *JSONFormatter {
	return &JSONFormatter{FieldMap: fieldMap}
}

// Format an log record
func (f *JSONFormatter) Format(r *slog.Record) ([]byte,error) {
	tplData := make(slog.M, len(f.FieldMap))

	for field, outName := range f.FieldMap {
		switch {
		case field == FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			tplData[outName] = r.Time.Format(time.RFC3339)
		case field == FieldKeyLevel:
			tplData[outName] = r.LevelName
		case field == FieldKeyChannel:
			tplData[outName] = r.Channel
		case field == FieldKeyMsg:
			tplData[outName] = r.Message
		case field == FieldKeyData:
			tplData[outName] = r.Data
		case field == FieldKeyExtra:
			tplData[outName] = r.Extra
		default:
			tplData[outName] = r.Fields[field]
		}
	}

	bs, err := json.Marshal(tplData)
	// with newline
	bs = append(bs, '\n')

	return bs, err
}


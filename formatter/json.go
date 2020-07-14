package formatter

import (
	"encoding/json"
	"time"

	"github.com/gookit/slog"
)

// JSONFormatter definition
type JSONFormatter struct {
	// Fields exported log fields.
	Fields []string
	// Aliases for output fields. you can change export field name.
	// item: `"field" : "output name"`
	// eg: {"message": "msg"} export field will display "msg"
	Aliases slog.StringMap

	// PrettyPrint will indent all json logs
	PrettyPrint bool
}

// NewJSONFormatter create new JSONFormatter
func NewJSONFormatter(aliases slog.StringMap) *JSONFormatter {
	return &JSONFormatter{
		Aliases: aliases,
		Fields: slog.DefaultFields,
	}
}

// Format an log record
func (f *JSONFormatter) Format(r *slog.Record) ([]byte,error) {
	logData := make(slog.M, len(f.Fields))
	for _, field := range f.Fields {
		outName, ok := f.Aliases[field]
		if !ok {
			outName = field
		}

		switch {
		case field == slog.FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			logData[outName] = r.Time.Format(time.RFC3339)
		case field == slog.FieldKeyLevel:
			logData[outName] = r.LevelName
		case field == slog.FieldKeyChannel:
			logData[outName] = r.Channel
		case field == slog.FieldKeyMessage:
			logData[outName] = r.Message
		case field == slog.FieldKeyData:
			logData[outName] = r.Data
		case field == slog.FieldKeyExtra:
			logData[outName] = r.Extra
		default:
			logData[outName] = r.Fields[field]
		}
	}

	// sort.Interface()

	buffer := r.NewBuffer()
	encoder := json.NewEncoder(buffer)

	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}

	// has been add newline.
	err := encoder.Encode(logData)
	// if err == nil {
		// with newline
		// buffer.Write([]byte("\n"))
	// }

	return buffer.Bytes(), err
}


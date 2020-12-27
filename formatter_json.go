package slog

import (
	"encoding/json"
	"time"
)

// JSONFormatter definition
type JSONFormatter struct {
	// Fields exported log fields.
	Fields []string
	// Aliases for output fields. you can change export field name.
	// item: `"field" : "output name"`
	// eg: {"message": "msg"} export field will display "msg"
	Aliases StringMap

	// PrettyPrint will indent all json logs
	PrettyPrint bool
	// TimeFormat the time format layout. default is time.RFC3339
	TimeFormat string
}

// NewJSONFormatter create new JSONFormatter
func NewJSONFormatter(fn ...func(*JSONFormatter)) *JSONFormatter {
	f := &JSONFormatter{
		// Aliases: make(StringMap, 0),
		Fields:     DefaultFields,
		TimeFormat: DefaultTimeFormat,
	}

	if len(fn) > 0 {
		fn[0](f)
	}

	return f
}

// Configure current formatter
func (f *JSONFormatter) Configure(fn func(*JSONFormatter)) *JSONFormatter {
	fn(f)
	return f
}

// Format an log record
func (f *JSONFormatter) Format(r *Record) ([]byte, error) {
	logData := make(M, len(f.Fields))

	for _, field := range f.Fields {
		outName, ok := f.Aliases[field]
		if !ok {
			outName = field
		}

		switch {
		case field == FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			logData[outName] = r.Time.Format(f.TimeFormat)
		case field == FieldKeyCaller && r.Caller != nil:
			logData[outName] = formatCaller(r.Caller, field) // "logger_test.go:48.TestLogger_ReportCaller"
		case field == FieldKeyFLine && r.Caller != nil:
			logData[outName] = formatCaller(r.Caller, field) // "logger_test.go:48"
		case field == FieldKeyFunc && r.Caller != nil:
			logData[outName] = r.Caller.Function // "github.com/gookit/slog_test.TestLogger_ReportCaller"
		case field == FieldKeyFile && r.Caller != nil:
			logData[outName] = formatCaller(r.Caller, field) // "/work/go/gookit/slog/logger_test.go:48"
		case field == FieldKeyLevel:
			logData[outName] = r.LevelName
		case field == FieldKeyChannel:
			logData[outName] = r.Channel
		case field == FieldKeyMessage:
			logData[outName] = r.Message
		case field == FieldKeyData:
			logData[outName] = r.Data
		case field == FieldKeyExtra:
			logData[outName] = r.Extra
			// default:
			// 	logData[outName] = r.Fields[field]
		}
	}

	// exported custom fields
	for field, value := range r.Fields {
		fieldKey := field
		if _, has := logData[field]; has {
			fieldKey = "fields." + field
		}

		logData[fieldKey] = value
	}

	// sort.Interface()

	buffer := r.NewBuffer()
	encoder := json.NewEncoder(buffer)

	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}

	// has been add newline in Encode().
	err := encoder.Encode(logData)
	return buffer.Bytes(), err
}

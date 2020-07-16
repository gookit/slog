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
func NewJSONFormatter(aliases StringMap) *JSONFormatter {
	return &JSONFormatter{
		Aliases: aliases,
		Fields:  DefaultFields,
		TimeFormat: DefaultTimeFormat,
	}
}

// Format an log record
func (f *JSONFormatter) Format(r *Record) ([]byte,error) {
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

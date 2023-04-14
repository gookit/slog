package slog

import (
	"encoding/json"

	"github.com/valyala/bytebufferpool"
)

var (
	// DefaultFields default log export fields for json formatter.
	DefaultFields = []string{
		FieldKeyDatetime,
		FieldKeyChannel,
		FieldKeyLevel,
		FieldKeyCaller,
		FieldKeyMessage,
		FieldKeyData,
		FieldKeyExtra,
	}

	// NoTimeFields log export fields without time
	NoTimeFields = []string{
		FieldKeyChannel,
		FieldKeyLevel,
		FieldKeyMessage,
		FieldKeyData,
		FieldKeyExtra,
	}
)

// JSONFormatter definition
type JSONFormatter struct {
	// Fields exported log fields. default is DefaultFields
	Fields []string
	// Aliases for output fields. you can change export field name.
	//
	// item: `"field" : "output name"`
	// eg: {"message": "msg"} export field will display "msg"
	Aliases StringMap

	// PrettyPrint will indent all json logs
	PrettyPrint bool
	// TimeFormat the time format layout. default is DefaultTimeFormat
	TimeFormat string
}

// NewJSONFormatter create new JSONFormatter
func NewJSONFormatter(fn ...func(f *JSONFormatter)) *JSONFormatter {
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

// AddField for export
func (f *JSONFormatter) AddField(name string) *JSONFormatter {
	f.Fields = append(f.Fields, name)
	return f
}

var jsonPool bytebufferpool.Pool

// Format an log record
func (f *JSONFormatter) Format(r *Record) ([]byte, error) {
	logData := make(M, len(f.Fields))

	// TODO perf: use buf write build JSON string.
	for _, field := range f.Fields {
		outName, ok := f.Aliases[field]
		if !ok {
			outName = field
		}

		switch {
		case field == FieldKeyDatetime:
			logData[outName] = r.Time.Format(f.TimeFormat)
		case field == FieldKeyTimestamp:
			logData[outName] = r.timestamp()
		case field == FieldKeyCaller && r.Caller != nil:
			logData[outName] = formatCaller(r.Caller, r.CallerFlag)
		case field == FieldKeyLevel:
			logData[outName] = r.LevelName()
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
	buf := jsonPool.Get()
	// buf.Reset()
	defer jsonPool.Put(buf)
	// buf := r.NewBuffer()
	// buf.Reset()
	// buf.Grow(256)

	encoder := json.NewEncoder(buf)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}

	// has been added newline in Encode().
	err := encoder.Encode(logData)
	return buf.Bytes(), err
}

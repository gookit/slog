package slog

import (
	"github.com/gookit/color"
	"github.com/valyala/bytebufferpool"
)

// there are built in text log template
const (
	DefaultTemplate = "[{{datetime}}] [{{channel}}] [{{level}}] [{{caller}}] {{message}} {{data}} {{extra}}\n"
	NamedTemplate   = "{{datetime}} channel={{channel}} level={{level}} [file={{caller}}] message={{message}} data={{data}}\n"
)

// ColorTheme for format log to console
var ColorTheme = map[Level]color.Color{
	PanicLevel:  color.FgRed,
	FatalLevel:  color.FgRed,
	ErrorLevel:  color.FgMagenta,
	WarnLevel:   color.FgYellow,
	NoticeLevel: color.OpBold,
	InfoLevel:   color.FgGreen,
	DebugLevel:  color.FgCyan,
	// TraceLevel:  color.FgLightGreen,
}

// TextFormatter definition
type TextFormatter struct {
	// template text template for render output log messages
	template string
	// fields list, parsed from template string.
	//
	// NOTE: contains no-field items in the list. eg: ["level", "}}"}
	fields []string

	// TimeFormat the time format layout. default is DefaultTimeFormat
	TimeFormat string
	// Enable color on print log to terminal
	EnableColor bool
	// ColorTheme setting on render color on terminal
	ColorTheme map[Level]color.Color
	// FullDisplay Whether to display when record.Data, record.Extra, etc. are empty
	FullDisplay bool
	// EncodeFunc data encode for Record.Data, Record.Extra, etc.
	//
	// Default is encode by EncodeToString()
	EncodeFunc func(v any) string
	// CallerFormatFunc the caller format layout. default is defined by CallerFlag
	CallerFormatFunc CallerFormatFn

	// TODO BeforeFunc call it before format, update fields or other
	// BeforeFunc func(r *Record)
}

// TextFormatterFn definition
type TextFormatterFn func(*TextFormatter)

// NewTextFormatter create new TextFormatter
func NewTextFormatter(template ...string) *TextFormatter {
	var fmtTpl string
	if len(template) > 0 {
		fmtTpl = template[0]
	} else {
		fmtTpl = DefaultTemplate
	}

	f := &TextFormatter{
		// default options
		ColorTheme: ColorTheme,
		TimeFormat: DefaultTimeFormat,
		// EnableColor: color.SupportColor(),
		// EncodeFunc: func(v any) string {
		// 	return fmt.Sprint(v)
		// },
		EncodeFunc: EncodeToString,
	}
	f.SetTemplate(fmtTpl)

	return f
}

// TextFormatterWith create new TextFormatter with options
func TextFormatterWith(fns ...TextFormatterFn) *TextFormatter {
	return NewTextFormatter().WithOptions(fns...)
}

// Configure the formatter
func (f *TextFormatter) Configure(fn TextFormatterFn) *TextFormatter {
	return f.WithOptions(fn)
}

// WithOptions func on the formatter
func (f *TextFormatter) WithOptions(fns ...TextFormatterFn) *TextFormatter {
	for _, fn := range fns {
		fn(f)
	}
	return f
}

// SetTemplate set the log format template and update field-map
func (f *TextFormatter) SetTemplate(fmtTpl string) {
	f.template = fmtTpl
	f.fields = parseTemplateToFields(fmtTpl)
}

// Template get
func (f *TextFormatter) Template() string {
	return f.template
}

// WithEnableColor enable color on print log to terminal
func (f *TextFormatter) WithEnableColor(enable bool) *TextFormatter {
	f.EnableColor = enable
	return f
}

// Fields get an export field list
func (f *TextFormatter) Fields() []string {
	ss := make([]string, 0, len(f.fields)/2)
	for _, s := range f.fields {
		if s[0] >= 'a' && s[0] <= 'z' {
			ss = append(ss, s)
		}
	}
	return ss
}

var textPool bytebufferpool.Pool

// Format a log record
//
//goland:noinspection GoUnhandledErrorResult
func (f *TextFormatter) Format(r *Record) ([]byte, error) {
	f.beforeFormat()
	buf := textPool.Get()
	defer textPool.Put(buf)

	for _, field := range f.fields {
		// is not field name. eg: "}}] "
		if field[0] < 'a' || field[0] > 'z' {
			// remove left "}}"
			if len(field) > 1 && field[0:2] == "}}" {
				buf.WriteString(field[2:])
			} else {
				buf.WriteString(field)
			}
			continue
		}

		switch {
		case field == FieldKeyDatetime:
			buf.B = r.Time.AppendFormat(buf.B, f.TimeFormat)
		case field == FieldKeyTimestamp:
			buf.WriteString(r.timestamp())
		case field == FieldKeyCaller && r.Caller != nil:
			buf.WriteString(formatCaller(r.Caller, r.CallerFlag, f.CallerFormatFunc))
		case field == FieldKeyLevel:
			buf.WriteString(f.renderColorText(field, r.LevelName(), r.Level))
		case field == FieldKeyChannel:
			buf.WriteString(r.Channel)
		case field == FieldKeyMessage:
			buf.WriteString(f.renderColorText(field, r.Message, r.Level))
		case field == FieldKeyData:
			if f.FullDisplay || len(r.Data) > 0 {
				buf.WriteString(f.EncodeFunc(r.Data))
			}
		case field == FieldKeyExtra:
			if f.FullDisplay || len(r.Extra) > 0 {
				buf.WriteString(f.EncodeFunc(r.Extra))
			}
		default:
			if _, ok := r.Fields[field]; ok {
				buf.WriteString(f.EncodeFunc(r.Fields[field]))
			} else {
				buf.WriteString(field)
			}
		}
	}

	// return buf.Bytes(), nil
	return buf.B, nil
}

func (f *TextFormatter) beforeFormat() {
	// if f.BeforeFunc == nil {}
	if f.EncodeFunc == nil {
		f.EncodeFunc = EncodeToString
	}
	if f.ColorTheme == nil {
		f.ColorTheme = ColorTheme
	}
}

func (f *TextFormatter) renderColorByLevel(s string, l Level) string {
	if theme, ok := f.ColorTheme[l]; ok {
		return theme.Render(s)
	}
	return s
}

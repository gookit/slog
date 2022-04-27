package slog

import (
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/valyala/bytebufferpool"
)

const DefaultTemplate = "[{{datetime}}] [{{channel}}] [{{level}}] [{{caller}}] {{message}} {{data}} {{extra}}\n"
const NamedTemplate = "{{datetime}} channel={{channel}} level={{level}} [file={{caller}}] message={{message}} data={{data}}\n"

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
	// Template text template for render output log messages
	Template string
	// fields list, parsed from Template string.
	// NOTICE: fields contains no-field items.
	// eg: ["level", "}}"}
	fields []string

	// TimeFormat the time format layout. default is time.RFC3339
	TimeFormat string
	// Enable color on print log to terminal
	EnableColor bool
	// ColorTheme setting on render color on terminal
	ColorTheme map[Level]color.Color
	// FullDisplay Whether to display when record.Data, record.Extra, etc. are empty
	FullDisplay bool
	// EncodeFunc data encode for Record.Data, Record.Extra, etc.
	// Default is encode by EncodeToString()
	EncodeFunc func(v interface{}) string
}

// NewTextFormatter create new TextFormatter
func NewTextFormatter(template ...string) *TextFormatter {
	var fmtTpl string
	if len(template) > 0 {
		fmtTpl = template[0]
	} else {
		fmtTpl = DefaultTemplate
	}

	return &TextFormatter{
		Template: fmtTpl,
		fields:   parseTemplateToFields(fmtTpl),
		// default options
		TimeFormat: DefaultTimeFormat,
		ColorTheme: ColorTheme,
		// EnableColor: color.SupportColor(),
		// EncodeFunc: func(v interface{}) string {
		// 	return fmt.Sprint(v)
		// },
		EncodeFunc: EncodeToString,
	}
}

// SetTemplate set the log format template and update field-map
func (f *TextFormatter) SetTemplate(fmtTpl string) {
	f.Template = fmtTpl
	f.fields = parseTemplateToFields(fmtTpl)
}

// Fields get export field list
func (f *TextFormatter) Fields() []string {
	ss := make([]string, 0, len(f.fields)/2)
	for _, s := range f.fields {
		if s[0] >= 'a' && s[0] <= 'z' {
			ss = append(ss, s)
		}
	}

	return ss
}

// Format a log record
//goland:noinspection GoUnhandledErrorResult
func (f *TextFormatter) Format(r *Record) ([]byte, error) {
	buf := bytebufferpool.Get()
	buf.Reset()
	defer bytebufferpool.Put(buf)

	for _, field := range f.fields {
		// is not field name.
		if field[0] < 'a' || field[0] > 'z' {
			// remove "}}"
			if len(field) > 2 && strings.HasPrefix(field, "}}") {
				buf.WriteString(field[2:])
			} else {
				buf.WriteString(field)
			}
			continue
		}

		switch {
		case field == FieldKeyDatetime:
			buf.WriteString(r.Time.Format(f.TimeFormat))
		case field == FieldKeyTimestamp:
			buf.WriteString(strconv.Itoa(r.MicroSecond()))
		case field == FieldKeyCaller && r.Caller != nil:
			buf.WriteString(formatCaller(r.Caller, r.CallerFlag))
		case field == FieldKeyLevel:
			// output colored logs for console
			if f.EnableColor {
				buf.WriteString(f.renderColorByLevel(r.LevelName(), r.Level))
			} else {
				buf.WriteString(r.LevelName())
			}
		case field == FieldKeyChannel:
			buf.WriteString(r.Channel)
		case field == FieldKeyMessage:
			// output colored logs for console
			if f.EnableColor {
				buf.WriteString(f.renderColorByLevel(r.Message, r.Level))
			} else {
				buf.WriteString(r.Message)
			}
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

func (f *TextFormatter) renderColorByLevel(text string, level Level) string {
	if theme, ok := f.ColorTheme[level]; ok {
		return theme.Render(text)
	}
	return text
}

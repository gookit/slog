package slog

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
)

const DefaultTemplate = "[{{datetime}}] [{{channel}}] [{{level}}] [{{caller}}] {{message}} {{data}} {{extra}}\n"
const NamedTemplate = "{{datetime}} channel={{channel}} level={{level}} [file={{caller}}] message={{message}} data={{data}}\n"

// ColorTheme for format log to console
var ColorTheme = map[Level]color.Color{
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
	// field map, parsed from format string.
	// eg: {"level": "{{level}}",}
	fieldMap StringMap

	// TimeFormat the time format layout. default is time.RFC3339
	TimeFormat string
	// Enable color on print log to terminal
	EnableColor bool
	// ColorTheme setting on render color on terminal
	ColorTheme map[Level]color.Color
	// FullDisplay Whether to display when record.Data, record.Extra, etc. are empty
	FullDisplay bool
	// EncodeFunc data encode for record.Data, record.Extra, etc.
	// Default is encode by `fmt.Sprint()`
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
		fieldMap: parseFieldMap(fmtTpl),
		// default options
		TimeFormat: DefaultTimeFormat,
		ColorTheme: ColorTheme,
		EncodeFunc: func(v interface{}) string {
			return fmt.Sprint(v)
		},
	}
}

// SetTemplate set the log format template and update field-map
func (f *TextFormatter) SetTemplate(fmtTpl string) {
	f.Template = fmtTpl
	f.fieldMap = parseFieldMap(fmtTpl)
}

// FieldMap get export field map
func (f *TextFormatter) FieldMap() StringMap {
	return f.fieldMap
}

// Format an log record
func (f *TextFormatter) Format(r *Record) ([]byte, error) {
	// output colored logs for console
	if f.EnableColor {
		return f.formatWithColor(r)
	}

	return f.formatNoColor(r)
}

func (f *TextFormatter) formatNoColor(r *Record) ([]byte, error) {
	tplData := make(map[string]string, len(f.fieldMap))
	for field, tplVar := range f.fieldMap {
		switch {
		case field == FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			tplData[tplVar] = r.Time.Format(f.TimeFormat)
		case field == FieldKeyCaller && r.Caller != nil: // caller eg: "logger_test.go:48"
			tplData[tplVar] = formatCaller(r.Caller, false)
		case field == FieldKeyFunc && r.Caller != nil:
			tplData[tplVar] = r.Caller.Function // "github.com/gookit/slog_test.TestLogger_ReportCaller"
		case field == FieldKeyFile && r.Caller != nil:
			tplData[tplVar] = formatCaller(r.Caller, true) // "/work/go/gookit/slog/logger_test.go:48"
		case field == FieldKeyLevel:
			tplData[tplVar] = r.LevelName
		case field == FieldKeyChannel:
			tplData[tplVar] = r.Channel
		case field == FieldKeyMessage:
			tplData[tplVar] = r.Message
		case field == FieldKeyData:
			if f.FullDisplay || len(r.Data) > 0 {
				tplData[tplVar] = f.EncodeFunc(r.Data)
			} else {
				tplData[tplVar] = ""
			}
		case field == FieldKeyExtra:
			if f.FullDisplay || len(r.Extra) > 0 {
				tplData[tplVar] = f.EncodeFunc(r.Extra)
			} else {
				tplData[tplVar] = ""
			}
		default:
			tplData[tplVar] = f.EncodeFunc(r.Fields[field])
		}
	}

	// TODO ... use r.Buffer
	// strings.NewReplacer().WriteString(buf)

	str := strutil.Replaces(f.Template, tplData)
	return []byte(str), nil
}

func (f *TextFormatter) formatWithColor(r *Record) ([]byte, error) {
	tplData := make(map[string]string, len(f.fieldMap))
	for field, tplVar := range f.fieldMap {
		switch {
		case field == FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			tplData[tplVar] = r.Time.Format(f.TimeFormat)
		case field == FieldKeyCaller && r.Caller != nil: // caller eg: "logger_test.go:48"
			tplData[tplVar] = formatCaller(r.Caller, false)
		case field == FieldKeyFunc && r.Caller != nil:
			tplData[tplVar] = r.Caller.Function // "github.com/gookit/slog_test.TestLogger_ReportCaller"
		case field == FieldKeyFile && r.Caller != nil:
			tplData[tplVar] = formatCaller(r.Caller, true) // "/work/go/gookit/slog/logger_test.go:48"
		case field == FieldKeyLevel:
			tplData[tplVar] = f.renderColorByLevel(r.LevelName, r.Level)
		case field == FieldKeyChannel:
			tplData[tplVar] = r.Channel
		case field == FieldKeyMessage:
			tplData[tplVar] = f.renderColorByLevel(r.Message, r.Level)
			// tplData[tplVar] = r.Message
			// if r.Level <= NoticeLevel {
			//
			// }
		case field == FieldKeyData:
			if f.FullDisplay || len(r.Data) > 0 {
				tplData[tplVar] = f.EncodeFunc(r.Data)
			} else {
				tplData[tplVar] = ""
			}
		case field == FieldKeyExtra:
			if f.FullDisplay || len(r.Extra) > 0 {
				tplData[tplVar] = f.EncodeFunc(r.Extra)
			} else {
				tplData[tplVar] = ""
			}
		default:
			tplData[tplVar] = f.EncodeFunc(r.Fields[field])
		}
	}

	str := strutil.Replaces(f.Template, tplData)
	return []byte(str), nil
}

func (f *TextFormatter) renderColorByLevel(text string, level Level) string {
	if theme, ok := f.ColorTheme[level]; ok {
		return theme.Render(text)
	}

	return text
}

// parse string "{{channel}}" to map { "channel": "{{channel}}" }
func parseFieldMap(format string) StringMap {
	rgp := regexp.MustCompile(`{{\w+}}`)

	ss := rgp.FindAllString(format, -1)
	fm := make(StringMap)
	for _, tplVar := range ss {
		field := strings.Trim(tplVar, "{}")
		fm[field] = tplVar
	}

	return fm
}

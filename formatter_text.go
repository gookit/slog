package slog

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gookit/color"
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
	// EncodeFunc data encode for Record.Data, Record.Extra, etc.
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
	f.fieldMap = parseFieldMap(fmtTpl)
}

// FieldMap get export field map
func (f *TextFormatter) FieldMap() StringMap {
	return f.fieldMap
}

// Format an log record
func (f *TextFormatter) Format(r *Record) ([]byte, error) {
	oldnew := make([]string, 0, len(f.fieldMap) * 2 + 1)
	// tplData := make(map[string]string, len(f.fieldMap))
	for field, tplVar := range f.fieldMap {
		switch {
		case field == FieldKeyDatetime:
			oldnew = append(oldnew, tplVar, r.Time.Format(f.TimeFormat))
			// tplData[tplVar] = r.Time.Format(f.TimeFormat)
		case field == FieldKeyTimestamp:
			// tplData[tplVar] = strconv.Itoa(r.MicroSecond())
			oldnew = append(oldnew, tplVar, strconv.Itoa(r.MicroSecond()))
		case field == FieldKeyCaller && r.Caller != nil:
			// caller eg: "logger_test.go:48,TestLogger_ReportCaller"
			// tplData[tplVar] = formatCaller(r.Caller, field)
			oldnew = append(oldnew, tplVar, formatCaller(r.Caller, field))
		case field == FieldKeyFLine && r.Caller != nil:
			// "logger_test.go:48"
			// tplData[tplVar] = formatCaller(r.Caller, field)
			oldnew = append(oldnew, tplVar, formatCaller(r.Caller, field))
		case field == FieldKeyFunc && r.Caller != nil:
			// "github.com/gookit/slog_test.TestLogger_ReportCaller"
			// tplData[tplVar] = r.Caller.Function
			oldnew = append(oldnew, tplVar,  r.Caller.Function)
		case field == FieldKeyFile && r.Caller != nil:
			// "/work/go/gookit/slog/logger_test.go:48"
			// tplData[tplVar] = formatCaller(r.Caller, field)
			oldnew = append(oldnew, tplVar, formatCaller(r.Caller, field))
		case field == FieldKeyLevel:
			oldnew = append(oldnew, tplVar)
			// output colored logs for console
			if f.EnableColor {
				// tplData[tplVar] = f.renderColorByLevel(r.LevelName(), r.Level)
				oldnew = append(oldnew, f.renderColorByLevel(r.LevelName(), r.Level))
			} else {
				// tplData[tplVar] = r.LevelName()
				oldnew = append(oldnew, r.LevelName())
			}
		case field == FieldKeyChannel:
			// tplData[tplVar] = r.Channel
			oldnew = append(oldnew, tplVar, r.Channel)
		case field == FieldKeyMessage:
			// output colored logs for console
			oldnew = append(oldnew, tplVar)
			if f.EnableColor {
				// tplData[tplVar] = f.renderColorByLevel(r.Message, r.Level)
				oldnew = append(oldnew, f.renderColorByLevel(r.Message, r.Level))
			} else {
				// tplData[tplVar] = r.Message
				oldnew = append(oldnew, r.Message)
			}
		case field == FieldKeyData:
			if f.FullDisplay || len(r.Data) > 0 {
				// tplData[tplVar] = f.EncodeFunc(r.Data)
				oldnew = append(oldnew, tplVar, f.EncodeFunc(r.Data))
			} else {
				// tplData[tplVar] = ""
				oldnew = append(oldnew, tplVar, "")
			}
		case field == FieldKeyExtra:
			if f.FullDisplay || len(r.Extra) > 0 {
				// tplData[tplVar] = f.EncodeFunc(r.Extra)
				oldnew = append(oldnew, tplVar, f.EncodeFunc(r.Extra))
			} else {
				// tplData[tplVar] = ""
				oldnew = append(oldnew, tplVar, "")
			}
		default:
			if _, ok := r.Fields[field]; ok {
				// tplData[tplVar] = f.EncodeFunc(r.Fields[field])
				oldnew = append(oldnew, tplVar, f.EncodeFunc(r.Fields[field]))
			}
		}
	}

	// str := strutil.Replaces(f.Template, tplData)
	str := strings.NewReplacer(oldnew...).Replace(f.Template)
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

package slog

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
)

const DefaultTemplate = "[{{datetime}}] [{{channel}}] [{{level}}] {{message}} {{data}} {{extra}}\n"

// ColorTheme for format log to console
var ColorTheme = map[Level]color.Color{
	ErrorLevel: color.FgRed,
	WarnLevel:  color.FgYellow,
	InfoLevel:  color.FgGreen,
	DebugLevel: color.FgCyan,
	TraceLevel: color.FgMagenta,
}

// TextFormatter definition
type TextFormatter struct {
	template string
	// field map, parsed from format string.
	// eg: {"level": "{{level}}",}
	fieldMap StringMap

	// TimeFormat the time format layout. default is time.RFC3339
	TimeFormat string
	// Enable color on print log to terminal
	EnableColor bool
	// ColorTheme setting on render color on terminal
	ColorTheme  map[Level]color.Color
	// FullDisplay Whether to display when record.Data, record.Extra, etc. are empty
	FullDisplay bool
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
		template:   fmtTpl,
		fieldMap:   parseFieldMap(fmtTpl),
		TimeFormat: DefaultTimeFormat,
		ColorTheme: ColorTheme,
	}
}

// FieldMap get export field map
func (f *TextFormatter) FieldMap() StringMap {
	return f.fieldMap
}

// Format an log record
func (f *TextFormatter) Format(r *Record) ([]byte, error) {
	if f.EnableColor {
		return f.formatForConsole(r)
	}

	tplData := make(StringMap, len(f.fieldMap))
	for field, tplVar := range f.fieldMap {
		switch {
		case field == FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			tplData[tplVar] = r.Time.Format(f.TimeFormat)
		case field == FieldKeyLevel:
			tplData[tplVar] = r.LevelName
		case field == FieldKeyChannel:
			tplData[tplVar] = r.Channel
		case field == FieldKeyMessage:
			tplData[tplVar] = r.Message
		case field == FieldKeyData:
			if f.FullDisplay || len(r.Data) > 0{
				tplData[tplVar] = fmt.Sprint(r.Data)
			} else {
				tplData[tplVar] = ""
			}
		case field == FieldKeyExtra:
			if f.FullDisplay || len(r.Extra) > 0 {
				tplData[tplVar] = fmt.Sprint(r.Extra)
			} else {
				tplData[tplVar] = ""
			}
		default:
			tplData[tplVar] = fmt.Sprint(r.Fields[field])
		}
	}

	// dump.Println(tplData, r.LevelName)

	// TODO ... use r.Buffer
	// strings.NewReplacer().WriteString(buf)

	str := strutil.Replaces(f.template, tplData)

	return []byte(str), nil
}

func (f *TextFormatter) renderColorByLevel(text string, level Level) string {
	if theme, ok := f.ColorTheme[level]; ok {
		return theme.Render(text)
	}

	return text
}

func (f *TextFormatter) formatForConsole(r *Record) ([]byte, error) {
	tplData := make(StringMap, len(f.fieldMap))
	for field, tplVar := range f.fieldMap {
		switch {
		case field == FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			tplData[tplVar] = r.Time.Format(f.TimeFormat)
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
			if f.FullDisplay || len(r.Data) > 0{
				tplData[tplVar] = fmt.Sprint(r.Data)
			} else {
				tplData[tplVar] = ""
			}
		case field == FieldKeyExtra:
			if f.FullDisplay || len(r.Extra) > 0 {
				tplData[tplVar] = fmt.Sprint(r.Extra)
			} else {
				tplData[tplVar] = ""
			}
		default:
			tplData[tplVar] = fmt.Sprint(r.Fields[field])
		}
	}

	str := strutil.Replaces(f.template, tplData)

	return []byte(str), nil
}

// parse "{{channel}}" to { "channel": "{{channel}}" }
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

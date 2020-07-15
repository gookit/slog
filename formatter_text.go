package slog

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
)

const DefaultTemplate = "[{{datetime}}] {{channel}}.{{level}}: {{message}} {{data}} {{extra}}\n"

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
	ColorTheme  map[Level]color.Color
}

// NewLineFormatter create new TextFormatter
func NewLineFormatter(template ...string) *TextFormatter {
	var fmtTpl string
	if len(template) > 0 {
		fmtTpl = template[0]
	} else {
		fmtTpl = DefaultTemplate
	}

	return &TextFormatter{
		template:     fmtTpl,
		fieldMap:   parseFieldMap(fmtTpl),
		TimeFormat: time.RFC3339,
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
			tplData[tplVar] = fmt.Sprint(r.Data)
		case field == FieldKeyExtra:
			tplData[tplVar] = fmt.Sprint(r.Extra)
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

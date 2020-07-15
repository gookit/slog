package slog

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
)

const SimpleFormat = "[{{datetime}}] {{channel}}.{{level}}: {{message}} {{data}} {{extra}}\n"

// ColorTheme for format log to console
var ColorTheme = map[Level]color.Color{
	ErrorLevel: color.FgRed,
	WarnLevel:  color.FgYellow,
	InfoLevel:  color.FgGreen,
	DebugLevel: color.FgCyan,
	TraceLevel: color.FgMagenta,
}

// LineFormatter definition
type LineFormatter struct {
	format string
	// field map, parsed from format string.
	// eg: {"level": "{{level}}",}
	fieldMap StringMap

	EnableColor bool
	ColorTheme map[Level]color.Color
}

// NewLineFormatter create new LineFormatter
func NewLineFormatter(format ...string) *LineFormatter  {
	var fmtTpl string
	if len(format) > 0 {
		fmtTpl = format[0]
	} else {
		fmtTpl = SimpleFormat
	}

	return &LineFormatter{
		format: fmtTpl,
		fieldMap: parseFieldMap(fmtTpl),
		ColorTheme: ColorTheme,
	}
}

// FieldMap get export field map
func (f *LineFormatter) FieldMap() StringMap {
	return f.fieldMap
}

// Format an log record
func (f *LineFormatter) Format(r *Record) ([]byte, error) {
	if f.EnableColor {

	}

	tplData := make(StringMap, len(f.fieldMap))

	for field, tplVar := range f.fieldMap {
		switch {
		case field == FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			tplData[tplVar] = r.Time.Format(time.RFC3339)
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

	str := strutil.Replaces(f.format, tplData)

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

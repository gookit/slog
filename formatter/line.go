package formatter

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/goutil/strutil"
	"github.com/gookit/slog"
)

const SimpleFormat = "[{{datetime}}] {{channel}}.{{level}}: {{message}} {{data}} {{extra}}\n"

// LineFormatter definition
type LineFormatter struct {
	format string
	// field map, parsed from format string.
	// eg: {"level": "{{level}}",}
	fieldMap slog.StringMap

	UseColor bool
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
	}
}

func (f *LineFormatter) FieldMap() slog.StringMap {
	return f.fieldMap
}

// Format an log record
func (f *LineFormatter) Format(r *slog.Record) ([]byte, error) {
	tplData := make(slog.StringMap, len(f.fieldMap))

	for field, tplVar := range f.fieldMap {
		switch {
		case field == slog.FieldKeyDatetime:
			if r.Time.IsZero() {
				r.Time = time.Now()
			}

			tplData[tplVar] = r.Time.Format(time.RFC3339)
		case field == slog.FieldKeyLevel:
			tplData[tplVar] = r.LevelName
		case field == slog.FieldKeyChannel:
			tplData[tplVar] = r.Channel
		case field == slog.FieldKeyMessage:
			tplData[tplVar] = r.Message
		case field == slog.FieldKeyData:
			tplData[tplVar] = fmt.Sprint(r.Data)
		case field == slog.FieldKeyExtra:
			tplData[tplVar] = fmt.Sprint(r.Extra)
		default:
			tplData[tplVar] = fmt.Sprint(r.Fields[field])
		}
	}

	// dump.Println(tplData, r.LevelName)

	// TODO ... use buffer
	// var buf *bytes.Buffer

	// strings.NewReplacer().WriteString(buf)

	str := strutil.Replaces(f.format, tplData)

	return []byte(str), nil
}

// parse "{{channel}}" to { "channel": "{{channel}}" }
func parseFieldMap(format string) slog.StringMap {
	rgp := regexp.MustCompile(`{{\w+}}`)

	ss := rgp.FindAllString(format, -1)
	fm := make(slog.StringMap)
	for _, tplVar := range ss {
		field := strings.Trim(tplVar, "{}")
		fm[field] = tplVar
	}

	return fm
}

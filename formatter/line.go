package formatter

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/slog"
)

const SimpleFormat = "[{{datetime}}] {{channel}}.{{level}}: {{message}} {{data}} {{extra}}\n"

// LineFormatter definition
type LineFormatter struct {
	format string
	// eg: {"level": "{{level}}",}
	FieldMap slog.FieldMap

	CliColor bool
}

func NewLineFormatter(format ...string) *LineFormatter  {
	var fmtTpl string
	if len(format) > 0 {
		fmtTpl = format[0]
	} else {
		fmtTpl = SimpleFormat
	}

	return &LineFormatter{
		format: fmtTpl,
		FieldMap: parseFieldMap(fmtTpl),
	}
}

func parseFieldMap(format string) slog.FieldMap {
	rgp := regexp.MustCompile(`{{\w+}}`)

	ss := rgp.FindAllString(format, -1)
	fm := make(slog.FieldMap)
	for _, tplVar := range ss {
		field := strings.Trim(tplVar, "{}")
		fm[field] = tplVar
	}

	return fm
}

// Format an log record
func (f *LineFormatter) Format(r *slog.Record) ([]byte, error) {
	tplData := make(slog.FieldMap, len(f.FieldMap))

	for field, tplVar := range f.FieldMap {
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
		case field == FieldKeyMsg:
			tplData[tplVar] = r.Message
		case field == FieldKeyData:
			tplData[tplVar] = fmt.Sprint(r.Data)
		case field == FieldKeyExtra:
			tplData[tplVar] = fmt.Sprint(r.Extra)
		default:
			tplData[tplVar] = fmt.Sprint(r.Fields[field])
		}
	}

	dump.Println(tplData, r.LevelName)

	// TODO ... use buffer
	// var buf *bytes.Buffer

	// strings.NewReplacer().WriteString(buf)

	str := strutil.Replaces(f.format, tplData)

	return []byte(str), nil
}

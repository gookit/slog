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

const SimpleFormat = "[{{datetime}}] {{channel}}.{{levelName}}: {{message}} {{data}} {{extra}}\n";

// LineFormatter definition
type LineFormatter struct {
	format string
	fieldMap slog.FieldMap

	CliColor bool
}

// FieldMap get field map
func (f *LineFormatter) FieldMap() slog.FieldMap {
	return f.fieldMap
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
		fieldMap: parseFieldMap(fmtTpl),
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

	dump.Println(ss, fm)

	return fm
}

// Format an log record
func (f *LineFormatter) Format(r *slog.Record) ([]byte, error) {
	tplData := make(slog.FieldMap, len(f.fieldMap))

	for field, tplVar := range f.fieldMap {
		switch {
		case field == FieldKeyDatetime:
			tplData[tplVar] = r.Time.Format(time.RFC3339)
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

	// TODO ... use buffer
	// var buf *bytes.Buffer

	// strings.NewReplacer().WriteString(buf)

	str := strutil.Replaces(f.format, tplData)  + "\n"

	return []byte(str), nil
}

package formatter

import (
	"regexp"

	"github.com/gookit/slog"
)

const SimpleFormat = "[{$datetime}] {$channel}.{$level_name}: {$message} {$data} {$extra}\n";

// LineFormatter definition
type LineFormatter struct {
	format string
	fieldMap slog.FieldMap
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
	rgp := regexp.MustCompile(`{$\d+}`)

	ss := rgp.FindAllString(format, -1)
	fm := make(slog.FieldMap)
	for _, field := range ss {
		fm[field] = field
	}

	return fm
}

// Format an log record
func (f *LineFormatter) Format(r *slog.Record) error {
	// TODO ... use buffer
	r.Formatted = []byte(r.Message + "\n")
	return nil
}

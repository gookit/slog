package formatter_test

import (
	"fmt"
	"testing"

	"github.com/gookit/slog/formatter"
)

func TestNewLineFormatter(t *testing.T) {
	lf := formatter.NewLineFormatter()

	fmt.Println(lf.FieldMap())
}

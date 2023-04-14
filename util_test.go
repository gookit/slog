package slog

import (
	"strings"
	"testing"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
)

func revertTemplateString(ss []string) string {
	var sb strings.Builder
	for _, s := range ss {
		// is field
		if s[0] >= 'a' && s[0] <= 'z' {
			sb.WriteString("{{")
			sb.WriteString(s)
			// sb.WriteString("}}")
		} else {
			sb.WriteString(s)
		}
	}

	// sb.WriteByte('\n')
	return sb.String()
}

func TestInner_parseTemplateToFields(t *testing.T) {
	ss := parseTemplateToFields(NamedTemplate)
	str := revertTemplateString(ss)
	// dump.P(ss, str)
	assert.Eq(t, NamedTemplate, str)

	ss = parseTemplateToFields(DefaultTemplate)
	str = revertTemplateString(ss)
	// dump.P(ss, str)
	assert.Eq(t, DefaultTemplate, str)

	testTemplate := "[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}"
	ss = parseTemplateToFields(testTemplate)
	str = revertTemplateString(ss)
	assert.Eq(t, testTemplate, str)
	// dump.P(ss, str)
}

func TestUtil_formatArgsWithSpaces(t *testing.T) {
	// tests for formatArgsWithSpaces
	tests := []struct {
		args []any
		want string
	}{
		{nil, ""},
		{[]any{"a", "b", "c"}, "a b c"},
		{[]any{"a", "b", "c", 1, 2, 3}, "a b c 1 2 3"},
		{[]any{"a", 1, nil}, "a 1 <nil>"},
		{[]any{12, int8(12), int16(12), int32(12), int64(12)}, "12 12 12 12 12"},
		{[]any{uint(12), uint8(12), uint16(12), uint32(12), uint64(12)}, "12 12 12 12 12"},
		{[]any{float32(12.12), 12.12}, "12.12 12.12"},
		{[]any{true, false}, "true false"},
		{[]any{[]byte("abc"), []byte("123")}, "abc 123"},
		{[]any{timex.OneHour}, "1h0m0s"},
		{[]any{errorx.Raw("a error message")}, "a error message"},
		{[]any{[]int{1, 2, 3}}, "[1 2 3]"},
	}

	for _, tt := range tests {
		assert.Eq(t, tt.want, formatArgsWithSpaces(tt.args))
	}

	assert.NotEmpty(t, formatArgsWithSpaces([]any{timex.Now()}))
}

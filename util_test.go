package slog

import (
	"strings"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
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

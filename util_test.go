package slog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, NamedTemplate, str)

	ss = parseTemplateToFields(DefaultTemplate)
	str = revertTemplateString(ss)
	// dump.P(ss, str)
	assert.Equal(t, DefaultTemplate, str)

	testTemplate := "[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}"
	ss = parseTemplateToFields(testTemplate)
	str = revertTemplateString(ss)
	assert.Equal(t, DefaultTemplate, str)
	// dump.P(ss, str)
}

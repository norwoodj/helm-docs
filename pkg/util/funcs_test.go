package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHtmlEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "html",
			input:    `<div>Hello & goodbye</div>`,
			expected: `&lt;div&gt;Hello &amp; goodbye&lt;/div&gt;`,
		},
		{
			name:     "unicode",
			input:    `\u003cdiv\u003eHello \u0026 goodbye\u003c/div\u003e`,
			expected: `&lt;div&gt;Hello &amp; goodbye&lt;/div&gt;`,
		},
		{
			name:     "no escaping needed",
			input:    `hello world`,
			expected: `hello world`,
		},
		{
			name:     "empty string",
			input:    ``,
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := htmlEscape(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

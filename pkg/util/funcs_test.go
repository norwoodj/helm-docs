package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMdTableCell(t *testing.T) {
	tests := map[string]struct{ in, want string }{
		"plain":        {"just text | with pipe", "just text | with pipe"},
		"newlines":     {"line one\nline two", "line one<br>line two"},
		"fence":        {"see:\n```yaml\na: |\n  x\n```\nend", "see:<br><pre lang=\"yaml\">a: \\|<br>  x</pre><br>end"},
		"no lang":      {"```\n<b> & c\n```", "<pre>&lt;b&gt; &amp; c</pre>"},
		"two fences":   {"```sh\nls\n```\n```sh\npwd\n```", "<pre lang=\"sh\">ls</pre><br><pre lang=\"sh\">pwd</pre>"},
		"unterminated": {"```yaml\na: 1", "```yaml<br>a: 1"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, toMdTableCell(tc.in))
		})
	}
}

func TestToMdDefault(t *testing.T) {
	tests := map[string]struct{ in, want string }{
		"nested map":    {"`{\"certPath\":\"\",\"enabled\":true,\"indexUrl\":{\"value\":\"\"}}`", "<pre lang=\"yaml\">certPath: &#34;&#34;<br>enabled: true<br>indexUrl:<br>  value: &#34;&#34;</pre>"},
		"list":          {"`[\"a\",{\"b\":1}]`", "<pre lang=\"yaml\">- a<br>- b: 1</pre>"},
		"empty map":     {"`{}`", "`{}`"},
		"empty list":    {"`[]`", "`[]`"},
		"scalar":        {"`\"x\"`", "`\"x\"`"},
		"not backticks": {"{\"a\":1}", "{\"a\":1}"},
		"big int":       {"`{\"count\":1000000}`", "<pre lang=\"yaml\">count: 1000000</pre>"},
		"escaping":      {"`{\"s\":\"a|<b>\"}`", "<pre lang=\"yaml\">s: a\\|&lt;b&gt;</pre>"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, toMdDefault(tc.in))
		})
	}
}

func TestToHtmlTableCell(t *testing.T) {
	tests := map[string]struct{ in, want string }{
		"newlines": {"line one\nline two", "line one<br>line two"},
		"fence":    {"see:\n```yaml\na: |\n  <x>\n```", "see:<br><pre lang=\"yaml\">a: |<br>  &lt;x&gt;</pre>"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, toHtmlTableCell(tc.in))
		})
	}
}

func TestYamlDefault(t *testing.T) {
	tests := map[string]struct{ in, want string }{
		"map":       {"`{\"a\":{\"b\":1}}`", "a:\n  b: 1"},
		"scalar":    {"`\"x\"`", ""},
		"empty map": {"`{}`", ""},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, yamlDefault(tc.in))
		})
	}
}

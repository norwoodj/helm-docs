package document

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyMarkDownFormat(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "collapses trailing blank lines into a single newline",
			in:   "# Chart\n\n| key | value |\n\n\n",
			want: "# Chart\n\n| key | value |\n",
		},
		{
			name: "adds a trailing newline when missing",
			in:   "# Chart",
			want: "# Chart\n",
		},
		{
			name: "keeps a single trailing newline untouched",
			in:   "# Chart\n",
			want: "# Chart\n",
		},
		{
			name: "still strips trailing space before a newline",
			in:   "trailing space \nnext\n",
			want: "trailing space\nnext\n",
		},
		{
			name: "still collapses three or more blank lines in the body",
			in:   "a\n\n\n\nb\n",
			want: "a\n\nb\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			buf.WriteString(tt.in)
			got := applyMarkDownFormat(buf)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

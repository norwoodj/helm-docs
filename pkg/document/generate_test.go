package document

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyMarkDownFormat_SingleTrailingNewline(t *testing.T) {
	for _, in := range []string{"hello\n\n\n", "hello\n\n", "hello\n", "hello"} {
		buf := bytes.NewBufferString(in)
		result := applyMarkDownFormat(*buf)
		assert.Equal(t, "hello\n", result.String())
	}

}

package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDocumentationTemplate(t *testing.T) {
	tpl, err := getDocumentationTemplate(".", ".", []string{"testdata/nonexistent.md.gotmpl"})

	require.NoError(t, err)
	assert.Equal(t, defaultDocumentationTemplate, tpl)
}

func TestGetDocumentationTemplate_LoadDefaultOnNotFound(t *testing.T) {
	tpl, err := getDocumentationTemplate(".", ".", []string{
		"testdata/README.md.gotmpl",
		"testdata/nonexistent.md.gotmpl",
		"testdata/README2.md.gotmpl",
	})

	const expected = "hello\nhello again\n" + defaultDocumentationTemplate

	require.NoError(t, err)
	assert.Equal(t, expected, tpl)
}

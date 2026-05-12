package document

import (
	"strings"
	"testing"
	"text/template"

	"github.com/norwoodj/helm-docs/pkg/util"
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

func TestValuesTableEscapesPipeInCells(t *testing.T) {
	tmpl := template.New("test").Funcs(util.FuncMap())
	_, err := tmpl.Parse(getValuesTableTemplates())
	require.NoError(t, err)

	data := chartTemplateData{
		Values: []valueRow{
			{
				Key:         "command",
				Type:        "string",
				Default:     "`\"a|b\"`",
				Description: "matches a|b",
			},
		},
	}

	var buf strings.Builder
	require.NoError(t, tmpl.ExecuteTemplate(&buf, "chart.valuesTable", data))

	out := buf.String()
	assert.Contains(t, out, "`\"a\\|b\"`")
	assert.Contains(t, out, "matches a\\|b")
	assert.NotContains(t, out, "`\"a|b\"`")
}

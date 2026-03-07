package document

import (
	"bytes"
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

func TestValuesTable_WithExampleColumn(t *testing.T) {
	tmpl := template.New("test").Funcs(util.FuncMap())
	_, err := tmpl.Parse(getValuesTableTemplates())
	require.NoError(t, err)

	data := chartTemplateData{
		HasExampleColumn: true,
		Values: []valueRow{
			{Key: "bar", Type: "int", Default: "1", Description: "bar", Example: "e.g. 1"},
		},
		Sections: sections{},
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "chart.valuesTable", data)
	require.NoError(t, err)
	out := buf.String()

	assert.Contains(t, out, "| Example |")
	assert.Contains(t, out, "e.g. 1")
	assert.Contains(t, out, "| Key | Type | Default | Description | Example |")
}

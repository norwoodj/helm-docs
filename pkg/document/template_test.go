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

func renderTemplate(t *testing.T, templateName string, templateBody string, data interface{}) string {
	t.Helper()

	tpl, err := template.New("values").Funcs(util.FuncMap()).Parse(templateBody)
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, tpl.ExecuteTemplate(&buf, templateName, data))

	return buf.String()
}

func TestValuesTable_DefaultValue(t *testing.T) {
	tests := []struct {
		name         string
		def          string
		notationType string
		want         string
	}{
		{
			name: "string",
			def:  "`\"bar\"`",
			want: "<pre lang=\"json\">&#34;bar&#34;</pre>",
		},
		{
			name: "int",
			def:  "`42`",
			want: "<pre lang=\"json\">42</pre>",
		},
		{
			name: "float",
			def:  "`3.14`",
			want: "<pre lang=\"json\">3.14</pre>",
		},
		{
			name: "bool",
			def:  "`true`",
			want: "<pre lang=\"json\">true</pre>",
		},
		{
			name: "object",
			def:  "`{\"admin\":true,\"name\":\"root\"}`",
			want: "<pre lang=\"json\">{<br/>  &#34;admin&#34;: true,<br/>  &#34;name&#34;: &#34;root&#34;<br/>}</pre>",
		},
		{
			name: "list",
			def:  "`[\"a\",\"b\",\"c\"]`",
			want: "<pre lang=\"json\">[<br/>  &#34;a&#34;,<br/>  &#34;b&#34;,<br/>  &#34;c&#34;<br/>]</pre>",
		},
		{
			name: "html escape",
			def:  "`\"This <span>HTML tag</span> should be escaped\"`",
			want: "<pre lang=\"json\">&#34;This &lt;span&gt;HTML tag&lt;/span&gt; should be escaped&#34;</pre>",
		},
		{
			name: "unicode",
			def:  "`\"\\u003chtml\\u003e\\u003c/html\\u003e\"`",
			want: "<pre lang=\"json\">&#34;&lt;html&gt;&lt;/html&gt;&#34;</pre>",
		},
		{
			name:         "custom",
			def:          "This is a custom default value with\nan <span>HTML tag</span> that should not be escaped",
			notationType: "custom",
			want:         "<pre lang=\"custom\">This is a custom default value with<br/>an <span>HTML tag</span> that should not be escaped</pre>",
		},
		{
			name:         "tpl",
			def:          "- name: DEBUG\n  value: {{ .Values.global.debug | quote }}",
			notationType: "tpl",
			want:         "<pre lang=\"tpl\">some.key: |<br/>  - name: DEBUG<br/>    value: {{ .Values.global.debug | quote }}</pre>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := []valueRow{{
				Key:          "some.key",
				Default:      tt.def,
				NotationType: tt.notationType,
			}}
			data := chartTemplateData{
				Values:   rows,
				Sections: sections{},
			}
			out := renderTemplate(t, "chart.valuesTable", getValuesTableTemplates(), data)
			assert.Contains(t, out, tt.want)
		})
		t.Run(tt.name+" with sections", func(t *testing.T) {
			rows := []valueRow{{
				Key:          "some.key",
				Default:      tt.def,
				NotationType: tt.notationType,
				Section:      "Some Section",
			}}
			data := chartTemplateData{
				Sections: sections{
					Sections: []section{
						{
							SectionName:  "Some Section",
							SectionItems: rows,
						},
					},
				},
			}
			out := renderTemplate(t, "chart.valuesTable", getValuesTableTemplates(), data)
			assert.Contains(t, out, tt.want)
		})
	}
}

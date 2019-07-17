package document

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/norwoodj/helm-docs/pkg/helm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const defaultDocumentationTemplate = `{{ template "chart.header" . }}
{{ template "chart.description" . }}

{{ template "chart.versionLine" . }}

{{ template "chart.sourceLinkLine" . }}

{{ template "chart.requirementsSection" . }}

{{ template "chart.valuesSection" . }}
`

func getHeaderTemplate() string {
	headerTemplateBuilder := strings.Builder{}
	headerTemplateBuilder.WriteString(`{{ define "chart.header" }}`)
	headerTemplateBuilder.WriteString("{{ .Name }}\n")
	headerTemplateBuilder.WriteString(`{{ repeat (len .Name) "=" }}`)
	headerTemplateBuilder.WriteString("{{ end }}")

	return headerTemplateBuilder.String()
}

func getDescriptionTemplate() string {
	descriptionBuilder := strings.Builder{}
	descriptionBuilder.WriteString(`{{ define "chart.description" }}`)
	descriptionBuilder.WriteString("{{ if .Description }}{{ .Description }}{{ end }}")
	descriptionBuilder.WriteString("{{ end }}")

	return descriptionBuilder.String()
}

func getVersionTemplates() string {
	versionBuilder := strings.Builder{}
	versionBuilder.WriteString(`{{ define "chart.version" }}{{ .Version }}{{ end }}\n`)
	versionBuilder.WriteString(`{{ define "chart.versionLine" }}`)
	versionBuilder.WriteString("Current chart version is `{{ .Version }}`")
	versionBuilder.WriteString("{{ end }}")

	return versionBuilder.String()
}

func getSourceLinkTemplates() string {
	sourceLinkBuilder := strings.Builder{}
	sourceLinkBuilder.WriteString(`{{ define "chart.sourceLink" }}`)
	sourceLinkBuilder.WriteString("{{ .Home }}")
	sourceLinkBuilder.WriteString("{{ end }}\n")

	sourceLinkBuilder.WriteString(`{{ define "chart.sourceLinkLine" }}`)
	sourceLinkBuilder.WriteString("{{ if .Home }}Source code can be found [here]({{ .Home }}){{ end }}")
	sourceLinkBuilder.WriteString("{{ end }}")

	return sourceLinkBuilder.String()
}

func getRequirementsTableTemplates() string {
	requirementsSectionBuilder := strings.Builder{}
	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsHeader" }}## Chart Requirements{{ end }}`)

	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsTable" }}`)
	requirementsSectionBuilder.WriteString("| Repository | Name | Version |\n")
	requirementsSectionBuilder.WriteString("|------------|------|---------|\n")
	requirementsSectionBuilder.WriteString("  {{- range .Dependencies }}")
	requirementsSectionBuilder.WriteString("\n| {{ .Repository }} | {{ .Name }} | {{ .Version }} |")
	requirementsSectionBuilder.WriteString("  {{- end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")

	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsSection" }}`)
	requirementsSectionBuilder.WriteString("{{ if .Dependencies }}")
	requirementsSectionBuilder.WriteString(`{{ template "chart.requirementsHeader" . }}`)
	requirementsSectionBuilder.WriteString("\n\n")
	requirementsSectionBuilder.WriteString(`{{ template "chart.requirementsTable" . }}`)
	requirementsSectionBuilder.WriteString("{{ end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")

	return requirementsSectionBuilder.String()
}

func getValuesTableTemplates() string {
	valuesSectionBuilder := strings.Builder{}
	valuesSectionBuilder.WriteString(`{{ define "chart.valuesHeader" }}## Chart Values{{ end }}`)

	valuesSectionBuilder.WriteString(`{{ define "chart.valuesTable" }}`)
	valuesSectionBuilder.WriteString("| Key | Type | Default | Description |\n")
	valuesSectionBuilder.WriteString("|-----|------|---------|-------------|\n")
	valuesSectionBuilder.WriteString("  {{- range .Values }}")
	valuesSectionBuilder.WriteString("\n| {{ .Key }} | {{ .Type }} | {{ .Default }} | {{ .Description }} |")
	valuesSectionBuilder.WriteString("  {{- end }}")
	valuesSectionBuilder.WriteString("{{ end }}")

	valuesSectionBuilder.WriteString(`{{ define "chart.valuesSection" }}`)
	valuesSectionBuilder.WriteString("{{ if .Values }}")
	valuesSectionBuilder.WriteString(`{{ template "chart.valuesHeader" . }}`)
	valuesSectionBuilder.WriteString("\n\n")
	valuesSectionBuilder.WriteString(`{{ template "chart.valuesTable" . }}`)
	valuesSectionBuilder.WriteString("{{ end }}")
	valuesSectionBuilder.WriteString("{{ end }}")

	return valuesSectionBuilder.String()
}

func getDocumentationTemplate(chartDirectory string) (string, error) {
	templateFile := viper.GetString("template-file")
	templateFileForChart := path.Join(chartDirectory, templateFile)

	if _, err := os.Stat(templateFileForChart); os.IsNotExist(err) {
		log.Debugf("Did not find template file %s for chart %s, using default template", templateFile, chartDirectory)
		return defaultDocumentationTemplate, nil
	}

	log.Debugf("Using template file %s for chart %s", templateFile, chartDirectory)
	templateContents, err := ioutil.ReadFile(templateFileForChart)
	if err != nil {
		return "", err
	}

	return string(templateContents), nil
}

func getDocumentationTemplates(chartDirectory string) ([]string, error) {
	documentationTemplate, err := getDocumentationTemplate(chartDirectory)

	if err != nil {
		log.Errorf("Failed to read documentation template for chart %s: %s", chartDirectory, err)
		return nil, err
	}

	return []string{
		getHeaderTemplate(),
		getDescriptionTemplate(),
		getVersionTemplates(),
		getSourceLinkTemplates(),
		getRequirementsTableTemplates(),
		getValuesTableTemplates(),
		documentationTemplate,
	}, nil
}

func newChartDocumentationTemplate(chartDocumentationInfo helm.ChartDocumentationInfo) (*template.Template, error) {
	documentationTemplate := template.New(chartDocumentationInfo.ChartDirectory)
	documentationTemplate.Funcs(sprig.TxtFuncMap())
	goTemplateList, err := getDocumentationTemplates(chartDocumentationInfo.ChartDirectory)

	if err != nil {
		return nil, err
	}

	for _, t := range goTemplateList {
		_, err := documentationTemplate.Parse(t)

		if err != nil {
			return nil, err
		}
	}

	return documentationTemplate, nil
}

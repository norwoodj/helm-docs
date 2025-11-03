package document

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/norwoodj/helm-docs/pkg/util"

	log "github.com/sirupsen/logrus"

	"github.com/norwoodj/helm-docs/pkg/helm"
)

const defaultDocumentationTemplate = `{{ template "chart.header" . }}
{{ template "chart.deprecationWarning" . }}

{{ template "chart.badgesSection" . }}

{{ template "chart.description" . }}

{{ template "chart.homepageLine" . }}

{{ template "chart.maintainersSection" . }}

{{ template "chart.sourcesSection" . }}

{{ template "chart.requirementsSection" . }}

{{ template "chart.valuesSection" . }}

{{ template "chart.extras" . }}

{{- if not .SkipVersionFooter }}
{{ template "helm-docs.versionFooter" . }}
{{- end }}
`

func getNameTemplate() string {
	nameBuilder := strings.Builder{}
	nameBuilder.WriteString(`{{ define "chart.name" }}`)
	nameBuilder.WriteString("{{ .Name }}")
	nameBuilder.WriteString("{{ end }}")

	return nameBuilder.String()
}

func getHeaderTemplate() string {
	headerTemplateBuilder := strings.Builder{}
	headerTemplateBuilder.WriteString(`{{ define "chart.header" }}`)
	headerTemplateBuilder.WriteString("# {{ .Name }}\n")
	headerTemplateBuilder.WriteString("{{ end }}")

	return headerTemplateBuilder.String()
}

func getDeprecatedTemplate() string {
	deprecatedTemplateBuilder := strings.Builder{}
	deprecatedTemplateBuilder.WriteString(`{{ define "chart.deprecationWarning" }}`)
	deprecatedTemplateBuilder.WriteString("{{ if .Deprecated }}> **:exclamation: This Helm Chart is deprecated!**{{ end }}")
	deprecatedTemplateBuilder.WriteString("{{ end }}")

	return deprecatedTemplateBuilder.String()
}

func getVersionTemplates(badgeStyle string) string {
	versionBuilder := strings.Builder{}
	versionBuilder.WriteString(`{{ define "chart.version" }}{{ .Version }}{{ end }}\n`)
	versionBuilder.WriteString(`{{ define "chart.versionBadge" }}`)
	versionBuilder.WriteString(fmt.Sprintf(`![Version: {{ .Version }}](https://img.shields.io/badge/Version-{{ .Version | replace "-" "--" }}-informational?style=%s) `, badgeStyle))
	versionBuilder.WriteString("{{ end }}")

	return versionBuilder.String()
}

func getTypeTemplate(badgeStyle string) string {
	typeBuilder := strings.Builder{}
	typeBuilder.WriteString(`{{ define "chart.type" }}{{ .Type }}{{ end }}\n`)
	typeBuilder.WriteString(`{{ define "chart.typeBadge" }}`)
	typeBuilder.WriteString(fmt.Sprintf("{{ if .Type }}![Type: {{ .Type }}](https://img.shields.io/badge/Type-{{ .Type }}-informational?style=%s) {{ end }}", badgeStyle))
	typeBuilder.WriteString("{{ end }}")

	return typeBuilder.String()
}

func getAppVersionTemplate(badgeStyle string) string {
	appVersionBuilder := strings.Builder{}
	appVersionBuilder.WriteString(`{{ define "chart.appVersion" }}{{ .AppVersion }}{{ end }}\n`)
	appVersionBuilder.WriteString(`{{ define "chart.appVersionBadge" }}`)
	appVersionBuilder.WriteString(fmt.Sprintf(`{{ if .AppVersion }}![AppVersion: {{ .AppVersion }}](https://img.shields.io/badge/AppVersion-{{ .AppVersion | replace "-" "--" }}-informational?style=%s) {{ end }}`, badgeStyle))
	appVersionBuilder.WriteString("{{ end }}")

	return appVersionBuilder.String()
}

func getBadgesTemplates() string {
	badgeBuilder := strings.Builder{}
	badgeBuilder.WriteString(`{{ define "chart.badgesSection" }}`)
	badgeBuilder.WriteString(`{{ template "chart.versionBadge" . }}{{ template "chart.typeBadge" . }}{{ template "chart.appVersionBadge" . }}`)
	badgeBuilder.WriteString("{{ end }}")

	return badgeBuilder.String()
}

func getDescriptionTemplate() string {
	descriptionBuilder := strings.Builder{}
	descriptionBuilder.WriteString(`{{ define "chart.description" }}`)
	descriptionBuilder.WriteString("{{ if .Description }}{{ .Description }}{{ end }}")
	descriptionBuilder.WriteString("{{ end }}")

	return descriptionBuilder.String()
}

func getHomepageTemplate() string {
	homepageBuilder := strings.Builder{}
	homepageBuilder.WriteString(`{{ define "chart.homepage" }}{{ .Home }}{{ end }}\n`)
	homepageBuilder.WriteString(`{{ define "chart.homepageLine" }}`)
	homepageBuilder.WriteString("{{ if .Home }}**Homepage:** <{{ .Home }}>{{ end }}")
	homepageBuilder.WriteString("{{ end }}")

	return homepageBuilder.String()
}

func getMaintainersTemplate() string {
	maintainerBuilder := strings.Builder{}
	maintainerBuilder.WriteString(`{{ define "chart.maintainersHeader" }}## Maintainers{{ end }}`)

	maintainerBuilder.WriteString(`{{ define "chart.maintainersTable" }}`)
	maintainerBuilder.WriteString("| Name | Email | Url |\n")
	maintainerBuilder.WriteString("| ---- | ------ | --- |\n")
	maintainerBuilder.WriteString("  {{- range .Maintainers }}")
	maintainerBuilder.WriteString("\n| {{ .Name }} | {{ if .Email }}<{{ .Email }}>{{ end }} | {{ if .Url }}<{{ .Url }}>{{ end }} |")
	maintainerBuilder.WriteString("  {{- end }}")
	maintainerBuilder.WriteString("{{ end }}")

	maintainerBuilder.WriteString(`{{ define "chart.maintainersSection" }}`)
	maintainerBuilder.WriteString("{{ if .Maintainers }}")
	maintainerBuilder.WriteString(`{{ template "chart.maintainersHeader" . }}`)
	maintainerBuilder.WriteString("\n\n")
	maintainerBuilder.WriteString(`{{ template "chart.maintainersTable" . }}`)
	maintainerBuilder.WriteString("{{ end }}")
	maintainerBuilder.WriteString("{{ end }}")

	return maintainerBuilder.String()
}

func getSourceLinkTemplates() string {
	sourceLinkBuilder := strings.Builder{}
	sourceLinkBuilder.WriteString(`{{ define "chart.sourcesHeader" }}## Source Code{{ end}}`)

	sourceLinkBuilder.WriteString(`{{ define "chart.sourcesList" }}`)
	sourceLinkBuilder.WriteString("{{- range .Sources }}")
	sourceLinkBuilder.WriteString("\n* <{{ . }}>")
	sourceLinkBuilder.WriteString("{{- end }}")
	sourceLinkBuilder.WriteString("{{ end }}")

	sourceLinkBuilder.WriteString(`{{ define "chart.sourcesSection" }}`)
	sourceLinkBuilder.WriteString("{{ if .Sources }}")
	sourceLinkBuilder.WriteString(`{{ template "chart.sourcesHeader" . }}`)
	sourceLinkBuilder.WriteString("\n")
	sourceLinkBuilder.WriteString(`{{ template "chart.sourcesList" . }}`)
	sourceLinkBuilder.WriteString("{{ end }}")
	sourceLinkBuilder.WriteString("{{ end }}")

	return sourceLinkBuilder.String()
}

func getRequirementsTableTemplates() string {
	requirementsSectionBuilder := strings.Builder{}
	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsHeader" }}## Requirements{{ end }}`)

	requirementsSectionBuilder.WriteString(`{{ define "chart.kubeVersion" }}{{ .KubeVersion }}{{ end }}\n`)
	requirementsSectionBuilder.WriteString(`{{ define "chart.kubeVersionLine" }}`)
	requirementsSectionBuilder.WriteString("{{ if .KubeVersion }}Kubernetes: `{{ .KubeVersion }}`{{ end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")

	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsTable" }}`)
	requirementsSectionBuilder.WriteString("| Repository | Name | Version |\n")
	requirementsSectionBuilder.WriteString("|------------|------|---------|")
	requirementsSectionBuilder.WriteString("  {{- range .Dependencies }}")
	requirementsSectionBuilder.WriteString("    {{- if .Alias }}")
	requirementsSectionBuilder.WriteString("\n| {{ .Repository }} | {{ .Alias }}({{ .Name }}) | {{ .Version }} |")
	requirementsSectionBuilder.WriteString("    {{- else }}")
	requirementsSectionBuilder.WriteString("\n| {{ .Repository }} | {{ .Name }} | {{ .Version }} |")
	requirementsSectionBuilder.WriteString("    {{- end }}")
	requirementsSectionBuilder.WriteString("  {{- end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")

	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsSection" }}`)
	requirementsSectionBuilder.WriteString("{{ if or .Dependencies .KubeVersion }}")
	requirementsSectionBuilder.WriteString(`{{ template "chart.requirementsHeader" . }}`)
	requirementsSectionBuilder.WriteString("\n\n")
	requirementsSectionBuilder.WriteString("{{ if .KubeVersion }}")
	requirementsSectionBuilder.WriteString(`{{ template "chart.kubeVersionLine" . }}`)
	requirementsSectionBuilder.WriteString("\n\n")
	requirementsSectionBuilder.WriteString("{{ end }}")
	requirementsSectionBuilder.WriteString("{{ if .Dependencies }}")
	requirementsSectionBuilder.WriteString(`{{ template "chart.requirementsTable" . }}`)
	requirementsSectionBuilder.WriteString("{{ end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")

	return requirementsSectionBuilder.String()
}

func getValuesTableTemplates() string {
	valuesSectionBuilder := strings.Builder{}
	valuesSectionBuilder.WriteString(`{{ define "chart.valuesHeader" }}## Values{{ end }}`)

	valuesSectionBuilder.WriteString(`{{ define "chart.valuesTable" }}`)
	valuesSectionBuilder.WriteString("{{ if .Sections.Sections }}")
	valuesSectionBuilder.WriteString("{{ range .Sections.Sections }}")
	valuesSectionBuilder.WriteString("\n")
	valuesSectionBuilder.WriteString("\n### {{ .SectionName }}\n")
	valuesSectionBuilder.WriteString("\n")
	valuesSectionBuilder.WriteString("| Key | Type | Default | Description |\n")
	valuesSectionBuilder.WriteString("|-----|------|---------|-------------|\n")
	valuesSectionBuilder.WriteString("  {{- range .SectionItems }}")
	valuesSectionBuilder.WriteString("\n| {{ .Key }} | {{ .Type }} | {{ if .Default }}{{ .Default }}{{ else }}{{ .AutoDefault }}{{ end }} | {{ if .Description }}{{ .Description }}{{ else }}{{ .AutoDescription }}{{ end }} |")
	valuesSectionBuilder.WriteString("  {{- end }}")
	valuesSectionBuilder.WriteString("{{- end }}")
	valuesSectionBuilder.WriteString("{{ if .Sections.DefaultSection.SectionItems}}")
	valuesSectionBuilder.WriteString("\n")
	valuesSectionBuilder.WriteString("\n### {{ .Sections.DefaultSection.SectionName }}\n")
	valuesSectionBuilder.WriteString("\n")
	valuesSectionBuilder.WriteString("| Key | Type | Default | Description |\n")
	valuesSectionBuilder.WriteString("|-----|------|---------|-------------|\n")
	valuesSectionBuilder.WriteString("  {{- range .Sections.DefaultSection.SectionItems }}")
	valuesSectionBuilder.WriteString("\n| {{ .Key }} | {{ .Type }} | {{ if .Default }}{{ .Default }}{{ else }}{{ .AutoDefault }}{{ end }} | {{ if .Description }}{{ .Description }}{{ else }}{{ .AutoDescription }}{{ end }} |")
	valuesSectionBuilder.WriteString("  {{- end }}")
	valuesSectionBuilder.WriteString("{{ end }}")
	valuesSectionBuilder.WriteString("{{ else }}")
	valuesSectionBuilder.WriteString("| Key | Type | Default | Description |\n")
	valuesSectionBuilder.WriteString("|-----|------|---------|-------------|\n")
	valuesSectionBuilder.WriteString("  {{- range .Values }}")
	valuesSectionBuilder.WriteString("\n| {{ .Key }} | {{ .Type }} | {{ if .Default }}{{ .Default }}{{ else }}{{ .AutoDefault }}{{ end }} | {{ if .Description }}{{ .Description }}{{ else }}{{ .AutoDescription }}{{ end }} |")
	valuesSectionBuilder.WriteString("  {{- end }}")
	valuesSectionBuilder.WriteString("{{ end }}")
	valuesSectionBuilder.WriteString("{{ end }}")

	valuesSectionBuilder.WriteString(`{{ define "chart.valuesSection" }}`)
	valuesSectionBuilder.WriteString("{{ if .Values }}")
	valuesSectionBuilder.WriteString(`{{ template "chart.valuesHeader" . }}`)
	valuesSectionBuilder.WriteString("\n\n")
	valuesSectionBuilder.WriteString(`{{ template "chart.valuesTable" . }}`)
	valuesSectionBuilder.WriteString("{{ end }}")
	valuesSectionBuilder.WriteString("{{ end }}")

	// For HTML tables
	valuesSectionBuilder.WriteString(`
{{ define "chart.valueDefaultColumnRender" }}
{{- $defaultValue := (default .Default .AutoDefault)  -}}
{{- $notationType := .NotationType }}
{{- if (and (hasPrefix "` + "`" + `" $defaultValue) (hasSuffix "` + "`" + `" $defaultValue) ) -}}
{{- $defaultValue = (toPrettyJson (fromJson (trimAll "` + "`" + `" (default .Default .AutoDefault) ) ) ) -}}
{{- $notationType = "json" }}
{{- end -}}
<pre lang="{{ $notationType }}">
{{- if (eq $notationType "tpl" ) }}
{{ .Key }}: |
{{- $defaultValue | nindent 2 }}
{{- else }}
{{ $defaultValue }}
{{- end }}
</pre>
{{ end }}

{{ define "chart.valuesTableHtml" }}
{{ if .Sections.Sections }}
{{- range .Sections.Sections }}
<h3>{{- .SectionName }}</h3>
<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
	{{- range .SectionItems }}
		<tr>
			<td>{{ .Key }}</td>
			<td>{{ .Type }}</td>
			<td>{{ template "chart.valueDefaultColumnRender" . }}</td>
			<td>{{ if .Description }}{{ .Description }}{{ else }}{{ .AutoDescription }}{{ end }}</td>
		</tr>
	{{- end }}
	</tbody>
</table>
{{- end }}
{{ if .Sections.DefaultSection.SectionItems }}
<h3>{{- .Sections.DefaultSection.SectionName }}</h3>
<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
	{{- range .Sections.DefaultSection.SectionItems }}
	<tr>
		<td>{{ .Key }}</td>
		<td>{{ .Type }}</td>
		<td>{{ template "chart.valueDefaultColumnRender" . }}</td>
		<td>{{ if .Description }}{{ .Description }}{{ else }}{{ .AutoDescription }}{{ end }}</td>
	</tr>
	{{- end }}
	</tbody>
</table>
{{ end }}
{{ else }}
<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
	{{- range .Values }}
		<tr>
			<td>{{ .Key }}</td>
			<td>{{ .Type }}</td>
			<td>{{ template "chart.valueDefaultColumnRender" . }}</td>
			<td>{{ if .Description }}{{ .Description }}{{ else }}{{ .AutoDescription }}{{ end }}</td>
		</tr>
	{{- end }}
	</tbody>
</table>
{{ end }}
{{ end }}

{{ define "chart.valuesSectionHtml" }}
{{ if .Sections }}
{{ template "chart.valuesHeader" . }}
{{ template "chart.valuesTableHtml" . }}
{{ end }}
{{ end }}
		`)

	return valuesSectionBuilder.String()
}

func getExtraTemplates() string {
	extrasBuilder := strings.Builder{}
	extrasBuilder.WriteString(`{{ define "chart.extras" }}`)
	extrasBuilder.WriteString("{{ range .Extras }}\n")
	extrasBuilder.WriteString("### {{ .Title }}\n")
	extrasBuilder.WriteString("{{ .MarkDown }}\n\n")
	extrasBuilder.WriteString("{{ end }}")
	extrasBuilder.WriteString("{{ end }}")

	return extrasBuilder.String()
}

func getHelmDocsVersionTemplates() string {
	versionSectionBuilder := strings.Builder{}
	versionSectionBuilder.WriteString(`{{ define "helm-docs.version" }}{{ if .HelmDocsVersion }}{{ .HelmDocsVersion }}{{ end }}{{ end }}`)
	versionSectionBuilder.WriteString(`{{ define "helm-docs.versionFooter" }}`)
	versionSectionBuilder.WriteString("{{ if .HelmDocsVersion }}\n")
	versionSectionBuilder.WriteString("----------------------------------------------\n")
	versionSectionBuilder.WriteString("Autogenerated from chart metadata using [helm-docs v{{ .HelmDocsVersion }}](https://github.com/norwoodj/helm-docs/releases/v{{ .HelmDocsVersion }})")
	versionSectionBuilder.WriteString("{{ end }}")
	versionSectionBuilder.WriteString("{{ end }}")

	return versionSectionBuilder.String()
}

func getDocumentationTemplate(chartDirectory string, chartSearchRoot string, templateFiles []string) (string, error) {
	templateFilesForChart := make([]string, 0)

	var templateNotFound bool

	for _, templateFile := range templateFiles {
		var fullTemplatePath string

		if util.IsRelativePath(templateFile) {
			fullTemplatePath = filepath.Join(chartSearchRoot, templateFile)
		} else if util.IsBaseFilename(templateFile) {
			fullTemplatePath = filepath.Join(chartDirectory, templateFile)
		} else {
			fullTemplatePath = templateFile
		}

		if _, err := os.Stat(fullTemplatePath); os.IsNotExist(err) {
			log.Debugf("Did not find template file %s for chart %s, using default template", templateFile, chartDirectory)

			templateNotFound = true
			continue
		}

		templateFilesForChart = append(templateFilesForChart, fullTemplatePath)
	}

	log.Debugf("Using template files %s for chart %s", templateFiles, chartDirectory)
	allTemplateContents := make([]byte, 0)
	for _, templateFileForChart := range templateFilesForChart {
		templateContents, err := ioutil.ReadFile(templateFileForChart)
		if err != nil {
			return "", err
		}
		allTemplateContents = append(allTemplateContents, templateContents...)
	}

	if templateNotFound {
		allTemplateContents = append(allTemplateContents, []byte(defaultDocumentationTemplate)...)
	}

	return string(allTemplateContents), nil
}

func getDocumentationTemplates(chartDirectory string, chartSearchRoot string, templateFiles []string, badgeStyle string) ([]string, error) {
	documentationTemplate, err := getDocumentationTemplate(chartDirectory, chartSearchRoot, templateFiles)
	if err != nil {
		log.Errorf("Failed to read documentation template for chart %s: %s", chartDirectory, err)
		return nil, err
	}

	return []string{
		getNameTemplate(),
		getHeaderTemplate(),
		getDeprecatedTemplate(),
		getAppVersionTemplate(badgeStyle),
		getBadgesTemplates(),
		getDescriptionTemplate(),
		getVersionTemplates(badgeStyle),
		getTypeTemplate(badgeStyle),
		getSourceLinkTemplates(),
		getRequirementsTableTemplates(),
		getValuesTableTemplates(),
		getHomepageTemplate(),
		getMaintainersTemplate(),
		getHelmDocsVersionTemplates(),
		getExtraTemplates(),
		documentationTemplate,
	}, nil
}

func newChartDocumentationTemplate(chartDocumentationInfo helm.ChartDocumentationInfo, chartSearchRoot string, templateFiles []string, badgeStyle string) (*template.Template, error) {
	documentationTemplate := template.New(chartDocumentationInfo.ChartDirectory)
	documentationTemplate.Funcs(util.FuncMap())
	goTemplateList, err := getDocumentationTemplates(chartDocumentationInfo.ChartDirectory, chartSearchRoot, templateFiles, badgeStyle)
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

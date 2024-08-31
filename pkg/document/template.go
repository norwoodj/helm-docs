package document

import (
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/norwoodj/helm-docs/pkg/util"

	log "github.com/sirupsen/logrus"

	"github.com/norwoodj/helm-docs/pkg/helm"
)

//go:embed templates
var embedTemplates embed.FS

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
		allTemplateContents = append(allTemplateContents, []byte(`{{ template "chart.defaultTemplate" }}`)...)
	}

	return string(allTemplateContents), nil
}

func getDocumentationTemplates(chartDirectory string, chartSearchRoot string, templateFiles []string) (string, error) {
	documentationTemplate, err := getDocumentationTemplate(chartDirectory, chartSearchRoot, templateFiles)

	if err != nil {
		log.Errorf("Failed to read documentation template for chart %s: %s", chartDirectory, err)
		return "", err
	}
	return documentationTemplate, nil
}

func newChartDocumentationTemplate(chartDocumentationInfo helm.ChartDocumentationInfo, chartSearchRoot string, templateFiles []string) (*template.Template, error) {
	documentationTemplate := template.New(chartDocumentationInfo.ChartDirectory)
	documentationTemplate.Funcs(util.FuncMap())
	docsTemplate, err := getDocumentationTemplates(chartDocumentationInfo.ChartDirectory, chartSearchRoot, templateFiles)

	if err != nil {
		return nil, err
	}

	if _, err := documentationTemplate.ParseFS(embedTemplates, "templates/*.tmpl"); err != nil {
		return nil, err

	}
	if _, err := documentationTemplate.Parse(docsTemplate); err != nil {
		return nil, err
	}

	return documentationTemplate, nil
}

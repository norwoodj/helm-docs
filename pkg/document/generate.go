package document

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/norwoodj/helm-docs/pkg/helm"
	log "github.com/sirupsen/logrus"
)

func getOutputFile(chartDirectory string, dryRun bool) (*os.File, error) {
	if dryRun {
		return os.Stdout, nil
	}

	f, err := os.Create(fmt.Sprintf("%s/README.md", chartDirectory))

	if err != nil {
		return nil, err
	}

	return f, err
}

func PrintDocumentation(chartDocumentationInfo helm.ChartDocumentationInfo, dryRun bool) {
	log.Infof("Generating README Documentation for chart %s", chartDocumentationInfo.ChartDirectory)

	outputFile, err := getOutputFile(chartDocumentationInfo.ChartDirectory, dryRun)
	if err != nil {
		log.Warnf("Could not open chart README file %s, skipping chart", filepath.Join(chartDocumentationInfo.ChartDirectory, "README.md"))
		return
	}

	if !dryRun {
		defer outputFile.Close()
	}

	chartDocumentationTemplate, err := newChartDocumentationTemplate(chartDocumentationInfo)
	if err != nil {
		log.Warnf("Error generating templates for chart %s: %s", chartDocumentationInfo.ChartDirectory, err)
		return
	}

	chartTemplateDataObject := getChartTemplateData(chartDocumentationInfo)
	err = chartDocumentationTemplate.Execute(outputFile, chartTemplateDataObject)

	if err != nil {
		log.Warnf("Error generating documentation for chart %s: %s", chartDocumentationInfo.ChartDirectory, err)
	}
}

package document

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/norwoodj/helm-docs/pkg/helm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getOutputFile(chartDirectory string, dryRun bool) (*os.File, error) {
	if dryRun {
		return os.Stdout, nil
	}

	outputFile := viper.GetString("output-file")
	f, err := os.Create(fmt.Sprintf("%s/%s", chartDirectory, outputFile))

	if err != nil {
		return nil, err
	}

	return f, err
}

func PrintDocumentation(chartDocumentationInfo helm.ChartDocumentationInfo, dryRun bool) {
	log.Infof("Generating README Documentation for chart %s", chartDocumentationInfo.ChartDirectory)

	chartDocumentationTemplate, err := newChartDocumentationTemplate(chartDocumentationInfo)
	if err != nil {
		log.Warnf("Error generating gotemplates for chart %s: %s", chartDocumentationInfo.ChartDirectory, err)
		return
	}

	chartTemplateDataObject, err := getChartTemplateData(chartDocumentationInfo)
	if err != nil {
		log.Warnf("Error generating template data for chart %s: %s", chartDocumentationInfo.ChartDirectory, err)
		return
	}

	outputFile, err := getOutputFile(chartDocumentationInfo.ChartDirectory, dryRun)
	if err != nil {
		log.Warnf("Could not open chart README file %s, skipping chart", filepath.Join(chartDocumentationInfo.ChartDirectory, "README.md"))
		return
	}

	if !dryRun {
		defer outputFile.Close()
	}

	err = chartDocumentationTemplate.Execute(outputFile, chartTemplateDataObject)
	if err != nil {
		log.Warnf("Error generating documentation for chart %s: %s", chartDocumentationInfo.ChartDirectory, err)
	}
}

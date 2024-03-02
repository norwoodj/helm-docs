package document

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"

	"github.com/norwoodj/helm-docs/pkg/helm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getOutputFile(chartDirectory string, dryRun bool) (*os.File, error) {
	if dryRun {
		return os.Stdout, nil
	}

	outputFile := viper.GetString("output-file")
	f, err := os.Create(filepath.Join(chartDirectory, outputFile))

	if err != nil {
		return nil, err
	}

	return f, err
}

func PrintDocumentation(chartDocumentationInfo helm.ChartDocumentationInfo, chartSearchRoot string, templateFiles []string, dryRun bool, helmDocsVersion string, badgeStyle string, dependencyValues []DependencyValues, skipVersionFooter bool) {
	log.Infof("Generating README Documentation for chart %s", chartDocumentationInfo.ChartDirectory)

	chartDocumentationTemplate, err := newChartDocumentationTemplate(
		chartDocumentationInfo,
		chartSearchRoot,
		templateFiles,
		badgeStyle,
	)

	if err != nil {
		log.Warnf("Error generating gotemplates for chart %s: %s", chartDocumentationInfo.ChartDirectory, err)
		return
	}

	chartTemplateDataObject, err := getChartTemplateData(chartDocumentationInfo, helmDocsVersion, dependencyValues, skipVersionFooter)
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

	var output bytes.Buffer
	err = chartDocumentationTemplate.Execute(&output, chartTemplateDataObject)
	if err != nil {
		log.Warnf("Error generating documentation for chart %s: %s", chartDocumentationInfo.ChartDirectory, err)
	}

	output = applyMarkDownFormat(output)
	_, err = output.WriteTo(outputFile)
	if err != nil {
		log.Warnf("Error generating documentation file for chart %s: %s", chartDocumentationInfo.ChartDirectory, err)
	}
}

func applyMarkDownFormat(output bytes.Buffer) bytes.Buffer {
	outputString := output.String()
	re := regexp.MustCompile(` \n`)
	outputString = re.ReplaceAllString(outputString, "\n")

	re = regexp.MustCompile(`\n{3,}`)
	outputString = re.ReplaceAllString(outputString, "\n\n")

	output.Reset()
	output.WriteString(outputString)
	return output
}

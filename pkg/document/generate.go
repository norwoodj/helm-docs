package document

import (
	"bytes"
	"fmt"
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

func PrintDocumentation(chartDocumentationInfo helm.ChartDocumentationInfo, chartSearchRoot string, templateFiles []string, dryRun bool, helmDocsVersion string, badgeStyle string, dependencyValues []DependencyValues, skipVersionFooter bool) (err error) {
	log.Infof("Generating README Documentation for chart %s", chartDocumentationInfo.ChartDirectory)

	chartDocumentationTemplate, err := newChartDocumentationTemplate(
		chartDocumentationInfo,
		chartSearchRoot,
		templateFiles,
		badgeStyle,
	)

	if err != nil {
		return fmt.Errorf("error generating gotemplates for chart %s: %w", chartDocumentationInfo.ChartDirectory, err)
	}

	chartTemplateDataObject, err := getChartTemplateData(chartDocumentationInfo, helmDocsVersion, dependencyValues, skipVersionFooter)
	if err != nil {
		return fmt.Errorf("error generating template data for chart %s: %w", chartDocumentationInfo.ChartDirectory, err)
	}

	outputFile, err := getOutputFile(chartDocumentationInfo.ChartDirectory, dryRun)
	if err != nil {
		return fmt.Errorf("could not open chart README file %s: %w", filepath.Join(chartDocumentationInfo.ChartDirectory, viper.GetString("output-file")), err)
	}
	if !dryRun {
		defer func() {
			if closeErr := outputFile.Close(); closeErr != nil && err == nil {
				err = fmt.Errorf("error closing documentation file for chart %s: %w", chartDocumentationInfo.ChartDirectory, closeErr)
			}
		}()
	}

	var output bytes.Buffer
	err = chartDocumentationTemplate.Execute(&output, chartTemplateDataObject)
	if err != nil {
		return fmt.Errorf("error generating documentation for chart %s: %w", chartDocumentationInfo.ChartDirectory, err)
	}

	output = applyMarkDownFormat(output)
	_, err = output.WriteTo(outputFile)
	if err != nil {
		return fmt.Errorf("error generating documentation file for chart %s: %w", chartDocumentationInfo.ChartDirectory, err)
	}

	return nil
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

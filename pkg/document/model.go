package document

import (
	"fmt"
	"strconv"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"gopkg.in/yaml.v3"
)

type valueRow struct {
	Key             string
	Type            string
	AutoDefault     string
	Default         string
	AutoDescription string
	Description     string
}

type chartTemplateData struct {
	helm.ChartDocumentationInfo
	HelmDocsVersion string
	Values []valueRow
}

func getChartTemplateData(chartDocumentationInfo helm.ChartDocumentationInfo, helmDocsVersion string) (chartTemplateData, error) {
	// handle empty values file case
	if chartDocumentationInfo.ChartValues.Kind == 0 {
		return chartTemplateData{
			ChartDocumentationInfo: chartDocumentationInfo,
			HelmDocsVersion: helmDocsVersion,
			Values:                 make([]valueRow, 0),
		}, nil
	}

	if chartDocumentationInfo.ChartValues.Kind != yaml.DocumentNode {
		return chartTemplateData{}, fmt.Errorf("invalid node kind supplied: %d", chartDocumentationInfo.ChartValues.Kind)
	}
	if chartDocumentationInfo.ChartValues.Content[0].Kind != yaml.MappingNode {
		return chartTemplateData{}, fmt.Errorf("values file must resolve to a map, not %s", strconv.Itoa(int(chartDocumentationInfo.ChartValues.Kind)))
	}

	valuesTableRows, err := createValueRowsFromObject(
		"",
		nil,
		chartDocumentationInfo.ChartValues.Content[0],
		chartDocumentationInfo.ChartValuesDescriptions,
		true,
	)

	if err != nil {
		return chartTemplateData{}, err
	}

	return chartTemplateData{
		ChartDocumentationInfo: chartDocumentationInfo,
		HelmDocsVersion: helmDocsVersion,
		Values:                 valuesTableRows,
	}, nil
}

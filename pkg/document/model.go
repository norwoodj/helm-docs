package document

import (
	"fmt"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type valueRow struct {
	Key             string
	Type            string
	AutoDefault     string
	Default         string
	AutoDescription string
	Description     string
	Column          int
	LineNumber      int
}

type chartTemplateData struct {
	helm.ChartDocumentationInfo
	HelmDocsVersion string
	Values          []valueRow
}

func getChartTemplateData(chartDocumentationInfo helm.ChartDocumentationInfo, helmDocsVersion string) (chartTemplateData, error) {
	// handle empty values file case
	if chartDocumentationInfo.ChartValues.Kind == 0 {
		return chartTemplateData{
			ChartDocumentationInfo: chartDocumentationInfo,
			HelmDocsVersion:        helmDocsVersion,
			Values:                 make([]valueRow, 0),
		}, nil
	}

	if chartDocumentationInfo.ChartValues.Kind != yaml.DocumentNode {
		return chartTemplateData{}, fmt.Errorf("invalid node kind supplied: %d", chartDocumentationInfo.ChartValues.Kind)
	}
	if chartDocumentationInfo.ChartValues.Content[0].Kind != yaml.MappingNode {
		return chartTemplateData{}, fmt.Errorf("values file must resolve to a map, not %s", strconv.Itoa(int(chartDocumentationInfo.ChartValues.Kind)))
	}

	valuesTableRows, err := createValueRowsFromField(
		"",
		nil,
		chartDocumentationInfo.ChartValues.Content[0],
		chartDocumentationInfo.ChartValuesDescriptions,
		true,
	)

	sortOrder := viper.GetString("sort-values-order")
	if sortOrder == FileSortOrder {
		sort.Slice(valuesTableRows, func(i, j int) bool {
			if valuesTableRows[i].LineNumber < valuesTableRows[j].LineNumber {
				return true
			} else if valuesTableRows[i].Column < valuesTableRows[j].Column {
				return true
			}

			return false
		})
	} else if sortOrder == AlphaNumSortOrder {
		sort.Slice(valuesTableRows, func(i, j int) bool {
			return valuesTableRows[i].Key < valuesTableRows[j].Key
		})
	} else {
		log.Warnf("Invalid sort order provided %s, defaulting to %s", sortOrder, AlphaNumSortOrder)
		sort.Slice(valuesTableRows, func(i, j int) bool {
			return valuesTableRows[i].Key < valuesTableRows[j].Key
		})
	}

	if err != nil {
		return chartTemplateData{}, err
	}

	return chartTemplateData{
		ChartDocumentationInfo: chartDocumentationInfo,
		HelmDocsVersion:        helmDocsVersion,
		Values:                 valuesTableRows,
	}, nil
}

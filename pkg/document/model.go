package document

import (
	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/spf13/viper"
)

type valueRow struct {
	Key         string
	Type        string
	Default     string
	Description string
}

type chartTemplateData struct {
	helm.ChartDocumentationInfo
	Values []valueRow
}

func getChartTemplateData(chartDocumentationInfo helm.ChartDocumentationInfo) (chartTemplateData, error) {
	documentLeafNodes := !viper.GetBool("omit-blanks")
	blankContainerDefaults := viper.GetBool("blank-container-defaults")

	valuesTableRows, err := createValueRowsFromObject(
		"",
		chartDocumentationInfo.ChartValues,
		chartDocumentationInfo.ChartValuesDescriptions,
		documentLeafNodes,
		blankContainerDefaults,
	)

	if err != nil {
		return chartTemplateData{}, err
	}

	return chartTemplateData{
		ChartDocumentationInfo: chartDocumentationInfo,
		Values:                 valuesTableRows,
	}, nil
}

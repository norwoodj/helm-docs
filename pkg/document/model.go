package document

import (
	"github.com/norwoodj/helm-docs/pkg/helm"
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

func getChartTemplateData(chartDocumentationInfo helm.ChartDocumentationInfo) chartTemplateData {
	valuesTableRows := createValueRows("", chartDocumentationInfo.ChartValues, chartDocumentationInfo.ChartValuesDescriptions)

	return chartTemplateData{
		ChartDocumentationInfo: chartDocumentationInfo,
		Values:                 valuesTableRows,
	}
}

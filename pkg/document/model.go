package document

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/norwoodj/helm-docs/pkg/helm"
)

type valueRow struct {
	Key             string
	Type            string
	NotationType    string
	AutoDefault     string
	Default         string
	AutoDescription string
	Description     string
	Column          int
	LineNumber      int
	Dependency      string
	IsGlobal        bool
}

type chartValuesSubsection struct {
	Title      string
	TitleLevel string
	Values     []valueRow
}

type chartTemplateData struct {
	helm.ChartDocumentationInfo
	HelmDocsVersion   string
	Values            []valueRow
	ValuesSubsections []chartValuesSubsection
}

func sortValueRows(valueRows []valueRow) {
	sortOrder := viper.GetString("sort-values-order")

	if sortOrder != FileSortOrder && sortOrder != AlphaNumSortOrder {
		log.Warnf("Invalid sort order provided %s, defaulting to %s", sortOrder, AlphaNumSortOrder)
		sortOrder = AlphaNumSortOrder
	}

	sort.Slice(valueRows, func(i, j int) bool {
		// Globals sort above non-globals.
		if valueRows[i].IsGlobal != valueRows[j].IsGlobal {
			return valueRows[i].IsGlobal
		}

		// Group by dependency for non-globals.
		if !valueRows[i].IsGlobal && !valueRows[j].IsGlobal {
			// Values for the main chart sort above values for dependencies.
			if (valueRows[i].Dependency == "") != (valueRows[j].Dependency == "") {
				return valueRows[i].Dependency == ""
			}

			// Group dependency values together.
			if valueRows[i].Dependency != valueRows[j].Dependency {
				return valueRows[i].Dependency < valueRows[j].Dependency
			}
		}

		// Sort the remaining values within the same section using the configured sort order.
		switch sortOrder {
		case FileSortOrder:
			if valueRows[i].LineNumber == valueRows[j].LineNumber {
				return valueRows[i].Column < valueRows[j].Column
			}
			return valueRows[i].LineNumber < valueRows[i].LineNumber
		case AlphaNumSortOrder:
			return valueRows[i].Key < valueRows[j].Key
		default:
			panic("cannot get here")
		}
	})
}

func getUnsortedValueRows(document *yaml.Node, descriptions map[string]helm.ChartValueDescription) ([]valueRow, error) {
	// Handle empty values file case.
	if document.Kind == 0 {
		return nil, nil
	}

	if document.Kind != yaml.DocumentNode {
		return nil, fmt.Errorf("invalid node kind supplied: %d", document.Kind)
	}

	if document.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("values file must resolve to a map (was %d)", document.Content[0].Kind)
	}

	return createValueRowsFromField("", nil, document.Content[0], descriptions, true)
}

func getChartTemplateData(info helm.ChartDocumentationInfo, helmDocsVersion string, dependencyValues []DependencyValues) (chartTemplateData, error) {
	dependencySectionsEnabled := viper.GetBool("dependency-values-as-sections")
	var valuesSubsections []chartValuesSubsection
	valuesTableRows, err := getUnsortedValueRows(info.ChartValues, info.ChartValuesDescriptions)
	if err != nil {
		return chartTemplateData{}, err
	}

	if len(dependencyValues) > 0 {
		seenGlobalKeys := make(map[string]bool)
		for i, row := range valuesTableRows {
			if strings.HasPrefix(row.Key, "global.") {
				valuesTableRows[i].IsGlobal = true
				seenGlobalKeys[row.Key] = true
			}
		}

		for _, dep := range dependencyValues {
			depValuesTableRows, err := getUnsortedValueRows(dep.ChartValues, dep.ChartValuesDescriptions)
			if err != nil {
				return chartTemplateData{}, err
			}

			for i := len(depValuesTableRows) -1; i >= 0; i-- {
				var row = &depValuesTableRows[i]
				if row.Key == "global" || strings.HasPrefix(row.Key, "global.") {
					if seenGlobalKeys[row.Key] {
						depValuesTableRows = append(depValuesTableRows[:i], depValuesTableRows[i+1:]...)
						continue
					}
					row.IsGlobal = true
					seenGlobalKeys[row.Key] = true
				} else {
					row.Key = dep.Prefix + "." + row.Key
				}

				row.Dependency = dep.Prefix
			}
			if dependencySectionsEnabled && len(depValuesTableRows) > 0 {
				valuesSubsections = append(valuesSubsections, chartValuesSubsection{
					Title: dep.Prefix,
					TitleLevel : "###",
					Values: depValuesTableRows,
				})
			} else {
				valuesTableRows = append(valuesTableRows, depValuesTableRows...)
			}

		}
	}

	sortValueRows(valuesTableRows)
	if dependencySectionsEnabled && len(valuesTableRows) > 0 {
		var section = chartValuesSubsection{
			Title: "",
			TitleLevel: "",
			Values: valuesTableRows,
		}
		valuesSubsections = append( []chartValuesSubsection{section}, valuesSubsections...)
	}

	return chartTemplateData{
		ChartDocumentationInfo: info,
		HelmDocsVersion:        helmDocsVersion,
		Values:                 valuesTableRows,
		ValuesSubsections:      valuesSubsections,
	}, nil
}

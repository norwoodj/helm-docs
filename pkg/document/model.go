package document

import (
	"fmt"
	"slices"
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
	Section         string
	Column          int
	LineNumber      int
	Dependency      string
	IsGlobal        bool
}

type chartTemplateData struct {
	helm.ChartDocumentationInfo
	HelmDocsVersion   string
	Values            []valueRow
	Sections          sections
	Extras            []extra
	Files             files
	SkipVersionFooter bool
}

type sections struct {
	DefaultSection section
	Sections       []section
}

type section struct {
	SectionName  string
	SectionItems []valueRow
}

type extra struct {
	Title    string
	MarkDown string
}

func sortValueRowsByOrder(valueRows []valueRow, sortOrder string) {
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
			return valueRows[i].LineNumber < valueRows[j].LineNumber
		case AlphaNumSortOrder:
			return valueRows[i].Key < valueRows[j].Key
		default:
			panic("cannot get here")
		}
	})
}

func sortValueRows(valueRows []valueRow) {
	sortOrder := viper.GetString("sort-values-order")

	if sortOrder != FileSortOrder && sortOrder != AlphaNumSortOrder {
		log.Warnf("Invalid sort order provided %s, defaulting to %s", sortOrder, AlphaNumSortOrder)
		sortOrder = AlphaNumSortOrder
	}

	sortValueRowsByOrder(valueRows, sortOrder)
}

func sortSectionedValueRows(sectionedValueRows sections) {
	sortOrder := viper.GetString("sort-values-order")

	if sortOrder != FileSortOrder && sortOrder != AlphaNumSortOrder {
		log.Warnf("Invalid sort order provided %s, defaulting to %s", sortOrder, AlphaNumSortOrder)
		sortOrder = AlphaNumSortOrder
	}

	sortValueRowsByOrder(sectionedValueRows.DefaultSection.SectionItems, sortOrder)

	for _, section := range sectionedValueRows.Sections {
		sortValueRowsByOrder(section.SectionItems, sortOrder)
	}
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

func getSectionedValueRows(valueRows []valueRow) sections {
	var valueRowsSectionSorted sections
	valueRowsSectionSorted.DefaultSection = section{
		SectionName:  "Other Values",
		SectionItems: []valueRow{},
	}

	for _, row := range valueRows {
		if row.Section == "" {
			valueRowsSectionSorted.DefaultSection.SectionItems = append(valueRowsSectionSorted.DefaultSection.SectionItems, row)
			continue
		}

		containsSection := false
		for i, section := range valueRowsSectionSorted.Sections {
			if section.SectionName == row.Section {
				containsSection = true
				valueRowsSectionSorted.Sections[i].SectionItems = append(valueRowsSectionSorted.Sections[i].SectionItems, row)
				break
			}
		}

		if !containsSection {
			valueRowsSectionSorted.Sections = append(valueRowsSectionSorted.Sections, section{
				SectionName:  row.Section,
				SectionItems: []valueRow{row},
			})
		}
	}

	return valueRowsSectionSorted
}

func getChartTemplateData(info helm.ChartDocumentationInfo, helmDocsVersion string, dependencyValues []DependencyValues, skipVersionFooter bool) (chartTemplateData, error) {
	valuesTableRows, err := getUnsortedValueRows(info.ChartValues, info.ChartValuesDescriptions)
	if err != nil {
		return chartTemplateData{}, err
	}

	if viper.GetBool("ignore-non-descriptions") {
		valuesTableRows = removeRowsWithoutDescription(valuesTableRows)
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

			for _, row := range depValuesTableRows {
				if row.Key == "global" || strings.HasPrefix(row.Key, "global.") {
					if seenGlobalKeys[row.Key] {
						continue
					}
					row.IsGlobal = true
					seenGlobalKeys[row.Key] = true
				} else {
					row.Key = dep.Prefix + "." + row.Key
				}

				row.Dependency = dep.Prefix
				valuesTableRows = append(valuesTableRows, row)
			}
		}
	}

	sortValueRows(valuesTableRows)
	valueRowsSectionSorted := getSectionedValueRows(valuesTableRows)
	sortSectionedValueRows(valueRowsSectionSorted)

	files, err := getFiles(info.ChartDirectory)
	if err != nil {
		return chartTemplateData{}, err
	}

	var extras []extra
	for title, markdown := range info.Extras {
		extras = append(extras, extra{Title: title, MarkDown: markdown})
	}
	slices.SortFunc(extras, func(a, b extra) int {
		return strings.Compare(a.Title, b.Title)
	})

	return chartTemplateData{
		ChartDocumentationInfo: info,
		HelmDocsVersion:        helmDocsVersion,
		Values:                 valuesTableRows,
		Extras:                 extras,
		Sections:               valueRowsSectionSorted,
		Files:                  files,
		SkipVersionFooter:      skipVersionFooter,
	}, nil
}

func removeRowsWithoutDescription(valuesTableRows []valueRow) []valueRow {
	var valuesTableRowsWithoutDescription []valueRow
	for i := range valuesTableRows {
		if valuesTableRows[i].AutoDescription != "" || valuesTableRows[i].Description != "" {
			valuesTableRowsWithoutDescription = append(valuesTableRowsWithoutDescription, valuesTableRows[i])
		}
	}
	return valuesTableRowsWithoutDescription
}

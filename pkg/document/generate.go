package document

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/norwoodj/helm-docs/pkg/helm"
	log "github.com/sirupsen/logrus"
)

type ValueRow struct {
	Key         string
	Type        string
	Default     string
	Description string
}

const boolType = "bool"
const floatType = "float"
const intType = "int"
const listType = "list"
const objectType = "object"
const stringType = "string"
const rowFormat = "| %s | %s | %s | %s |\n"


func printRequirementsHeader(f *os.File) {
	f.WriteString("| Repository | Name | Version |\n")
	f.WriteString("|------------|------|---------|\n")
}

func requirementKey(requirement helm.ChartRequirementsItem) string {
	return fmt.Sprintf("%s/%s", requirement.Repository, requirement.Name)
}

func printRequirementsRows(outputFile *os.File, requirements helm.ChartRequirements) {
	sort.Slice(requirements.Dependencies[:], func(i, j int) bool {
		return requirementKey(requirements.Dependencies[i]) < requirementKey(requirements.Dependencies[j])
	})

	for _, r := range requirements.Dependencies {
		outputFile.WriteString(fmt.Sprintf("| %s | %s | %s |\n", r.Repository, r.Name, r.Version))
	}

	outputFile.WriteString("\n\n")
}

func printValuesHeader(outputFile *os.File) {
	outputFile.WriteString("| Key | Type | Default | Description |\n")
	outputFile.WriteString("|-----|------|---------|-------------|\n")
}

func createAtomRow(value interface{}, prefix string, keysToDescriptions map[string]string) ValueRow {
	description := keysToDescriptions[prefix]

	switch value.(type) {
	case bool:
		return ValueRow{
			Key:         prefix,
			Type:        boolType,
			Default:     fmt.Sprintf("%t", value),
			Description: description,
		}
	case float64:
		return ValueRow{
			Key:         prefix,
			Type:        floatType,
			Default:     strconv.FormatFloat(value.(float64), 'f', -1, 64),
			Description: description,
		}
	case int:
		return ValueRow{
			Key:         prefix,
			Type:        intType,
			Default:     fmt.Sprintf("%d", value),
			Description: description,
		}
	case string:
		return ValueRow{
			Key:         prefix,
			Type:        stringType,
			Default:     fmt.Sprintf("\"%s\"", value),
			Description: description,
		}
	case []interface{}:
		return ValueRow{
			Key:         prefix,
			Type:        listType,
			Default:     "[]",
			Description: description,
		}
	case helm.ChartValues:
		return ValueRow{
			Key:         prefix,
			Type:        objectType,
			Default:     "{}",
			Description: description,
		}
	case nil:
		return parseNilValueType(prefix, description)
	}

	return ValueRow{}
}

func parseNilValueType(prefix string, description string) ValueRow {
	// Grab whatever's in between the parentheses of the description and treat it as the type
	r, _ := regexp.Compile("^\\(.*?\\)")
	t := r.FindString(description)

	if len(t) > 0 {
		t = t[1 : len(t)-1]
		description = description[len(t)+3:]
	} else {
		t = stringType
	}

	return ValueRow{
		Key:         prefix,
		Type:        t,
		Default:     "\\<nil\\>",
		Description: description,
	}
}

func createListRows(values []interface{}, prefix string, keysToDescriptions map[string]string) []ValueRow {
	if len(values) == 0 {
		return []ValueRow{createAtomRow(values, prefix, keysToDescriptions)}
	}

	valueRows := []ValueRow{}

	for i, v := range values {
		var nextPrefix string
		if prefix != "" {
			nextPrefix = fmt.Sprintf("%s[%d]", prefix, i)
		} else {
			nextPrefix = fmt.Sprintf("[%d]", i)
		}

		switch v.(type) {
		case helm.ChartValues:
			valueRows = append(valueRows, createValueRows(v.(helm.ChartValues), nextPrefix, keysToDescriptions)...)
		case []interface{}:
			valueRows = append(valueRows, createListRows(v.([]interface{}), nextPrefix, keysToDescriptions)...)
		case bool:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
		case float64:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
		case int:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
		case string:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
			break
		}
	}

	return valueRows
}

func createValueRows(values helm.ChartValues, prefix string, keysToDescriptions map[string]string) []ValueRow {
	if len(values) == 0 {
		return []ValueRow{createAtomRow(values, prefix, keysToDescriptions)}
	}

	valueRows := []ValueRow{}

	for k, v := range values {
		var escapedKey string
		var nextPrefix string

		key := k.(string)
		if strings.Contains(key, ".") {
			escapedKey = fmt.Sprintf("\"%s\"", k)
		} else {
			escapedKey = key
		}

		if prefix != "" {
			nextPrefix = fmt.Sprintf("%s.%s", prefix, escapedKey)
		} else {
			nextPrefix = fmt.Sprintf("%s", escapedKey)
		}

		switch v.(type) {
		case helm.ChartValues:
			valueRows = append(valueRows, createValueRows(v.(helm.ChartValues), nextPrefix, keysToDescriptions)...)
		case []interface{}:
			valueRows = append(valueRows, createListRows(v.([]interface{}), nextPrefix, keysToDescriptions)...)
		case bool:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
		case float64:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
		case int:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
		case string:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
		default:
			valueRows = append(valueRows, createAtomRow(v, nextPrefix, keysToDescriptions))
		}
	}

	sort.Slice(valueRows[:], func(i, j int) bool {
		return valueRows[i].Key < valueRows[j].Key
	})

	return valueRows
}

func printValueRows(f *os.File, values helm.ChartValues, keysToDescriptions map[string]string) {
	valueRows := createValueRows(values, "", keysToDescriptions)
	for _, valueRow := range valueRows {
		f.WriteString(fmt.Sprintf(rowFormat, valueRow.Key, valueRow.Type, valueRow.Default, valueRow.Description))
	}
}

func withNewline(s string) string {
	return fmt.Sprintln(s)
}

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

	chartMeta := chartDocumentationInfo.ChartMeta
	outputFile.WriteString(withNewline(chartMeta.Name))
	outputFile.WriteString(withNewline(strings.Repeat("=", len(chartMeta.Name))))
	outputFile.WriteString(withNewline(chartMeta.Description))

	if chartMeta.Home != "" {
		outputFile.WriteString(fmt.Sprintf("\nThis chart's source code can be found [here](%s)\n\n\n", chartMeta.Home))
	}

	chartRequirements := chartDocumentationInfo.ChartRequirements
	if len(chartRequirements.Dependencies) > 0 {
		outputFile.WriteString("## Chart Requirements\n\n")
		printRequirementsHeader(outputFile)
		printRequirementsRows(outputFile, chartRequirements)
	}

	outputFile.WriteString("## Chart Values\n\n")
	printValuesHeader(outputFile)
	printValueRows(outputFile, chartDocumentationInfo.ChartValues, chartDocumentationInfo.ChartValuesDescriptions)
}

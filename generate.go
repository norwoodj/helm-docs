package main

import (
    "fmt"
    "strings"
    "strconv"
    "sort"
    "os"
)

type ValueRow struct {
    Key string
    Type string
    Default string
    Description string
}

const BOOL_TYPE = "bool"
const FLOAT_TYPE = "float"
const INT_TYPE = "int"
const LIST_TYPE = "list"
const OBJECT_TYPE = "object"
const STRING_TYPE = "string"

const ROW_FORMAT = "| %s | %s | %s | %s |\n"


func printRequirementsHeader(f *os.File) {
    f.WriteString("| Repository | Name | Version |\n")
    f.WriteString("|------------|------|---------|\n")
}

func requirementKey(requirement ChartRequirementsItem) string {
    return fmt.Sprintf("%s/%s", requirement.Repository, requirement.Name)
}

func printRequirementsRows(outputFile *os.File, requirements ChartRequirements) {
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
            return ValueRow {
                Key: prefix,
                Type: BOOL_TYPE,
                Default: fmt.Sprintf("%t", value),
                Description: description,
            }
        case float64:
            return ValueRow {
                Key: prefix,
                Type: FLOAT_TYPE,
                Default: strconv.FormatFloat(value.(float64), 'f', -1, 64),
                Description: description,
            }
        case int:
            return ValueRow {
                Key: prefix,
                Type: INT_TYPE,
                Default: fmt.Sprintf("%d", value),
                Description: description,
            }
        case string:
            return ValueRow {
                Key: prefix,
                Type: STRING_TYPE,
                Default: fmt.Sprintf("\"%s\"", value),
                Description: description,
            }
        case []interface{}:
            return ValueRow {
                Key: prefix,
                Type: LIST_TYPE,
                Default: "[]",
                Description: description,
            }
        case ChartValues:
            return ValueRow {
                Key: prefix,
                Type: OBJECT_TYPE,
                Default: "{}",
                Description: description,
            }
    }

    return ValueRow{}
}

func createListRows(values []interface{}, prefix string, keysToDescriptions map[string]string) []ValueRow {
    if len(values) == 0 {
        return []ValueRow {createAtomRow(values, prefix, keysToDescriptions)}
    }

    valueRows := []ValueRow {}

    for i, v := range values {
        var nextPrefix string
        if prefix != "" {
            nextPrefix = fmt.Sprintf("%s[%d]", prefix, i)
        } else {
            nextPrefix = fmt.Sprintf("[%d]", i)
        }

        switch v.(type) {
            case ChartValues:
                valueRows = append(valueRows, createValueRows(v.(ChartValues), nextPrefix, keysToDescriptions)...)
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

func createValueRows(values ChartValues, prefix string, keysToDescriptions map[string]string) []ValueRow {
    if len(values) == 0 {
        return []ValueRow {createAtomRow(values, prefix, keysToDescriptions)}
    }

    valueRows := []ValueRow {}

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
            case ChartValues:
                valueRows = append(valueRows, createValueRows(v.(ChartValues), nextPrefix, keysToDescriptions)...)
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
        }
    }

    sort.Slice(valueRows[:], func(i, j int) bool {
        return valueRows[i].Key < valueRows[j].Key
    })

    return valueRows
}

func printValueRows(f *os.File, values ChartValues, keysToDescriptions map[string]string) {
    valueRows := createValueRows(values, "", keysToDescriptions)
    for _, valueRow := range valueRows {
        f.WriteString(fmt.Sprintf(ROW_FORMAT, valueRow.Key, valueRow.Type, valueRow.Default, valueRow.Description))
    }
}

func withNewline(s string) string {
    return fmt.Sprintln(s)
}

func getOutputFile(dryRun bool) *os.File {
    if dryRun {
        return os.Stdout
    }

    f, err := os.Create("README.md")

    if err != nil {
        panic(err)
    }

    return f
}

func printDocumentation(debug bool, dryRun bool) {
    outputFile := getOutputFile(dryRun)
    defer outputFile.Close()
    chartMeta := parseChartFile(debug)

    outputFile.WriteString(withNewline(chartMeta.Name))
    outputFile.WriteString(withNewline(strings.Repeat("=", len(chartMeta.Name))))
    outputFile.WriteString(withNewline(chartMeta.Description))
    outputFile.WriteString(fmt.Sprintf("\nThis chart's source code can be found [here](%s)\n\n\n", chartMeta.Home))

    chartRequirements := parseChartRequirementsFile(debug)

    if len(chartRequirements.Dependencies) > 0 {
        outputFile.WriteString("## Chart Requirements\n\n")
        printRequirementsHeader(outputFile)
        printRequirementsRows(outputFile, chartRequirements)
    }

    outputFile.WriteString("## Chart Values\n\n")
    values := parseValuesFile(debug)
    printValuesHeader(outputFile)
    keysToDescriptions := parseValuesFileComments()
    printValueRows(outputFile, values, keysToDescriptions)
}

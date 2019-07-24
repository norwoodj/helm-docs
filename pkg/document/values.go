package document

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/norwoodj/helm-docs/pkg/helm"
)

const (
	boolType   = "bool"
	floatType  = "float"
	intType    = "int"
	listType   = "list"
	objectType = "object"
	stringType = "string"
)

func createAtomRow(prefix string, value interface{}, keysToDescriptions map[string]string) valueRow {
	description := keysToDescriptions[prefix]

	switch value.(type) {
	case bool:
		return valueRow{
			Key:         prefix,
			Type:        boolType,
			Default:     fmt.Sprintf("%t", value),
			Description: description,
		}
	case float64:
		return valueRow{
			Key:         prefix,
			Type:        floatType,
			Default:     strconv.FormatFloat(value.(float64), 'f', -1, 64),
			Description: description,
		}
	case int:
		return valueRow{
			Key:         prefix,
			Type:        intType,
			Default:     fmt.Sprintf("%d", value),
			Description: description,
		}
	case string:
		return valueRow{
			Key:         prefix,
			Type:        stringType,
			Default:     fmt.Sprintf("\"%s\"", value),
			Description: description,
		}
	case []interface{}:
		return valueRow{
			Key:         prefix,
			Type:        listType,
			Default:     "[]",
			Description: description,
		}
	case helm.ChartValues:
		return valueRow{
			Key:         prefix,
			Type:        objectType,
			Default:     "{}",
			Description: description,
		}
	case nil:
		return parseNilValueType(prefix, description)
	}

	return valueRow{}
}

func parseNilValueType(prefix string, description string) valueRow {
	// Grab whatever's in between the parentheses of the description and treat it as the type
	r, _ := regexp.Compile("^\\(.*?\\)")
	t := r.FindString(description)

	if len(t) > 0 {
		t = t[1 : len(t)-1]
		description = description[len(t)+3:]
	} else {
		t = stringType
	}

	return valueRow{
		Key:         prefix,
		Type:        t,
		Default:     "\\<nil\\>",
		Description: description,
	}
}

func createListRows(prefix string, values []interface{}, keysToDescriptions map[string]string) []valueRow {
	valueRows := []valueRow{createAtomRow(prefix, values, keysToDescriptions)}

	if len(values) == 0 {
		return valueRows
	}

	for i, v := range values {
		var nextPrefix string
		if prefix != "" {
			nextPrefix = fmt.Sprintf("%s[%d]", prefix, i)
		} else {
			nextPrefix = fmt.Sprintf("[%d]", i)
		}

		switch v.(type) {
		case helm.ChartValues:
			valueRows = append(valueRows, createValueRows(nextPrefix, v.(helm.ChartValues), keysToDescriptions)...)
		case []interface{}:
			valueRows = append(valueRows, createListRows(nextPrefix, v.([]interface{}), keysToDescriptions)...)
		case bool:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
		case float64:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
		case int:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
		case string:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
			break
		}
	}

	return valueRows
}

func createValueRows(prefix string, values helm.ChartValues, keysToDescriptions map[string]string) []valueRow {
	valueRows := make([]valueRow, 0)

	if prefix != "" {
		valueRows = append(valueRows, createAtomRow(prefix, values, keysToDescriptions))
	}

	if len(values) == 0 {
		return valueRows
	}

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
			valueRows = append(valueRows, createValueRows(nextPrefix, v.(helm.ChartValues), keysToDescriptions)...)
		case []interface{}:
			valueRows = append(valueRows, createListRows(nextPrefix, v.([]interface{}), keysToDescriptions)...)
		case bool:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
		case float64:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
		case int:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
		case string:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
		default:
			valueRows = append(valueRows, createAtomRow(nextPrefix, v, keysToDescriptions))
		}
	}

	sort.Slice(valueRows[:], func(i, j int) bool {
		return valueRows[i].Key < valueRows[j].Key
	})

	return valueRows
}

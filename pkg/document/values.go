package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
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

var nilValueTypeRegex = regexp.MustCompile("^\\(.*?\\)")

func formatNextListKeyPrefix(prefix string, index int) string {
	return fmt.Sprintf("%s[%d]", prefix, index)
}

func formatNextObjectKeyPrefix(prefix string, key string) string {
	var escapedKey string
	var nextPrefix string

	if strings.Contains(key, ".") || strings.Contains(key, " ") {
		escapedKey = fmt.Sprintf(`"%s"`, key)
	} else {
		escapedKey = key
	}

	if prefix != "" {
		nextPrefix = fmt.Sprintf("%s.%s", prefix, escapedKey)
	} else {
		nextPrefix = fmt.Sprintf("%s", escapedKey)
	}

	return nextPrefix
}

func getTypeName(value interface{}) string {
	switch value.(type) {
	case bool:
		return boolType
	case float64:
		return floatType
	case int:
		return intType
	case string:
		return stringType
	case []interface{}:
		return listType
	case jsonableMap:
		return objectType
	}

	return ""
}

func parseNilValueType(key string, description helm.ChartValueDescription) valueRow {
	// Grab whatever's in between the parentheses of the description and treat it as the type
	t := nilValueTypeRegex.FindString(description.Description)

	if len(t) > 0 {
		t = t[1 : len(t)-1]
		description.Description = description.Description[len(t)+3:]
	} else {
		t = stringType
	}

	if description.Default == "" {
		description.Default = "`nil`"
	}

	return valueRow{
		Key:         key,
		Type:        t,
		Default:     description.Default,
		Description: description.Description,
	}
}

func jsonMarshalNoEscape(key string, value interface{}) (string, error) {
	outputBuffer := &bytes.Buffer{}
	valueEncoder := json.NewEncoder(outputBuffer)
	valueEncoder.SetEscapeHTML(false)
	err := valueEncoder.Encode(value)

	if err != nil {
		return "", fmt.Errorf("failed to marshal default value for %s to json: %s", key, err)
	}

	return strings.TrimRight(outputBuffer.String(), "\n"), nil
}

func createValueRow(
	key string,
	value interface{},
	description helm.ChartValueDescription,
) (valueRow, error) {
	if value == nil {
		return parseNilValueType(key, description), nil
	}

	defaultValue := description.Default
	if defaultValue == "" {
		jsonEncodedValue, err := jsonMarshalNoEscape(key, value)
		if err != nil {
			return valueRow{}, fmt.Errorf("failed to marshal default value for %s to json: %s", key, err)
		}

		defaultValue = fmt.Sprintf("`%s`", jsonEncodedValue)
	}

	return valueRow{
		Key:         key,
		Type:        getTypeName(value),
		Default:     defaultValue,
		Description: description.Description,
	}, nil
}

func createRowsFromField(
	nextPrefix string,
	value interface{},
	keysToDescriptions map[string]helm.ChartValueDescription,
	documentLeafNodes bool,
	containerDefaults string,
) ([]valueRow, error) {
	switch value.(type) {
	case map[interface{}]interface{}:
		return createValueRowsFromObject(nextPrefix, value.(map[interface{}]interface{}), keysToDescriptions, documentLeafNodes, containerDefaults)

	case []interface{}:
		return createValueRowsFromList(nextPrefix, value.([]interface{}), keysToDescriptions, documentLeafNodes, containerDefaults)

	default:
		description, hasDescription := keysToDescriptions[nextPrefix]
		if !(documentLeafNodes || hasDescription) {
			return []valueRow{}, nil
		}

		leafValueRow, err := createValueRow(nextPrefix, value, description)
		return []valueRow{leafValueRow}, err
	}
}

func createValueRowsFromList(
	prefix string,
	values []interface{},
	keysToDescriptions map[string]helm.ChartValueDescription,
	documentLeafNodes bool,
	containerDefaults string,
) ([]valueRow, error) {
	description, hasDescription := keysToDescriptions[prefix]

	// If we encounter an empty list, it should be documented if no parent object or list had a description or if this
	// list has a description
	if len(values) == 0 {
		if !(documentLeafNodes || hasDescription) {
			return []valueRow{}, nil
		}

		emptyListRow, err := createValueRow(prefix, values, description)
		if err != nil {
			return nil, err
		}

		return []valueRow{emptyListRow}, nil
	}

	valueRows := make([]valueRow, 0)

	// We have a nonempty list with a description, document it, and mark that leaf nodes underneath it should not be
	// documented without descriptions
	if hasDescription {
		jsonableObject := convertHelmValuesToJsonable(values)

		if containerDefaults != "" {
			// Don't document complex objects
			if description.Default == "" {
				description.Default = containerDefaults
			}
		}

		listRow, err := createValueRow(prefix, jsonableObject, description)

		if err != nil {
			return nil, err
		}

		valueRows = append(valueRows, listRow)
		documentLeafNodes = false
	}

	// Generate documentation rows for all list items and their potential sub-fields
	for i, v := range values {
		nextPrefix := formatNextListKeyPrefix(prefix, i)
		valueRowsForListField, err := createRowsFromField(nextPrefix, v, keysToDescriptions, documentLeafNodes, containerDefaults)

		if err != nil {
			return nil, err
		}

		valueRows = append(valueRows, valueRowsForListField...)
	}

	return valueRows, nil
}

func createValueRowsFromObject(
	prefix string,
	values map[interface{}]interface{},
	keysToDescriptions map[string]helm.ChartValueDescription,
	documentLeafNodes bool,
	containerDefaults string,
) ([]valueRow, error) {
	description, hasDescription := keysToDescriptions[prefix]

	if len(values) == 0 {
		// if the first level of recursion has no values, then there are no values at all, and so we return zero rows of documentation
		if prefix == "" {
			return []valueRow{}, nil
		}

		// Otherwise, we have a leaf empty object node that should be documented if no object up the recursion chain had
		// a description or if this object has a description
		if !(documentLeafNodes || hasDescription) {
			return []valueRow{}, nil
		}

		documentedRow, err := createValueRow(prefix, jsonableMap{}, description)
		return []valueRow{documentedRow}, err
	}

	valueRows := make([]valueRow, 0)

	// We have a nonempty object with a description, document it, and mark that leaf nodes underneath it should not be
	// documented without descriptions
	if hasDescription {
		jsonableObject := convertHelmValuesToJsonable(values)

		if containerDefaults != "" {
			// Don't document complex objects
			if description.Default == "" {
				description.Default = containerDefaults
			}
		}

		objectRow, err := createValueRow(prefix, jsonableObject, description)

		if err != nil {
			return nil, err
		}

		valueRows = append(valueRows, objectRow)
		documentLeafNodes = false
	}

	for k, v := range values {
		nextPrefix := formatNextObjectKeyPrefix(prefix, convertMapKeyToString(k))
		valueRowsForObjectField, err := createRowsFromField(nextPrefix, v, keysToDescriptions, documentLeafNodes, containerDefaults)

		if err != nil {
			return nil, err
		}

		valueRows = append(valueRows, valueRowsForObjectField...)
	}

	// At the top level of recursion, sort value rows by key
	if prefix == "" {
		sort.Slice(valueRows[:], func(i, j int) bool {
			return valueRows[i].Key < valueRows[j].Key
		})
	}

	return valueRows, nil
}

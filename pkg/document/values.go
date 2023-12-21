package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"gopkg.in/yaml.v3"
)

const (
	boolType   = "bool"
	floatType  = "float"
	intType    = "int"
	listType   = "list"
	objectType = "object"
	stringType = "string"
	yamlType   = "yaml"
	tplType    = "tpl"
)

// Yaml tags that differentiate the type of scalar object in the node
const (
	nullTag      = "!!null"
	boolTag      = "!!bool"
	strTag       = "!!str"
	intTag       = "!!int"
	floatTag     = "!!float"
	timestampTag = "!!timestamp"
)

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
	case map[string]interface{}:
		return objectType
	}

	return ""
}

func parseNilValueType(key string, description helm.ChartValueDescription, autoDescription helm.ChartValueDescription, column int, lineNumber int) valueRow {
	if len(description.Description) == 0 {
		description.Description = autoDescription.Description
	}

	var t string
	if description.ValueType != "" {
		t = description.ValueType
	} else if autoDescription.ValueType != "" {
		// Use whatever the type recognized by autoDescription parser
		t = autoDescription.ValueType
	} else {
		t = stringType
	}

	// only set description.Default if no fallback (autoDescription.Default) is available
	if description.Default == "" && autoDescription.Default == "" {
		description.Default = "`nil`"
	}

	section := description.Section
	if section == "" && autoDescription.Section != "" {
		section = autoDescription.Section
	}

	return valueRow{
		Key:             key,
		Type:            t,
		NotationType:    autoDescription.NotationType,
		AutoDefault:     autoDescription.Default,
		Default:         description.Default,
		AutoDescription: autoDescription.Description,
		Description:     description.Description,
		Section:         section,
		Column:          column,
		LineNumber:      lineNumber,
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

func getDescriptionFromNode(node *yaml.Node) helm.ChartValueDescription {
	if node == nil {
		return helm.ChartValueDescription{}
	}

	if node.HeadComment == "" {
		return helm.ChartValueDescription{}
	}

	if !strings.Contains(node.HeadComment, helm.PrefixComment) {
		return helm.ChartValueDescription{}
	}
	commentLines := strings.Split(node.HeadComment, "\n")
	keyFromComment, c := helm.ParseComment(commentLines)
	if keyFromComment != "" {
		return helm.ChartValueDescription{}
	}

	return c
}

func createValueRow(
	key string,
	value interface{},
	description helm.ChartValueDescription,
	autoDescription helm.ChartValueDescription,
	column int,
	lineNumber int,
) (valueRow, error) {
	if value == nil {
		return parseNilValueType(key, description, autoDescription, column, lineNumber), nil
	}

	autoDefaultValue := autoDescription.Default
	defaultValue := description.Default
	notationType := autoDescription.NotationType
	defaultType := getTypeName(value)
	if description.ValueType != "" {
		defaultType = description.ValueType
	} else if autoDescription.ValueType != "" {
		defaultType = autoDescription.ValueType
	} else if notationType != "" {
		// If nothing can be inferred then infer from notationType
		defaultType = notationType
	}

	if defaultValue == "" && autoDefaultValue == "" && notationType == "" {
		jsonEncodedValue, err := jsonMarshalNoEscape(key, value)
		if err != nil {
			return valueRow{}, fmt.Errorf("failed to marshal default value for %s to json: %s", key, err)
		}

		defaultValue = fmt.Sprintf("`%s`", jsonEncodedValue)
	}

	if defaultValue == "" && autoDefaultValue == "" && notationType != "" {
		// We want to render custom styles for custom NotationType
		// So, output a raw default value for this and let the template handle it
		defaultValue = fmt.Sprintf("%s", value)
	}

	section := description.Section
	if section == "" && autoDescription.Section != "" {
		section = autoDescription.Section
	}

	return valueRow{
		Key:             key,
		Type:            defaultType,
		NotationType:    notationType,
		AutoDefault:     autoDescription.Default,
		Default:         defaultValue,
		AutoDescription: autoDescription.Description,
		Description:     description.Description,
		Section:         section,
		Column:          column,
		LineNumber:      lineNumber,
	}, nil
}

func createValueRowsFromList(
	prefix string,
	key *yaml.Node,
	values *yaml.Node,
	keysToDescriptions map[string]helm.ChartValueDescription,
	documentLeafNodes bool,
) ([]valueRow, error) {
	description, hasDescription := keysToDescriptions[prefix]
	autoDescription := getDescriptionFromNode(key)

	// If we encounter an empty list, it should be documented if no parent object or list had a description or if this
	// list has a description
	if len(values.Content) == 0 {
		if !(documentLeafNodes || hasDescription || autoDescription.Description != "") {
			return []valueRow{}, nil
		}

		emptyListRow, err := createValueRow(prefix, make([]interface{}, 0), description, autoDescription, key.Column, key.Line)
		if err != nil {
			return nil, err
		}

		return []valueRow{emptyListRow}, nil
	}

	valueRows := make([]valueRow, 0)

	// We have a nonempty list with a description, document it, and mark that leaf nodes underneath it should not be
	// documented without descriptions
	if hasDescription || (autoDescription.Description != "" && autoDescription.NotationType == "") {
		jsonableObject := convertHelmValuesToJsonable(values)
		listRow, err := createValueRow(prefix, jsonableObject, description, autoDescription, key.Column, key.Line)

		if err != nil {
			return nil, err
		}

		valueRows = append(valueRows, listRow)
		documentLeafNodes = false
	} else if hasDescription || (autoDescription.Description != "" && autoDescription.NotationType != "") {
		// If it has NotationType described, then use that
		var notationValue interface{}
		var err error
		var listRow valueRow
		switch autoDescription.NotationType {
		case yamlType:
			notationValue, err = yaml.Marshal(values)
			if err != nil {
				return nil, err
			}

			listRow, err = createValueRow(prefix, notationValue, description, autoDescription, key.Column, key.Line)

			if err != nil {
				return nil, err
			}
		default:
			// Any other case means we let the template renderer to decide how to
			// format the default value. But the value are stored as raw string
			fallthrough
		case tplType:
			notationValue = values.Value
			listRow, err = createValueRow(prefix, notationValue, description, autoDescription, key.Column, key.Line)

			if err != nil {
				return nil, err
			}
		}

		valueRows = append(valueRows, listRow)
		documentLeafNodes = false
	}

	// Generate documentation rows for all list items and their potential sub-fields
	for i, v := range values.Content {
		nextPrefix := formatNextListKeyPrefix(prefix, i)
		valueRowsForListField, err := createValueRowsFromField(nextPrefix, v, v, keysToDescriptions, documentLeafNodes)

		if err != nil {
			return nil, err
		}

		valueRows = append(valueRows, valueRowsForListField...)
	}

	return valueRows, nil
}

func createValueRowsFromObject(
	nextPrefix string,
	key *yaml.Node,
	values *yaml.Node,
	keysToDescriptions map[string]helm.ChartValueDescription,
	documentLeafNodes bool,
) ([]valueRow, error) {
	description, hasDescription := keysToDescriptions[nextPrefix]
	autoDescription := getDescriptionFromNode(key)

	if len(values.Content) == 0 {
		// if the first level of recursion has no values, then there are no values at all, and so we return zero rows of documentation
		if nextPrefix == "" {
			return []valueRow{}, nil
		}

		// Otherwise, we have a leaf empty object node that should be documented if no object up the recursion chain had
		// a description or if this object has a description
		if !(documentLeafNodes || hasDescription || autoDescription.Description != "") {
			return []valueRow{}, nil
		}

		documentedRow, err := createValueRow(nextPrefix, make(map[string]interface{}), description, autoDescription, key.Column, key.Line)
		return []valueRow{documentedRow}, err
	}

	valueRows := make([]valueRow, 0)

	// We have a nonempty object with a description, document it, and mark that leaf nodes underneath it should not be
	// documented without descriptions
	if hasDescription || (autoDescription.Description != "" && autoDescription.NotationType == "") {
		jsonableObject := convertHelmValuesToJsonable(values)
		objectRow, err := createValueRow(nextPrefix, jsonableObject, description, autoDescription, key.Column, key.Line)

		if err != nil {
			return nil, err
		}

		valueRows = append(valueRows, objectRow)
		documentLeafNodes = false
	} else if hasDescription || (autoDescription.Description != "" && autoDescription.NotationType != "") {

		// If it has NotationType described, then use that
		var notationValue interface{}
		var err error
		var objectRow valueRow
		switch autoDescription.NotationType {
		case yamlType:
			notationValue, err = yaml.Marshal(values)
			if err != nil {
				return nil, err
			}

			objectRow, err = createValueRow(nextPrefix, notationValue, description, autoDescription, key.Column, key.Line)

			if err != nil {
				return nil, err
			}

		default:
			// Any other case means we let the template renderer to decide how to
			// format the default value. But the value are stored as raw string
			fallthrough
		case tplType:
			notationValue = values.Value
			objectRow, err = createValueRow(nextPrefix, notationValue, description, autoDescription, key.Column, key.Line)

			if err != nil {
				return nil, err
			}
		}

		valueRows = append(valueRows, objectRow)
		documentLeafNodes = false
	}

	for i := 0; i < len(values.Content); i += 2 {
		k := values.Content[i]
		v := values.Content[i+1]
		nextPrefix := formatNextObjectKeyPrefix(nextPrefix, k.Value)
		valueRowsForObjectField, err := createValueRowsFromField(nextPrefix, k, v, keysToDescriptions, documentLeafNodes)

		if err != nil {
			return nil, err
		}

		valueRows = append(valueRows, valueRowsForObjectField...)
	}

	return valueRows, nil
}

func createValueRowsFromField(
	prefix string,
	key *yaml.Node,
	value *yaml.Node,
	keysToDescriptions map[string]helm.ChartValueDescription,
	documentLeafNodes bool,
) ([]valueRow, error) {
	switch value.Kind {
	case yaml.MappingNode:
		return createValueRowsFromObject(prefix, key, value, keysToDescriptions, documentLeafNodes)
	case yaml.SequenceNode:
		return createValueRowsFromList(prefix, key, value, keysToDescriptions, documentLeafNodes)
	case yaml.AliasNode:
		return createValueRowsFromField(prefix, key, value.Alias, keysToDescriptions, documentLeafNodes)
	case yaml.ScalarNode:
		autoDescription := getDescriptionFromNode(key)
		description, hasDescription := keysToDescriptions[prefix]
		if !(documentLeafNodes || hasDescription || autoDescription.Description != "") {
			return []valueRow{}, nil
		}

		switch value.Tag {
		case nullTag:
			leafValueRow, err := createValueRow(prefix, nil, description, autoDescription, key.Column, key.Line)
			return []valueRow{leafValueRow}, err
		case strTag:
			// extra check to see if the node is a string, but @notationType was declared
			if autoDescription.NotationType != "" {
				var notationValue interface{}
				var err error
				var leafValueRow valueRow
				switch autoDescription.NotationType {
				case yamlType:
					notationValue, err = yaml.Marshal(value)
					if err != nil {
						return nil, err
					}

					leafValueRow, err = createValueRow(prefix, notationValue, description, autoDescription, key.Column, key.Line)

					if err != nil {
						return nil, err
					}

					return []valueRow{leafValueRow}, err
				default:
					// Any other case means we let the template renderer to decide how to
					// format the default value. But the value are stored as raw string
					fallthrough
				case tplType:
					notationValue = value.Value
					leafValueRow, err = createValueRow(prefix, notationValue, description, autoDescription, key.Column, key.Line)

					if err != nil {
						return nil, err
					}

					return []valueRow{leafValueRow}, err
				}
			}
			fallthrough
		case timestampTag:
			leafValueRow, err := createValueRow(prefix, value.Value, description, autoDescription, key.Column, key.Line)
			return []valueRow{leafValueRow}, err
		case intTag:
			var decodedValue int
			err := value.Decode(&decodedValue)
			if err != nil {
				return []valueRow{}, err
			}

			leafValueRow, err := createValueRow(prefix, decodedValue, description, autoDescription, key.Column, key.Line)
			return []valueRow{leafValueRow}, err
		case floatTag:
			var decodedValue float64
			err := value.Decode(&decodedValue)
			if err != nil {
				return []valueRow{}, err
			}
			leafValueRow, err := createValueRow(prefix, decodedValue, description, autoDescription, key.Column, key.Line)
			return []valueRow{leafValueRow}, err

		case boolTag:
			var decodedValue bool
			err := value.Decode(&decodedValue)
			if err != nil {
				return []valueRow{}, err
			}
			leafValueRow, err := createValueRow(prefix, decodedValue, description, autoDescription, key.Column, key.Line)
			return []valueRow{leafValueRow}, err
		}
	}

	return []valueRow{}, fmt.Errorf("invalid node type %d received", value.Kind)
}

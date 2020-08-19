package document

import (
	"strings"
	"testing"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func parseYamlValues(yamlValues string) map[interface{}]interface{} {
	var chartValues map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(strings.TrimSpace(yamlValues)), &chartValues)

	if err != nil {
		panic(err)
	}

	return chartValues
}

func TestEmptyValues(t *testing.T) {
	valuesRows, err := createValueRowsFromObject("", make(map[interface{}]interface{}), make(map[string]helm.ChartValueDescription), true, "")
	assert.Nil(t, err)
	assert.Len(t, valuesRows, 0)
}

func TestSimpleValues(t *testing.T) {
	helmValues := parseYamlValues(`
echo: 0
foxtrot: true
hello: "world"
oscar: 3.14159
	`)

	valuesRows, err := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "`0`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "`true`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"world\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "`3.14159`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].Description)

	// Check omit-blanks behaviour - no fields output
	valuesRowsBlank, errBlank := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), false, "")

	assert.Nil(t, errBlank)
	assert.Len(t, valuesRowsBlank, 0)
}

func TestSimpleValuesWithDescriptions(t *testing.T) {
	helmValues := parseYamlValues(`
echo: 0
foxtrot: true
hello: "world"
oscar: 3.14159
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"echo":    {Description: "echo"},
		"foxtrot": {Description: "foxtrot"},
		"hello":   {Description: "hello"},
		"oscar":   {Description: "oscar"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")
	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "`0`", valuesRows[0].Default)
	assert.Equal(t, "echo", valuesRows[0].Description)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "`true`", valuesRows[1].Default)
	assert.Equal(t, "foxtrot", valuesRows[1].Description)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"world\"`", valuesRows[2].Default)
	assert.Equal(t, "hello", valuesRows[2].Description)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "`3.14159`", valuesRows[3].Default)
	assert.Equal(t, "oscar", valuesRows[3].Description)

	// Check skip-blanks behaviour - no fields output
	valuesRowsBlank, errBlank := createValueRowsFromObject("", helmValues, descriptions, false, "")
	assert.Nil(t, errBlank)
	assert.Len(t, valuesRowsBlank, 4)
}

func TestSimpleValuesWithDescriptionsAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
echo: 0
foxtrot: true
hello: "world"
oscar: 3.14159
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"echo":    {Description: "echo", Default: "some"},
		"foxtrot": {Description: "foxtrot", Default: "explicit"},
		"hello":   {Description: "hello", Default: "default"},
		"oscar":   {Description: "oscar", Default: "values"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "echo", valuesRows[0].Description)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "foxtrot", valuesRows[1].Description)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "hello", valuesRows[2].Description)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "values", valuesRows[3].Default)
	assert.Equal(t, "oscar", valuesRows[3].Description)
}

func TestRecursiveValues(t *testing.T) {
	helmValues := parseYamlValues(`
recursive:
  echo: cat
oscar: dog
	`)

	valuesRows, err := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}

func TestRecursiveValuesWithDescriptions(t *testing.T) {
	helmValues := parseYamlValues(`
recursive:
  echo: cat
oscar: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"recursive.echo": {Description: "echo"},
		"oscar":          {Description: "oscar"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "oscar", valuesRows[0].Description)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "echo", valuesRows[1].Description)
}

func TestRecursiveValuesWithDescriptionsAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
recursive:
  echo: cat
oscar: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"recursive.echo": {Description: "echo", Default: "custom"},
		"oscar":          {Description: "oscar", Default: "default"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "default", valuesRows[0].Default)
	assert.Equal(t, "oscar", valuesRows[0].Description)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "custom", valuesRows[1].Default)
	assert.Equal(t, "echo", valuesRows[1].Description)
}

func TestEmptyObject(t *testing.T) {
	helmValues := parseYamlValues(`
recursive: {}
oscar: dog
	`)

	valuesRows, err := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key, "oscar")
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{}`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}

func TestEmptyObjectWithDescription(t *testing.T) {
	helmValues := parseYamlValues(`
recursive: {}
oscar: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"recursive": {Description: "an empty object"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{}`", valuesRows[1].Default)
	assert.Equal(t, "an empty object", valuesRows[1].Description)
}

func TestEmptyObjectWithDescriptionAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
recursive: {}
oscar: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"recursive": {Description: "an empty object", Default: "default"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "default", valuesRows[1].Default)
	assert.Equal(t, "an empty object", valuesRows[1].Description)

	// Check skip-blanks behaviour - only fields with descriptions output
	valuesRowsBlanks, errBlanks := createValueRowsFromObject("", helmValues, descriptions, false, "")

	assert.Nil(t, errBlanks)
	assert.Len(t, valuesRowsBlanks, 1)

	assert.Equal(t, "recursive", valuesRowsBlanks[0].Key)
	assert.Equal(t, objectType, valuesRowsBlanks[0].Type)
	assert.Equal(t, "default", valuesRowsBlanks[0].Default)
	assert.Equal(t, "an empty object", valuesRowsBlanks[0].Description)

}
func TestEmptyList(t *testing.T) {
	helmValues := parseYamlValues(`
birds: []
echo: cat
	`)

	valuesRows, err := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[]`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}

func TestEmptyListWithDescriptions(t *testing.T) {
	helmValues := parseYamlValues(`
birds: []
echo: cat
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"birds": {Description: "birds"},
		"echo":  {Description: "echo"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[]`", valuesRows[0].Default)
	assert.Equal(t, "birds", valuesRows[0].Description)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "echo", valuesRows[1].Description)
}

func TestEmptyListWithDescriptionsAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
birds: []
echo: cat
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"birds": {Description: "birds", Default: "explicit"},
		"echo":  {Description: "echo", Default: "default value"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "birds", valuesRows[0].Description)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default value", valuesRows[1].Default)
	assert.Equal(t, "echo", valuesRows[1].Description)
}

func TestListOfStrings(t *testing.T) {
	helmValues := parseYamlValues(`
cats: [echo, foxtrot]
	`)

	valuesRows, err := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

}

func TestListOfStringsWithDescriptions(t *testing.T) {
	helmValues := parseYamlValues(`
cats: [echo, foxtrot]
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"cats[0]": {Description: "the black one"},
		"cats[1]": {Description: "the friendly one"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "the black one", valuesRows[0].Description)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)

}

func TestListOfStringsWithDescriptionsAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
cats: [echo, foxtrot]
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"cats[0]": {Description: "the black one", Default: "explicit"},
		"cats[1]": {Description: "the friendly one", Default: "default value"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "the black one", valuesRows[0].Description)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default value", valuesRows[1].Default)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)

}

func TestListOfObjects(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	valuesRows, err := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].Description)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].Description)
}

func TestListOfObjectsWithDescriptions(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals[0].elements[0]": {Description: "the black one"},
		"animals[0].elements[1]": {Description: "the friendly one"},
		"animals[1].elements[0]": {Description: "the sleepy one"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "the black one", valuesRows[0].Description)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[3].Default)
	assert.Equal(t, "the sleepy one", valuesRows[3].Description)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].Description)
}

func TestListOfObjectsWithDescriptionsAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals[0].elements[0]": {Description: "the black one", Default: "explicit"},
		"animals[0].elements[1]": {Description: "the friendly one", Default: "default"},
		"animals[1].elements[0]": {Description: "the sleepy one", Default: "value"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "the black one", valuesRows[0].Description)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default", valuesRows[1].Default)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "value", valuesRows[3].Default)
	assert.Equal(t, "the sleepy one", valuesRows[3].Description)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].Description)
}

func TestDescriptionOnList(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals": {Description: "all the animals of the house"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[{\"elements\":[\"echo\",\"foxtrot\"],\"type\":\"cat\"},{\"elements\":[\"oscar\"],\"type\":\"dog\"}]`", valuesRows[0].Default)
	assert.Equal(t, "all the animals of the house", valuesRows[0].Description)

	// Repeat with containerDefaults set
	valuesRows2, err2 := createValueRowsFromObject("", helmValues, descriptions, true, "EMPTY")

	assert.Nil(t, err2)
	assert.Len(t, valuesRows2, 1)

	assert.Equal(t, "animals", valuesRows2[0].Key)
	assert.Equal(t, listType, valuesRows2[0].Type)
	assert.Equal(t, "EMPTY", valuesRows2[0].Default)
	assert.Equal(t, "all the animals of the house", valuesRows2[0].Description)
}

func TestDescriptionAndDefaultOnList(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals": {Description: "all the animals of the house", Default: "cat and dog"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "cat and dog", valuesRows[0].Default)
	assert.Equal(t, "all the animals of the house", valuesRows[0].Description)
}

func TestDescriptionAndDefaultOnObjectUnderList(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals[0]": {Description: "all the cats of the house", Default: "only cats here"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals[0]", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "only cats here", valuesRows[0].Default)
	assert.Equal(t, "all the cats of the house", valuesRows[0].Description)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

	assert.Equal(t, "animals[1].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)
}

func TestDescriptionOnObjectUnderObject(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  byTrait:
    friendly: [foxtrot, oscar]
    mean: [echo]
    sleepy: [oscar]
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals.byTrait": {Description: "animals listed by their various characteristics"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.byTrait", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "`{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}`", valuesRows[0].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[0].Description)

	// Repeat with containerDefaults set.  default will be string ""
	valuesRows2, err2 := createValueRowsFromObject("", helmValues, descriptions, true, "BLANK")

	assert.Nil(t, err2)
	assert.Len(t, valuesRows2, 1)

	assert.Equal(t, "animals.byTrait", valuesRows2[0].Key)
	assert.Equal(t, objectType, valuesRows2[0].Type)
	assert.Equal(t, "BLANK", valuesRows2[0].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows2[0].Description)
}

func TestDescriptionAndDefaultOnObjectUnderObject(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  byTrait:
    friendly: [foxtrot, oscar]
    mean: [echo]
    sleepy: [oscar]
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals.byTrait": {Description: "animals listed by their various characteristics", Default: "animals, you know"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.byTrait", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "animals, you know", valuesRows[0].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[0].Description)
}

func TestDescriptionsDownChain(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  byTrait:
    friendly: [foxtrot, oscar]
    mean: [echo]
    sleepy: [oscar]
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals":                     {Description: "animal stuff"},
		"animals.byTrait":             {Description: "animals listed by their various characteristics"},
		"animals.byTrait.friendly":    {Description: "the friendly animals of the house"},
		"animals.byTrait.friendly[0]": {Description: "best cat ever"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "`{\"byTrait\":{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}}`", valuesRows[0].Default)
	assert.Equal(t, "animal stuff", valuesRows[0].Description)

	assert.Equal(t, "animals.byTrait", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}`", valuesRows[1].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[1].Description)

	assert.Equal(t, "animals.byTrait.friendly", valuesRows[2].Key)
	assert.Equal(t, listType, valuesRows[2].Type)
	assert.Equal(t, "`[\"foxtrot\",\"oscar\"]`", valuesRows[2].Default)
	assert.Equal(t, "the friendly animals of the house", valuesRows[2].Description)

	assert.Equal(t, "animals.byTrait.friendly[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[3].Default)
	assert.Equal(t, "best cat ever", valuesRows[3].Description)
}

func TestDescriptionsDownChainEmptyDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  byTrait:
    friendly: [foxtrot, oscar]
    mean: [echo]
    sleepy: [oscar]
  `)

	descriptions := map[string]helm.ChartValueDescription{
		"animals":                     {Description: "animal stuff"},
		"animals.byTrait":             {Description: "animals listed by their various characteristics"},
		"animals.byTrait.friendly":    {Description: "the friendly animals of the house"},
		"animals.byTrait.friendly[0]": {Description: "best cat ever"},
	}

	// containerDefaults set - Defaults output will be ""
	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "EMPTY")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "EMPTY", valuesRows[0].Default)
	assert.Equal(t, "animal stuff", valuesRows[0].Description)

	assert.Equal(t, "animals.byTrait", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "EMPTY", valuesRows[1].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[1].Description)

	assert.Equal(t, "animals.byTrait.friendly", valuesRows[2].Key)
	assert.Equal(t, listType, valuesRows[2].Type)
	assert.Equal(t, "EMPTY", valuesRows[2].Default)
	assert.Equal(t, "the friendly animals of the house", valuesRows[2].Description)

	assert.Equal(t, "animals.byTrait.friendly[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[3].Default)
	assert.Equal(t, "best cat ever", valuesRows[3].Description)
}

func TestDescriptionsAndDefaultsDownChain(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  byTrait:
    friendly: [foxtrot, oscar]
    mean: [echo]
    sleepy: [oscar]
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals":                     {Description: "animal stuff", Default: "some"},
		"animals.byTrait":             {Description: "animals listed by their various characteristics", Default: "explicit"},
		"animals.byTrait.friendly":    {Description: "the friendly animals of the house", Default: "default"},
		"animals.byTrait.friendly[0]": {Description: "best cat ever", Default: "value"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "animal stuff", valuesRows[0].Description)

	assert.Equal(t, "animals.byTrait", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[1].Description)

	assert.Equal(t, "animals.byTrait.friendly", valuesRows[2].Key)
	assert.Equal(t, listType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "the friendly animals of the house", valuesRows[2].Description)

	assert.Equal(t, "animals.byTrait.friendly[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "value", valuesRows[3].Default)
	assert.Equal(t, "best cat ever", valuesRows[3].Description)
}

func TestNilValues(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  birds:
  birdCount:
  nonWeirdCats:
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals.birdCount":    {Description: "(int) the number of birds we have"},
		"animals.birds":        {Description: "(list) the list of birds we have"},
		"animals.nonWeirdCats": {Description: "the cats that we have that are not weird"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals.birdCount", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type)
	assert.Equal(t, "`nil`", valuesRows[0].Default)
	assert.Equal(t, "the number of birds we have", valuesRows[0].Description)

	assert.Equal(t, "animals.birds", valuesRows[1].Key)
	assert.Equal(t, listType, valuesRows[1].Type)
	assert.Equal(t, "`nil`", valuesRows[1].Default)
	assert.Equal(t, "the list of birds we have", valuesRows[1].Description)

	assert.Equal(t, "animals.nonWeirdCats", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`nil`", valuesRows[2].Default)
	assert.Equal(t, "the cats that we have that are not weird", valuesRows[2].Description)
}

func TestNilValuesWithDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  birds:
  birdCount:
  nonWeirdCats:
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"animals.birdCount":    {Description: "(int) the number of birds we have", Default: "some"},
		"animals.birds":        {Description: "(list) the list of birds we have", Default: "explicit"},
		"animals.nonWeirdCats": {Description: "the cats that we have that are not weird", Default: "default"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals.birdCount", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "the number of birds we have", valuesRows[0].Description)

	assert.Equal(t, "animals.birds", valuesRows[1].Key)
	assert.Equal(t, listType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "the list of birds we have", valuesRows[1].Description)

	assert.Equal(t, "animals.nonWeirdCats", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "the cats that we have that are not weird", valuesRows[2].Description)
}

func TestKeysWithSpecialCharacters(t *testing.T) {
	helmValues := parseYamlValues(`
websites:
  stupidchess.jmn23.com: defunct
fullNames:
  John Norwood: me
`)

	valuesRows, err := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"me\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"defunct\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}

func TestKeysWithSpecialCharactersWithDescriptions(t *testing.T) {
	helmValues := parseYamlValues(`
websites:
  stupidchess.jmn23.com: defunct
fullNames:
  John Norwood: me
`)

	descriptions := map[string]helm.ChartValueDescription{
		`fullNames."John Norwood"`:         {Description: "who am I"},
		`websites."stupidchess.jmn23.com"`: {Description: "status of the stupidchess website"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"me\"`", valuesRows[0].Default)
	assert.Equal(t, "who am I", valuesRows[0].Description)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"defunct\"`", valuesRows[1].Default)
	assert.Equal(t, "status of the stupidchess website", valuesRows[1].Description)
}

func TestKeysWithSpecialCharactersWithDescriptionsAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
websites:
  stupidchess.jmn23.com: defunct
fullNames:
  John Norwood: me
`)

	descriptions := map[string]helm.ChartValueDescription{
		`fullNames."John Norwood"`:         {Description: "who am I", Default: "default"},
		`websites."stupidchess.jmn23.com"`: {Description: "status of the stupidchess website", Default: "value"},
	}

	valuesRows, err := createValueRowsFromObject("", helmValues, descriptions, true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "default", valuesRows[0].Default)
	assert.Equal(t, "who am I", valuesRows[0].Description)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "value", valuesRows[1].Default)
	assert.Equal(t, "status of the stupidchess website", valuesRows[1].Description)
}

func TestNonStringKeys(t *testing.T) {
	helmValues := parseYamlValues(`
3: three
3.14159: pi
true: "true"
`)

	valuesRows, err := createValueRowsFromObject("", helmValues, make(map[string]helm.ChartValueDescription), true, "")

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, `"float(3.141590)"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"pi\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "bool(true)", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"true\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

	assert.Equal(t, "int(3)", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"three\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)
}

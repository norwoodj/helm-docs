package document

import (
	"strings"
	"testing"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func parseYamlValues(yamlValues string) *yaml.Node {
	var chartValues yaml.Node
	err := yaml.Unmarshal([]byte(strings.TrimSpace(yamlValues)), &chartValues)

	if err != nil {
		panic(err)
	}

	return chartValues.Content[0]
}

func TestEmptyValues(t *testing.T) {
	helmValues := parseYamlValues(`{}`)
	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))
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

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "`0`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "`true`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"world\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "`3.14159`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "", valuesRows[3].Description)
	assert.Equal(t, "", valuesRows[3].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "`0`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "echo", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "`true`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "foxtrot", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"world\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "hello", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "`3.14159`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "oscar", valuesRows[3].Description)
	assert.Equal(t, "", valuesRows[3].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "echo", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "foxtrot", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "hello", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "values", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "oscar", valuesRows[3].Description)
	assert.Equal(t, "", valuesRows[3].AutoDescription)
}

func TestRecursiveValues(t *testing.T) {
	helmValues := parseYamlValues(`
recursive:
  echo: cat
oscar: dog
	`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "oscar", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "echo", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "default", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "oscar", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "custom", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "echo", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
}

func TestEmptyObject(t *testing.T) {
	helmValues := parseYamlValues(`
recursive: {}
oscar: dog
	`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key, "oscar")
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{}`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
}

func TestEmptyObjectWithDescription(t *testing.T) {
	helmValues := parseYamlValues(`
recursive: {}
oscar: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"recursive": {Description: "an empty object"},
	}

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{}`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "an empty object", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
}

func TestEmptyObjectWithDescriptionAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
recursive: {}
oscar: dog
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"recursive": {Description: "an empty object", Default: "default"},
	}

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "default", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "an empty object", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
}
func TestEmptyList(t *testing.T) {
	helmValues := parseYamlValues(`
birds: []
echo: cat
	`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[]`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[]`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "birds", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "echo", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "birds", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default value", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "echo", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
}

func TestListOfStrings(t *testing.T) {
	helmValues := parseYamlValues(`
cats: [echo, foxtrot]
	`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

}

func TestListOfStringsWithDescriptions(t *testing.T) {
	helmValues := parseYamlValues(`
cats: [echo, foxtrot]
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"cats[0]": {Description: "the black one"},
		"cats[1]": {Description: "the friendly one"},
	}

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "the black one", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

}

func TestListOfStringsWithDescriptionsAndDefaults(t *testing.T) {
	helmValues := parseYamlValues(`
cats: [echo, foxtrot]
	`)

	descriptions := map[string]helm.ChartValueDescription{
		"cats[0]": {Description: "the black one", Default: "explicit"},
		"cats[1]": {Description: "the friendly one", Default: "default value"},
	}

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "the black one", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default value", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

}

func TestListOfObjects(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "", valuesRows[3].Description)
	assert.Equal(t, "", valuesRows[3].AutoDescription)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].AutoDefault)
	assert.Equal(t, "", valuesRows[4].Description)
	assert.Equal(t, "", valuesRows[4].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "the black one", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "the sleepy one", valuesRows[3].Description)
	assert.Equal(t, "", valuesRows[3].AutoDescription)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].AutoDefault)
	assert.Equal(t, "", valuesRows[4].Description)
	assert.Equal(t, "", valuesRows[4].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "the black one", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "value", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "the sleepy one", valuesRows[3].Description)
	assert.Equal(t, "", valuesRows[3].AutoDescription)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].AutoDefault)
	assert.Equal(t, "", valuesRows[4].Description)
	assert.Equal(t, "", valuesRows[4].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[{\"elements\":[\"echo\",\"foxtrot\"],\"type\":\"cat\"},{\"elements\":[\"oscar\"],\"type\":\"dog\"}]`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "all the animals of the house", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "cat and dog", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "all the animals of the house", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals[0]", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "only cats here", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "all the cats of the house", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals[1].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.byTrait", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "`{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.byTrait", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "animals, you know", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "`{\"byTrait\":{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}}`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "animal stuff", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals.byTrait", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals.byTrait.friendly", valuesRows[2].Key)
	assert.Equal(t, listType, valuesRows[2].Type)
	assert.Equal(t, "`[\"foxtrot\",\"oscar\"]`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "the friendly animals of the house", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "animals.byTrait.friendly[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "best cat ever", valuesRows[3].Description)
	assert.Equal(t, "", valuesRows[3].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "animal stuff", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals.byTrait", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals.byTrait.friendly", valuesRows[2].Key)
	assert.Equal(t, listType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "the friendly animals of the house", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "animals.byTrait.friendly[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "value", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "best cat ever", valuesRows[3].Description)
	assert.Equal(t, "", valuesRows[3].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals.birdCount", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type)
	assert.Equal(t, "`nil`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "the number of birds we have", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals.birds", valuesRows[1].Key)
	assert.Equal(t, listType, valuesRows[1].Type)
	assert.Equal(t, "`nil`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "the list of birds we have", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals.nonWeirdCats", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`nil`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "the cats that we have that are not weird", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals.birdCount", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "the number of birds we have", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals.birds", valuesRows[1].Key)
	assert.Equal(t, listType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "the list of birds we have", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals.nonWeirdCats", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "the cats that we have that are not weird", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)
}

func TestKeysWithSpecialCharacters(t *testing.T) {
	helmValues := parseYamlValues(`
websites:
  stupidchess.jmn23.com: defunct
fullNames:
  John Norwood: me
`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"me\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"defunct\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"me\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "who am I", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"defunct\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "status of the stupidchess website", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
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

	valuesRows, err := getSortedValuesTableRows(helmValues, descriptions)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "default", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "who am I", valuesRows[0].Description)
	assert.Equal(t, "", valuesRows[0].AutoDescription)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "value", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "status of the stupidchess website", valuesRows[1].Description)
	assert.Equal(t, "", valuesRows[1].AutoDescription)
}

func TestSimpleAutoDoc(t *testing.T) {
	helmValues := parseYamlValues(`
# -- on a scale of 0 to 9 how mean is echo
echo: 8

# -- is she friendly?
foxtrot: true

# doesn't show up
hello: "world"

# -- his favorite food in number format
oscar: 3.14159
	`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "`8`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "on a scale of 0 to 9 how mean is echo", valuesRows[0].AutoDescription)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "`true`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "is she friendly?", valuesRows[1].AutoDescription)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"world\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "", valuesRows[2].Description)
	assert.Equal(t, "", valuesRows[2].AutoDescription)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "`3.14159`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].AutoDefault)
	assert.Equal(t, "", valuesRows[3].Description)
	assert.Equal(t, "his favorite food in number format", valuesRows[3].AutoDescription)
}

func TestAutoDocNested(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  cats:
    # -- on a scale of 0 to 9 how mean is echo
    echo: 8

# -- is she friendly?
    foxtrot: true

  dogs:
# -- his favorite food in number format
    oscar: 3.14159
`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals.cats.echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "`8`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "on a scale of 0 to 9 how mean is echo", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals.cats.foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "`true`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "is she friendly?", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals.dogs.oscar", valuesRows[2].Key)
	assert.Equal(t, floatType, valuesRows[2].Type)
	assert.Equal(t, "`3.14159`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].AutoDefault)
	assert.Equal(t, "", valuesRows[2].Description)
	assert.Equal(t, "his favorite food in number format", valuesRows[2].AutoDescription)
}

func TestAutoDocList(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  cats:
    # -- best cat, really
    - echo
    # -- trash cat, really
    - foxtrot
`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "animals.cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "best cat, really", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals.cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "trash cat, really", valuesRows[1].AutoDescription)
}

func TestAutoDocListOfObjects(t *testing.T) {
	helmValues := parseYamlValues(`
animalLocations:
  # -- place with the most cats
  - place: home
    cats:
      - echo
      - foxtrot

  # -- place with the fewest cats
  - place: work
    cats: []
`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "animalLocations[0]", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "`{\"cats\":[\"echo\",\"foxtrot\"],\"place\":\"home\"}`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "place with the most cats", valuesRows[0].AutoDescription)

	assert.Equal(t, "animalLocations[1]", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{\"cats\":[],\"place\":\"work\"}`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].AutoDefault)
	assert.Equal(t, "", valuesRows[1].Description)
	assert.Equal(t, "place with the fewest cats", valuesRows[1].AutoDescription)
}

func TestAutoMultilineDescription(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  # -- The best kind of animal probably, allow me to list their many varied benefits.
  # Cats are very funny, and quite friendly, in almost all cases
  # @default -- The list of cats that _I_ own
  cats:
      - echo
      - foxtrot
`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.cats", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "The list of cats that _I_ own", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)
	assert.Equal(t, "The best kind of animal probably, allow me to list their many varied benefits. Cats are very funny, and quite friendly, in almost all cases", valuesRows[0].AutoDescription)
}

func TestAutoMultilineDescriptionWithoutValue(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  # -- (list) I mean, dogs are quite nice too...
  # @default -- The list of dogs that _I_ own
  dogs:
`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.dogs", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "The list of dogs that _I_ own", valuesRows[0].AutoDefault)
	assert.Equal(t, "", valuesRows[0].Default)
	assert.Equal(t, "I mean, dogs are quite nice too...", valuesRows[0].Description)
}

func TestExtractValueNotationType(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  # -- (list) My animals lists
  # @notationType -- yaml
  cats:
   - mike
   - ralph
  # -- (list) My animal lists, but in tpl string
  # @notationType -- tpl
  catsInTpl: |
   {{- .Values.animals.cats }}

  # -- (object) Declaring object as tpl (to be cascaded with tpl function)
  # @notationType -- tpl
  dinosaur: |
    name: hockney
    dynamicVar: {{ .Values.fromOtherProperty }}

  # -- (object) Declaring object as yaml
  # @notationType -- yaml
  fish:
    name: nomoby
`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "animals.cats", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, yamlType, valuesRows[0].NotationType)
	assert.Equal(t, "- mike\n- ralph\n", valuesRows[0].Default)
	assert.Equal(t, "My animals lists", valuesRows[0].AutoDescription)

	assert.Equal(t, "animals.catsInTpl", valuesRows[1].Key)
	assert.Equal(t, listType, valuesRows[1].Type)
	assert.Equal(t, tplType, valuesRows[1].NotationType)
	assert.Equal(t, "{{- .Values.animals.cats }}\n", valuesRows[1].Default)
	assert.Equal(t, "My animal lists, but in tpl string", valuesRows[1].AutoDescription)

	assert.Equal(t, "animals.dinosaur", valuesRows[2].Key)
	assert.Equal(t, objectType, valuesRows[2].Type)
	assert.Equal(t, tplType, valuesRows[2].NotationType)
	assert.Equal(t, "name: hockney\ndynamicVar: {{ .Values.fromOtherProperty }}\n", valuesRows[2].Default)
	assert.Equal(t, "Declaring object as tpl (to be cascaded with tpl function)", valuesRows[2].AutoDescription)

	assert.Equal(t, "animals.fish", valuesRows[3].Key)
	assert.Equal(t, objectType, valuesRows[3].Type)
	assert.Equal(t, yamlType, valuesRows[3].NotationType)
	assert.Equal(t, "name: nomoby\n", valuesRows[3].Default)
	assert.Equal(t, "My animals lists", valuesRows[0].AutoDescription)
}

func TestExtractCustomDeclaredType(t *testing.T) {
	helmValues := parseYamlValues(`
animals:
  # -- (list/csv) My animals lists but annotated as csv field
  cats: mike,ralph

owner:
  # -- (string/email) This has to be email address
  # @notationType -- email
  email: "owner@home.org"
`)

	valuesRows, err := getSortedValuesTableRows(helmValues, make(map[string]helm.ChartValueDescription))

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "animals.cats", valuesRows[0].Key)
	// With custom value type, we can convey to the reader that this value is a list, but in a csv format
	assert.Equal(t, "list/csv", valuesRows[0].Type)
	assert.Equal(t, "`\"mike,ralph\"`", valuesRows[0].Default)
	assert.Equal(t, "My animals lists but annotated as csv field", valuesRows[0].AutoDescription)

	assert.Equal(t, "owner.email", valuesRows[1].Key)
	assert.Equal(t, "string/email", valuesRows[1].Type)
	assert.Equal(t, "email", valuesRows[1].NotationType)
	// In case of custom notation type, value in Default must be raw string
	// So that template can handle the formatting.
	// In this case, email might be reformatted as <a href="mailto:owner@home.org">owner@home.org</a>
	assert.Equal(t, "owner@home.org", valuesRows[1].Default)
	assert.Equal(t, "This has to be email address", valuesRows[1].AutoDescription)
}

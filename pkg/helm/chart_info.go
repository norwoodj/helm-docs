package helm

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var columnDescriptionRegex = regexp.MustCompile("# (.*) -- (.*)")
var commentContinuationRegex = regexp.MustCompile("#(.*)")

type ChartMetaMaintainer struct {
	Email string
	Name  string
}

type ChartMeta struct {
	ApiVersion  string `yaml:"apiVersion"`
	Name        string
	Description string
	Version     string
	Home        string
	Type        string
	Sources     []string
	Engine      string
	Maintainers []ChartMetaMaintainer
}

type ChartRequirementsItem struct {
	Name       string
	Version    string
	Repository string
}

type ChartRequirements struct {
	Dependencies []ChartRequirementsItem
}

type ChartDocumentationInfo struct {
	ChartMeta
	ChartRequirements

	ChartDirectory          string
	ChartValues             map[interface{}]interface{}
	ChartValuesDescriptions []ChartValueDescription
}

type ChartValueDescription struct {
	Key     string
	Columns map[string]string
}

func (desc *ChartValueDescription) Description() string {
	return desc.Columns["Description"]
}

func getYamlFileContents(filename string) ([]byte, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	yamlFileContents, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return []byte(yamlFileContents), nil
}

func yamlLoadAndCheck(yamlFileContents []byte, out interface{}) {
	err := yaml.Unmarshal(yamlFileContents, out)

	if err != nil {
		panic(err)
	}
}

func isErrorInReadingNecessaryFile(filePath string, loadError error) bool {
	if loadError != nil {
		if os.IsNotExist(loadError) {
			log.Printf("Required chart file %s missing. Skipping documentation for chart", filePath)
			return true
		} else {
			log.Printf("Error occurred in reading chart file %s. Skipping documentation for chart", filePath)
			return true
		}
	}

	return false
}

func parseChartFile(chartDirectory string) (ChartMeta, error) {
	chartYamlPath := path.Join(chartDirectory, "Chart.yaml")
	chartMeta := ChartMeta{}
	yamlFileContents, err := getYamlFileContents(chartYamlPath)

	if isErrorInReadingNecessaryFile(chartYamlPath, err) {
		return chartMeta, err
	}

	yamlLoadAndCheck(yamlFileContents, &chartMeta)
	return chartMeta, nil
}

func requirementKey(requirement ChartRequirementsItem) string {
	return fmt.Sprintf("%s/%s", requirement.Repository, requirement.Name)
}

func parseChartRequirementsFile(chartDirectory string, apiVersion string) (ChartRequirements, error) {
	var requirementsPath string

	if apiVersion == "v1" {
		requirementsPath = path.Join(chartDirectory, "requirements.yaml")

		if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
			return ChartRequirements{Dependencies: []ChartRequirementsItem{}}, nil
		}
	} else {
		requirementsPath = path.Join(chartDirectory, "Chart.yaml")
	}

	chartRequirements := ChartRequirements{}
	yamlFileContents, err := getYamlFileContents(requirementsPath)

	if isErrorInReadingNecessaryFile(requirementsPath, err) {
		return chartRequirements, err
	}

	yamlLoadAndCheck(yamlFileContents, &chartRequirements)

	sort.Slice(chartRequirements.Dependencies[:], func(i, j int) bool {
		return requirementKey(chartRequirements.Dependencies[i]) < requirementKey(chartRequirements.Dependencies[j])
	})

	return chartRequirements, nil
}

func parseChartValuesFile(chartDirectory string) (map[interface{}]interface{}, error) {
	valuesPath := path.Join(chartDirectory, "values.yaml")
	values := make(map[interface{}]interface{})
	yamlFileContents, err := getYamlFileContents(valuesPath)

	if isErrorInReadingNecessaryFile(valuesPath, err) {
		return values, err
	}

	yamlLoadAndCheck(yamlFileContents, &values)
	return values, nil
}

func parseChartValuesFileComments(chartDirectory string) ([]ChartValueDescription, error) {
	valuesPath := path.Join(chartDirectory, "values.yaml")
	valuesFile, err := os.Open(valuesPath)

	if isErrorInReadingNecessaryFile(valuesPath, err) {
		return []ChartValueDescription{}, err
	}

	defer valuesFile.Close()

	chartValueDescriptions := []ChartValueDescription{}
	scanner := bufio.NewScanner(valuesFile)

	key := ""
	column := ""
	for scanner.Scan() {
		match := columnDescriptionRegex.FindStringSubmatch(scanner.Text())

		if len(match) > 2 {
			if key == "" {
				key = strings.TrimSpace(match[1])
				column = "Description"

				chartValueDescriptions = append(chartValueDescriptions, ChartValueDescription{
					Key:     key,
					Columns: map[string]string{column: strings.TrimSpace(match[2])},
				})
			} else {
				column = strings.TrimSpace(match[1])

				chartValueDescriptions[len(chartValueDescriptions)-1].Columns[column] = strings.TrimSpace(match[2])
			}
		} else {
			match = commentContinuationRegex.FindStringSubmatch(scanner.Text())

			if (len(match)) > 1 && key != "" {
				chartValueDescriptions[len(chartValueDescriptions)-1].Columns[column] =
					chartValueDescriptions[len(chartValueDescriptions)-1].Columns[column] +
						" " +
						strings.TrimSpace(match[1])
			} else {
				key = ""
				column = ""
			}
		}

	}

	return chartValueDescriptions, nil
}

func ParseChartInformation(chartDirectory string) (ChartDocumentationInfo, error) {
	var chartDocInfo ChartDocumentationInfo
	var err error

	chartDocInfo.ChartDirectory = chartDirectory
	chartDocInfo.ChartMeta, err = parseChartFile(chartDirectory)
	if err != nil {
		return chartDocInfo, err
	}

	chartDocInfo.ChartRequirements, err = parseChartRequirementsFile(chartDirectory, chartDocInfo.ApiVersion)
	if err != nil {
		return chartDocInfo, err
	}

	chartDocInfo.ChartValues, err = parseChartValuesFile(chartDirectory)
	if err != nil {
		return chartDocInfo, err
	}

	chartDocInfo.ChartValuesDescriptions, err = parseChartValuesFileComments(chartDirectory)
	if err != nil {
		return chartDocInfo, err
	}

	return chartDocInfo, nil
}

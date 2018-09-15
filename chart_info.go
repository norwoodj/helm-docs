package main

import (
    "bufio"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
    "os"
    "path"
    "regexp"
)

var valuesDescriptionRegex = regexp.MustCompile("# (.*) -- (.*)")

type ChartMetaMaintainer struct {
    Email string
    Name string
}

type ChartMeta struct {
    ApiVersion string `yaml:"apiVersion"`
    Name string
    Description string
    Version string
    Home string
    Sources []string
    Engine string
    Maintainers []ChartMetaMaintainer
}

type ChartRequirementsItem struct {
    Name string
    Version string
    Repository string
}

type ChartRequirements struct {
    Dependencies []ChartRequirementsItem
}

type ChartValues map[interface{}]interface{}


func getYamlFileContents(filename string, debug bool) ([]byte, error) {
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return nil, err
    }

    yamlFileContents, err := ioutil.ReadFile(filename)

    if err != nil {
        panic(err)
    }

    if debug {
        log.Printf("Reading %s file contents:\n%s\n", filename, yamlFileContents)
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

func parseChartFile(chartDirectory string, debug bool) (ChartMeta, error) {
    chartYamlPath := path.Join(chartDirectory, "Chart.yaml")
    chartMeta := ChartMeta{}
    yamlFileContents, err := getYamlFileContents(chartYamlPath, debug)

    if isErrorInReadingNecessaryFile(chartYamlPath, err) {
        return chartMeta, err
    }

    yamlLoadAndCheck(yamlFileContents, &chartMeta)
    return chartMeta, nil
}

func parseChartRequirementsFile(chartDirectory string, debug bool) (ChartRequirements, error) {
    requirementsPath := path.Join(chartDirectory, "requirements.yaml")
    if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
        return ChartRequirements{Dependencies: []ChartRequirementsItem{}}, nil
    }

    chartRequirements := ChartRequirements{}
    yamlFileContents, err := getYamlFileContents(requirementsPath, debug)

    if isErrorInReadingNecessaryFile(requirementsPath, err) {
        return chartRequirements, err
    }

    yamlLoadAndCheck(yamlFileContents, &chartRequirements)
    return chartRequirements, nil
}

func parseValuesFile(chartDirectory string, debug bool) (ChartValues, error) {
    valuesPath := path.Join(chartDirectory, "values.yaml")
    values := ChartValues{}
    yamlFileContents, err := getYamlFileContents(valuesPath, debug)

    if isErrorInReadingNecessaryFile(valuesPath, err) {
        return values, err
    }

    yamlLoadAndCheck(yamlFileContents, &values)
    return values, nil
}

func parseValuesFileComments(chartDirectory string) (map[string]string, error) {
    valuesPath := path.Join(chartDirectory, "values.yaml")
    valuesFile, err := os.Open(valuesPath)

    if isErrorInReadingNecessaryFile(valuesPath, err) {
        return map[string]string{}, err
    }

    defer valuesFile.Close()

    keyToDescriptions := make(map[string]string)
    scanner := bufio.NewScanner(valuesFile)

    for scanner.Scan() {
        match := valuesDescriptionRegex.FindStringSubmatch(scanner.Text())

        if len(match) > 2 {
            keyToDescriptions[match[1]] = match[2]
        }
    }

    return keyToDescriptions, nil
}

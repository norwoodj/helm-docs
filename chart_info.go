package main

import (
    "bufio"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
    "os"
    "regexp"
)

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


func getYamlFileContents(filename string, debug bool) []byte {
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        log.Fatalf("ERROR: No %s file exists!", filename)
    }

    yamlFileContents, err := ioutil.ReadFile(filename)

    if err != nil {
        panic(err)
    }

    if debug {
        log.Printf("Reading %s file contents:\n%s\n", filename, yamlFileContents)
    }

    return []byte(yamlFileContents)
}

func yamlLoadAndCheck(yamlFileContents []byte, out interface{}) {
    err := yaml.Unmarshal(yamlFileContents, out)

    if err != nil {
        panic(err)
    }
}

func parseValuesFile(debug bool) ChartValues {
    yamlFileContents := getYamlFileContents("values.yaml", debug)
    values := make(ChartValues)
    yamlLoadAndCheck(yamlFileContents, &values)
    return values
}

func parseChartFile(debug bool) ChartMeta {
    yamlFileContents := getYamlFileContents("Chart.yaml", debug)
    chartMeta := ChartMeta{}
    yamlLoadAndCheck(yamlFileContents, &chartMeta)
    return chartMeta
}

func parseChartRequirementsFile(debug bool) ChartRequirements {
    if _, err := os.Stat("requirements.yaml"); os.IsNotExist(err) {
        return ChartRequirements{Dependencies: []ChartRequirementsItem{}}
    }

    yamlFileContents := getYamlFileContents("requirements.yaml", debug)
    chartRequirements := ChartRequirements{}
    yamlLoadAndCheck(yamlFileContents, &chartRequirements)
    return chartRequirements
}

func parseValuesFileComments() map[string]string {
    valuesFile, err := os.Open("values.yaml")
    if err != nil {
        log.Fatal(err)
    }

    defer valuesFile.Close()

    keyToDescriptions := make(map[string]string)
    scanner := bufio.NewScanner(valuesFile)
    re := regexp.MustCompile("# (.*) -- (.*)")

    for scanner.Scan() {
        match := re.FindStringSubmatch(scanner.Text())

        if len(match) > 2 {
            keyToDescriptions[match[1]] = match[2]
        }
    }

    return keyToDescriptions
}

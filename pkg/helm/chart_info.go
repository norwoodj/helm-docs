package helm

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	valuesDescriptionRegex   = regexp.MustCompile("^\\s*#\\s*(.*)\\s+--\\s*(.*)$")
	rawDescriptionRegex      = regexp.MustCompile("^\\s*#\\s+@raw")
	commentContinuationRegex = regexp.MustCompile("^\\s*#(\\s?)(.*)$")
	defaultValueRegex        = regexp.MustCompile("^\\s*# @default -- (.*)$")
	valueTypeRegex           = regexp.MustCompile("^\\((.*?)\\)\\s*(.*)$")
	valueNotationTypeRegex   = regexp.MustCompile("^\\s*#\\s+@notationType\\s+--\\s+(.*)$")
	sectionRegex             = regexp.MustCompile("^\\s*# @section -- (.*)$")
	extraRegex               = regexp.MustCompile("^\\s*# @extra -- (.*)$")
	extraContinuationRegex   = regexp.MustCompile("^\\s*# ?(.*)$")
)

type ChartMetaMaintainer struct {
	Email string
	Name  string
	Url   string
}

type ChartMeta struct {
	ApiVersion  string `yaml:"apiVersion"`
	AppVersion  string `yaml:"appVersion"`
	KubeVersion string `yaml:"kubeVersion"`
	Name        string
	Deprecated  bool
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
	Alias      string
}

type ChartRequirements struct {
	Dependencies []ChartRequirementsItem
}

type ChartValueDescription struct {
	Description  string
	Default      string
	Section      string
	ValueType    string
	NotationType string
}

type ChartDocumentationInfo struct {
	ChartMeta
	ChartRequirements

	ChartDirectory          string
	ChartValues             *yaml.Node
	ChartValuesDescriptions map[string]ChartValueDescription
	Extras                  map[string]string
}

type ChartValuesDocumentationParsingConfig struct {
	StrictMode                 bool
	AllowedMissingValuePaths   []string
	AllowedMissingValueRegexps []*regexp.Regexp
}

func getYamlFileContents(filename string) ([]byte, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	yamlFileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return []byte(strings.Replace(string(yamlFileContents), "\r\n", "\n", -1)), nil
}

func isErrorInReadingNecessaryFile(filePath string, loadError error) bool {
	if loadError != nil {
		if os.IsNotExist(loadError) {
			log.Warnf("Required chart file %s missing. Skipping documentation for chart", filePath)
			return true
		} else {
			log.Warnf("Error occurred in reading chart file %s. Skipping documentation for chart", filePath)
			return true
		}
	}

	return false
}

func parseChartFile(chartDirectory string) (ChartMeta, error) {
	chartYamlPath := filepath.Join(chartDirectory, "Chart.yaml")
	chartMeta := ChartMeta{}
	yamlFileContents, err := getYamlFileContents(chartYamlPath)

	if isErrorInReadingNecessaryFile(chartYamlPath, err) {
		return chartMeta, err
	}

	err = yaml.Unmarshal(yamlFileContents, &chartMeta)
	return chartMeta, err
}

func requirementKey(requirement ChartRequirementsItem) string {
	return fmt.Sprintf("%s/%s", requirement.Repository, requirement.Name)
}

func parseChartRequirementsFile(chartDirectory string, apiVersion string) (ChartRequirements, error) {
	var requirementsPath string

	if apiVersion == "v1" {
		requirementsPath = filepath.Join(chartDirectory, "requirements.yaml")

		if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
			return ChartRequirements{Dependencies: []ChartRequirementsItem{}}, nil
		}
	} else {
		requirementsPath = filepath.Join(chartDirectory, "Chart.yaml")
	}

	chartRequirements := ChartRequirements{}
	yamlFileContents, err := getYamlFileContents(requirementsPath)

	if isErrorInReadingNecessaryFile(requirementsPath, err) {
		return chartRequirements, err
	}

	err = yaml.Unmarshal(yamlFileContents, &chartRequirements)
	if err != nil {
		return chartRequirements, err
	}

	sort.Slice(chartRequirements.Dependencies[:], func(i, j int) bool {
		return requirementKey(chartRequirements.Dependencies[i]) < requirementKey(chartRequirements.Dependencies[j])
	})

	return chartRequirements, nil
}

func removeIgnored(rootNode *yaml.Node, parentKind yaml.Kind) {
	newContent := make([]*yaml.Node, 0, len(rootNode.Content))
	for i := 0; i < len(rootNode.Content); i++ {
		node := rootNode.Content[i]
		if !strings.Contains(node.HeadComment, "@ignore") {
			removeIgnored(node, node.Kind)
			newContent = append(newContent, node)
		} else if parentKind == yaml.MappingNode {
			// for parentKind each yaml key is represented by two nodes
			i++
		}
	}
	rootNode.Content = newContent
}

func parseChartValuesFile(chartDirectory string) (yaml.Node, error) {
	valuesPath := filepath.Join(chartDirectory, viper.GetString("values-file"))
	yamlFileContents, err := getYamlFileContents(valuesPath)

	var values yaml.Node
	if isErrorInReadingNecessaryFile(valuesPath, err) {
		return values, err
	}

	err = yaml.Unmarshal(yamlFileContents, &values)
	removeIgnored(&values, values.Kind)
	return values, err
}

func checkDocumentation(rootNode *yaml.Node, comments map[string]ChartValueDescription, config ChartValuesDocumentationParsingConfig) error {
	if len(rootNode.Content) == 0 {
		return nil
	}
	valuesWithoutDocs := collectValuesWithoutDoc(rootNode.Content[0], comments, make([]string, 0))
	valuesWithoutDocsAfterIgnore := make([]string, 0)
	for _, valueWithoutDoc := range valuesWithoutDocs {
		ignored := false
		for _, ignorableValuePath := range config.AllowedMissingValuePaths {
			ignored = ignored || valueWithoutDoc == ignorableValuePath
		}
		for _, ignorableValueRegexp := range config.AllowedMissingValueRegexps {
			ignored = ignored || ignorableValueRegexp.MatchString(valueWithoutDoc)
		}
		if !ignored {
			valuesWithoutDocsAfterIgnore = append(valuesWithoutDocsAfterIgnore, valueWithoutDoc)
		}
	}
	if len(valuesWithoutDocsAfterIgnore) > 0 {
		return fmt.Errorf("values without documentation: \n%s", strings.Join(valuesWithoutDocsAfterIgnore, "\n"))
	}
	return nil
}

func collectValuesWithoutDoc(node *yaml.Node, comments map[string]ChartValueDescription, currentPath []string) []string {
	valuesWithoutDocs := make([]string, 0)
	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode, valueNode := node.Content[i], node.Content[i+1]
			currentPath = append(currentPath, keyNode.Value)
			pathString := strings.Join(currentPath, ".")
			if _, ok := comments[pathString]; !ok {
				valuesWithoutDocs = append(valuesWithoutDocs, pathString)
			}

			childValuesWithoutDoc := collectValuesWithoutDoc(valueNode, comments, currentPath)
			valuesWithoutDocs = append(valuesWithoutDocs, childValuesWithoutDoc...)

			currentPath = currentPath[:len(currentPath)-1]
		}
	case yaml.SequenceNode:
		for i := 0; i < len(node.Content); i++ {
			valueNode := node.Content[i]
			currentPath = append(currentPath, fmt.Sprintf("[%d]", i))
			childValuesWithoutDoc := collectValuesWithoutDoc(valueNode, comments, currentPath)
			valuesWithoutDocs = append(valuesWithoutDocs, childValuesWithoutDoc...)
			currentPath = currentPath[:len(currentPath)-1]
		}
	}
	return valuesWithoutDocs
}

func parseChartValuesFileComments(chartDirectory string, values *yaml.Node, lintingConfig ChartValuesDocumentationParsingConfig) (map[string]ChartValueDescription, map[string]string, error) {
	valuesPath := filepath.Join(chartDirectory, viper.GetString("values-file"))
	valuesFile, err := os.Open(valuesPath)

	if isErrorInReadingNecessaryFile(valuesPath, err) {
		return map[string]ChartValueDescription{}, map[string]string{}, err
	}

	defer valuesFile.Close()

	extras := make(map[string]string)
	currentExtra := ""
	keyToDescriptions := make(map[string]ChartValueDescription)
	scanner := bufio.NewScanner(valuesFile)
	foundValuesComment := false
	commentLines := make([]string, 0)
	currentLineIdx := -1

	for scanner.Scan() {
		currentLineIdx++
		currentLine := scanner.Text()

		if currentExtra != "" {
			match := extraContinuationRegex.FindStringSubmatch(currentLine)
			if len(match) > 1 {
				commentLines = append(commentLines, match[1])
				continue
			}
			extras[currentExtra] = strings.Join(commentLines, "\n")
			commentLines = make([]string, 0)
			currentExtra = ""
		} else {
			match := extraRegex.FindStringSubmatch(currentLine)
			if len(match) == 2 && match[1] != "" {
				currentExtra = match[1]
				continue
			}
		}

		// If we've not yet found a values comment with a key name, try and find one on each line
		if !foundValuesComment {
			match := valuesDescriptionRegex.FindStringSubmatch(currentLine)
			if len(match) < 3 || match[1] == "" {
				continue
			}
			foundValuesComment = true
			commentLines = append(commentLines, currentLine)
			continue
		}

		// If we've already found a values comment, on the next line try and parse a comment continuation, a custom default value, or a section comment.
		// If we find continuations we can add them to the list and continue to the next line until we find a section comment or default value.
		// If we find a default value, we can add it to the list and continue to the next line. In the case we don't find one, we continue looking for a section comment.
		// When we eventually find a section comment, we add it to the list and conclude matching for the current key. If we don't find one, matching is also concluded.
		//
		// NOTE: This isn't readily enforced yet, because we can match the section comment and custom default value more than once and in another order, although this is just overwriting it.
		// Values comment, possible continuation, default value once or none then section comment once or none should be the preferred order.
		defaultCommentMatch := defaultValueRegex.FindStringSubmatch(currentLine)
		sectionCommentMatch := sectionRegex.FindStringSubmatch(currentLine)
		commentContinuationMatch := commentContinuationRegex.FindStringSubmatch(currentLine)

		if len(defaultCommentMatch) > 1 || len(sectionCommentMatch) > 1 || len(commentContinuationMatch) > 1 {
			commentLines = append(commentLines, currentLine)
			continue
		}

		// If we haven't continued by this point, we didn't match any of the comment formats we want, so we need to add
		// the in progress value to the map, and reset to looking for a new key
		key, description := ParseComment(commentLines)
		if key != "" {
			keyToDescriptions[key] = description
		}

		commentLines = make([]string, 0)
		foundValuesComment = false
	}

	if currentExtra != "" {
		extras[currentExtra] = strings.Join(commentLines, "\n")
	}

	if lintingConfig.StrictMode {
		err := checkDocumentation(values, keyToDescriptions, lintingConfig)
		if err != nil {
			return nil, nil, err
		}
	}
	return keyToDescriptions, extras, nil
}

func ParseChartInformation(chartDirectory string, documentationParsingConfig ChartValuesDocumentationParsingConfig) (ChartDocumentationInfo, error) {
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

	chartValues, err := parseChartValuesFile(chartDirectory)
	if err != nil {
		return chartDocInfo, err
	}

	chartDocInfo.ChartValues = &chartValues
	chartDocInfo.ChartValuesDescriptions, chartDocInfo.Extras, err = parseChartValuesFileComments(chartDirectory, &chartValues, documentationParsingConfig)
	if err != nil {
		return chartDocInfo, err
	}

	return chartDocInfo, nil
}

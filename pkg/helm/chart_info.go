package helm

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/nlepage/go-tarfs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var valuesDescriptionRegex = regexp.MustCompile(`^\s*#\s*(.*)\s+--\s*(.*)$`)
var rawDescriptionRegex = regexp.MustCompile(`^\s*#\s+@raw`)
var commentContinuationRegex = regexp.MustCompile(`^\s*#(\s?)(.*)$`)
var defaultValueRegex = regexp.MustCompile(`^\s*# @default -- (.*)$`)
var valueTypeRegex = regexp.MustCompile(`^\((.*?)\)\s*(.*)$`)
var valueNotationTypeRegex = regexp.MustCompile(`^\s*#\s+@notationType\s+--\s+(.*)$`)

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
	ValueType    string
	NotationType string
}

type ChartDocumentationInfo struct {
	ChartMeta
	ChartRequirements

	ChartDirectory          string
	ChartValues             *yaml.Node
	ChartValuesDescriptions map[string]ChartValueDescription

	ChartFS   fs.FS
	CloseFunc func() error // closes FS if necessary
}

func parseChartFS(chartDirectory string) (ChartDocumentationInfo, error) {
	chartDocInfo := ChartDocumentationInfo{
		ChartDirectory: chartDirectory,
	}

	if strings.HasSuffix(chartDirectory, ".tgz") {
		tf, err := os.Open(chartDirectory)
		if err != nil {
			return chartDocInfo, fmt.Errorf("could not open Chart archive %s: %w", chartDirectory, err)
		}
		chartDocInfo.CloseFunc = tf.Close

		tfs, err := tarfs.New(tf)
		if err != nil {
			return chartDocInfo, fmt.Errorf("could not open Chart archive %s: %w", chartDirectory, err)
		}

		// Do not need to check these errors any further,
		// they have already been checked by
		// pkg/helm/chart_finder.go:checkArchiveIsChart
		dirs, _ := fs.ReadDir(tfs, ".")
		chartDocInfo.ChartFS, err = fs.Sub(tfs, dirs[0].Name())
		if err != nil {
			return chartDocInfo, fmt.Errorf("could not read Chart archive %s: %w", chartDirectory, err)
		}
	} else {
		chartDocInfo.ChartFS = os.DirFS(chartDirectory)
	}

	return chartDocInfo, nil
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

	yamlFileContents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return []byte(strings.Replace(string(yamlFileContents), "\r\n", "\n", -1)), nil
}

func getYamlFileContentsFS(fsys fs.FS, filename string) ([]byte, error) {
	if _, err := fs.Stat(fsys, filename); errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	yamlFileContents, err := fs.ReadFile(fsys, filename)
	if err != nil {
		panic(err)
	}

	return []byte(strings.Replace(string(yamlFileContents), "\r\n", "\n", -1)), nil
}

func isErrorInReadingNecessaryFile(filePath string, loadError error) bool {
	if loadError != nil {
		if os.IsNotExist(loadError) {
			log.Warn(loadError)
			log.Warnf("Required chart file %s missing. Skipping documentation for chart", filePath)
			return true
		} else {
			log.Warn(loadError)
			log.Warnf("Error occurred in reading chart file %s. Skipping documentation for chart", filePath)
			return true
		}
	}

	return false
}

func (chartDocInfo ChartDocumentationInfo) parseChartFile() error {
	chartYaml := "Chart.yaml"
	chartYamlPath := filepath.Join(chartDocInfo.ChartDirectory, "Chart.yaml")
	chartDocInfo.ChartMeta = ChartMeta{}
	yamlFileContents, err := getYamlFileContents(chartYaml)
	if isErrorInReadingNecessaryFile(chartYamlPath, err) {
		return err
	}

	err = yaml.Unmarshal(yamlFileContents, &chartDocInfo.ChartMeta)
	return err
}

func requirementKey(requirement ChartRequirementsItem) string {
	return fmt.Sprintf("%s/%s", requirement.Repository, requirement.Name)
}

func (chartDocInfo ChartDocumentationInfo) parseChartRequirementsFile() error {
	var requirementsYaml string
	var requirementsYamlPath string

	if chartDocInfo.ApiVersion == "v1" {
		requirementsYaml = "requirements.yaml"
		requirementsYamlPath = filepath.Join(chartDocInfo.ChartDirectory, requirementsYaml)

		if _, err := fs.Stat(chartDocInfo.ChartFS, requirementsYaml); os.IsNotExist(err) {
			chartDocInfo.ChartRequirements = ChartRequirements{Dependencies: []ChartRequirementsItem{}}
			return nil
		}
	} else {
		requirementsYaml = "Chart.yaml"
		requirementsYamlPath = filepath.Join(chartDocInfo.ChartDirectory, requirementsYaml)
	}

	chartDocInfo.ChartRequirements = ChartRequirements{}
	yamlFileContents, err := getYamlFileContentsFS(chartDocInfo.ChartFS, requirementsYaml)

	if isErrorInReadingNecessaryFile(requirementsYamlPath, err) {
		return err
	}

	err = yaml.Unmarshal(yamlFileContents, &chartDocInfo.ChartRequirements)
	if err != nil {
		return err
	}

	sort.Slice(chartDocInfo.ChartRequirements.Dependencies[:], func(i, j int) bool {
		return requirementKey(chartDocInfo.ChartRequirements.Dependencies[i]) < requirementKey(chartDocInfo.ChartRequirements.Dependencies[j])
	})

	return nil
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

func (chartDocInfo ChartDocumentationInfo) parseChartValuesFile() error {
	valuesFile := viper.GetString("values-file")
	valuesFilePath := filepath.Join(chartDocInfo.ChartDirectory, valuesFile)
	yamlFileContents, err := getYamlFileContentsFS(chartDocInfo.ChartFS, valuesFile)
	if isErrorInReadingNecessaryFile(valuesFilePath, err) {
		return err
	}

	chartDocInfo.ChartValues = &yaml.Node{}
	err = yaml.Unmarshal(yamlFileContents, chartDocInfo.ChartValues)
	if err != nil {
		return err
	}

	removeIgnored(chartDocInfo.ChartValues, chartDocInfo.ChartValues.Kind)

	return nil
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

func (chartDocInfo ChartDocumentationInfo) parseChartValuesFileComments(lintingConfig ChartValuesDocumentationParsingConfig) error {
	valuesFile := viper.GetString("values-file")
	valuesFilePath := filepath.Join(chartDocInfo.ChartDirectory, valuesFile)
	values, err := chartDocInfo.ChartFS.Open(valuesFile)
	if isErrorInReadingNecessaryFile(valuesFilePath, err) {
		return err
	}

	defer values.Close()

	chartDocInfo.ChartValuesDescriptions = make(map[string]ChartValueDescription)
	scanner := bufio.NewScanner(values)
	foundValuesComment := false
	commentLines := make([]string, 0)
	currentLineIdx := -1

	for scanner.Scan() {
		currentLineIdx++
		currentLine := scanner.Text()

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

		// If we've already found a values comment, on the next line try and parse a custom default value. If we find one
		// that completes parsing for this key, add it to the list and reset to searching for a new key
		defaultCommentMatch := defaultValueRegex.FindStringSubmatch(currentLine)
		commentContinuationMatch := commentContinuationRegex.FindStringSubmatch(currentLine)

		if len(defaultCommentMatch) > 1 || len(commentContinuationMatch) > 1 {
			commentLines = append(commentLines, currentLine)
			continue
		}

		// If we haven't continued by this point, we didn't match any of the comment formats we want, so we need to add
		// the in progress value to the map, and reset to looking for a new key
		key, description := ParseComment(commentLines)
		chartDocInfo.ChartValuesDescriptions[key] = description
		commentLines = make([]string, 0)
		foundValuesComment = false
	}
	if lintingConfig.StrictMode {
		err := checkDocumentation(chartDocInfo.ChartValues, chartDocInfo.ChartValuesDescriptions, lintingConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseChartInformation(chartDirectory string, documentationParsingConfig ChartValuesDocumentationParsingConfig) (ChartDocumentationInfo, error) {
	chartDocInfo, err := parseChartFS(chartDirectory)
	if err != nil {
		return chartDocInfo, err
	}

	// chartDocInfo.ChartMeta, err = parseChartFile(chartDirectory)
	// if err != nil {
	// 	return chartDocInfo, err
	// }
	err = chartDocInfo.parseChartFile()
	if err != nil {
		return chartDocInfo, err
	}

	// chartDocInfo.ChartRequirements, err = parseChartRequirementsFile(chartDirectory, chartDocInfo.ApiVersion)
	// if err != nil {
	// 	return chartDocInfo, err
	// }
	err = chartDocInfo.parseChartRequirementsFile()
	if err != nil {
		return chartDocInfo, err
	}

	// chartValues, err := parseChartValuesFile(chartDirectory)
	// if err != nil {
	// 	return chartDocInfo, err
	// }
	// chartDocInfo.ChartValues = &chartValues
	err = chartDocInfo.parseChartValuesFile()
	if err != nil {
		return chartDocInfo, err
	}

	// chartDocInfo.ChartValuesDescriptions, err = parseChartValuesFileComments(chartDirectory, &chartValues, documentationParsingConfig)
	// if err != nil {
	// 	return chartDocInfo, err
	// }
	err = chartDocInfo.parseChartValuesFileComments(documentationParsingConfig)
	if err != nil {
		return chartDocInfo, err
	}

	return chartDocInfo, nil
}

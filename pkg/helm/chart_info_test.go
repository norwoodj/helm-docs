package helm_test

import (
	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

type ChartParsingTestSuite struct {
	suite.Suite
}

func (_ *ChartParsingTestSuite) SetupTest() {
	viper.Set("values-file", "values.yaml")
}

func (suite *ChartParsingTestSuite) resolveRelativePath(chartPath string) string {
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		suite.T().Fatal("PROJECT_ROOT environment variable is not set")
	}
	return filepath.Join(projectRoot, chartPath)
}
func TestChartParsingTestSuite(t *testing.T) {
	suite.Run(t, new(ChartParsingTestSuite))
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOff() {
	chartPath := suite.resolveRelativePath("example-charts/full-template/")
	_, error := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: false,
	})
	suite.NoError(error)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOn() {
	chartPath := suite.resolveRelativePath("example-charts/full-template/")
	_, error := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
	})
	expectedError := `values without documentation: 
controller
controller.name
controller.image
controller.image.repository
controller.image.tag
controller.extraVolumes
controller.extraVolumes.[0].name
controller.extraVolumes.[0].configMap
controller.extraVolumes.[0].configMap.name
controller.publishService
controller.service
controller.service.annotations
controller.service.annotations.external-dns.alpha.kubernetes.io/hostname
controller.service.type`
	suite.EqualError(error, expectedError)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOnIgnores() {
	chartPath := suite.resolveRelativePath("example-charts/full-template/")
	_, error := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
		AllowedMissingValuePaths: []string{
			"controller",
			"controller.image",
			"controller.name",
			"controller.image.repository",
			"controller.image.tag",
			"controller.extraVolumes",
			"controller.extraVolumes.[0].name",
			"controller.extraVolumes.[0].configMap",
			"controller.extraVolumes.[0].configMap.name",
			"controller.publishService",
			"controller.service",
			"controller.service.annotations",
			"controller.service.annotations.external-dns.alpha.kubernetes.io/hostname",
			"controller.service.type",
		},
	})
	suite.NoError(error)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOnIgnoresRegexp() {
	chartPath := suite.resolveRelativePath("example-charts/full-template/")
	_, error := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
		AllowedMissingValueRegexps: []*regexp.Regexp{
			regexp.MustCompile("controller.*"),
		},
	})
	suite.NoError(error)
}

func (suite *ChartParsingTestSuite) TestFullyDocumentedChartStrictModeOn() {
	chartPath := suite.resolveRelativePath("example-charts/fully-documented/")
	_, error := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
	})
	suite.NoError(error)
}

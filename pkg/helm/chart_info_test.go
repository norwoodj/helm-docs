package helm_test

import (
	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
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

func TestChartParsingTestSuite(t *testing.T) {
	suite.Run(t, new(ChartParsingTestSuite))
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOff() {
	chartPath := filepath.Join("test-fixtures", "full-template")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: false,
	})
	suite.NoError(err)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOn() {
	chartPath := filepath.Join("test-fixtures", "full-template")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
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
	suite.EqualError(err, expectedError)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOnIgnores() {
	chartPath := filepath.Join("test-fixtures", "full-template")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
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
	suite.NoError(err)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOnIgnoresRegexp() {
	chartPath := filepath.Join("test-fixtures", "full-template")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
		AllowedMissingValueRegexps: []*regexp.Regexp{
			regexp.MustCompile("controller.*"),
		},
	})
	suite.NoError(err)
}

func (suite *ChartParsingTestSuite) TestFullyDocumentedChartStrictModeOn() {
	chartPath := filepath.Join("test-fixtures", "fully-documented")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
	})
	suite.NoError(err)
}

func (suite *ChartParsingTestSuite) TestExtraSectionsWithDescriptionAndDefault() {
	chartPath := filepath.Join("test-fixtures", "extra-sections")
	info, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{})
	suite.Require().NoError(err)

	pvc, ok := info.ChartValuesDescriptions["controller.persistentVolumeClaims"]
	suite.Require().True(ok)
	suite.Equal("List of persistent volume claims to create. For very long comments, break them into multiple lines.", pvc.Description)
	suite.Equal("the chart will construct this list internally unless specified", pvc.Default)

	replicas, ok := info.ChartValuesDescriptions["controller.replicas"]
	suite.Require().True(ok)
	suite.Equal("Number of replicas.", replicas.Description)

	suite.Equal(map[string]string{
		"Volume Claim Template": "The volume claim template is used to create a `PersistentVolumeClaim`.\n\n```\n( N_enclaves * 200 ) + ( N_keys * 200 ) = total number of bytes\n```",
		"First Extra":           "First line.",
		"Second Extra":          "Second line.",
	}, info.Extras)
}

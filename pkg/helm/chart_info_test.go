package helm_test

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
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

func (suite *ChartParsingTestSuite) TestLockFileForChartWithMultipleDeps() {
	chartPath := filepath.Join("test-fixtures", "helm-3")
	want := []helm.ChartRequirementsItem{
		{
			Name:       "airflow",
			Version:    "1.17.0",
			Repository: "https://airflow.apache.org",
			Alias:      "nginx-but-actually-airflow",
			Constraint: "~1.17.0",
		},
		{
			Name:       "nginx-ingress",
			Version:    "0.22.1",
			Repository: "https://charts.helm.sh/stable",
			Alias:      "",
			Constraint: "",
		},
		{
			Name:       "nginx-ingress",
			Version:    "0.22.1",
			Repository: "https://charts.helm.sh/stable",
			Alias:      "nginx-2",
			Constraint: "~0.22.1",
		},
	}
	got, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: false,
	})
	suite.Equal(want, got.ChartRequirements.Dependencies)
	suite.NoError(err)
}

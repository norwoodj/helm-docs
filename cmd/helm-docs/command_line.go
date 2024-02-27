package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/norwoodj/helm-docs/pkg/document"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version string

func possibleLogLevels() []string {
	levels := make([]string, 0)

	for _, l := range log.AllLevels {
		levels = append(levels, l.String())
	}

	return levels
}

func initializeCli() {
	logLevelName := viper.GetString("log-level")
	logLevel, err := log.ParseLevel(logLevelName)
	if err != nil {
		log.Errorf("Failed to parse provided log level %s: %s", logLevelName, err)
		os.Exit(1)
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(logLevel)
}

func newHelmDocsCommand(run func(cmd *cobra.Command, args []string)) (*cobra.Command, error) {
	command := &cobra.Command{
		Use:     "helm-docs",
		Short:   "helm-docs automatically generates markdown documentation for helm charts from requirements and values files",
		Version: version,
		Run:     run,
	}

	logLevelUsage := fmt.Sprintf("Level of logs that should printed, one of (%s)", strings.Join(possibleLogLevels(), ", "))
	command.PersistentFlags().StringP("chart-search-root", "c", ".", "directory to search recursively within for charts")
	command.PersistentFlags().BoolP("dry-run", "d", false, "don't actually render any markdown files just print to stdout passed")
	command.PersistentFlags().Bool("ignore-non-descriptions", false, "ignore values without a comment, this values will not be included in the README")
	command.PersistentFlags().StringP("ignore-file", "i", ".helmdocsignore", "The filename to use as an ignore file to exclude chart directories")
	command.PersistentFlags().StringP("log-level", "l", "info", logLevelUsage)
	command.PersistentFlags().StringP("output-file", "o", "README.md", "markdown file path relative to each chart directory to which rendered documentation will be written")
	command.PersistentFlags().StringP("sort-values-order", "s", document.AlphaNumSortOrder, fmt.Sprintf("order in which to sort the values table (\"%s\" or \"%s\")", document.AlphaNumSortOrder, document.FileSortOrder))
	command.PersistentFlags().StringSliceP("template-files", "t", []string{"README.md.gotmpl"}, "gotemplate file paths relative to each chart directory from which documentation will be generated")
	command.PersistentFlags().StringP("badge-style", "b", "flat-square", "badge style to use for charts")
	command.PersistentFlags().StringP("values-file", "f", "values.yaml", "Path to values file")
	command.PersistentFlags().BoolP("document-dependency-values", "u", false, "For charts with dependencies, include the dependency values in the chart values documentation")
	command.PersistentFlags().StringSliceP("chart-to-generate", "g", []string{}, "List of charts that will have documentation generated. Comma separated, no space. Empty list - generate for all charts in chart-search-root")
	command.PersistentFlags().BoolP("documentation-strict-mode", "x", false, "Fail the generation of docs if there are undocumented values")
	command.PersistentFlags().StringSliceP("documentation-strict-ignore-absent", "y", []string{"service.type", "image.repository", "image.tag"}, "A comma separate values which are allowed not to be documented in strict mode")
	command.PersistentFlags().StringSliceP("documentation-strict-ignore-absent-regex", "z", []string{".*service\\.type", ".*image\\.repository", ".*image\\.tag"}, "A comma separate values which are allowed not to be documented in strict mode")
	command.PersistentFlags().Bool("skip-version-footer", false, "if true the helm-docs version footer will not be shown in the default README template")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HELM_DOCS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err := viper.BindPFlags(command.PersistentFlags())

	return command, err
}

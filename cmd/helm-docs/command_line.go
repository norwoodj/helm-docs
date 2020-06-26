package main

import (
	"fmt"
	"os"
	"strings"

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
	command.PersistentFlags().BoolP("dry-run", "d", false, "don't actually render any markdown files just print to stdout passed")
	command.PersistentFlags().StringP("ignore-file", "i", ".helmdocsignore", "The filename to use as an ignore file to exclude chart directories")
	command.PersistentFlags().StringP("path", "p", ".", "The path to search for chart directories")
	command.PersistentFlags().BoolP("omit-blanks", "s", false, "Do not output fields with no description")
	command.PersistentFlags().StringP("log-level", "l", "info", logLevelUsage)
	command.PersistentFlags().StringP("output-file", "o", "README.md", "markdown file path relative to each chart directory to which rendered documentation will be written")
	command.PersistentFlags().StringP("template-file", "t", "README.md.gotmpl", "gotemplate file path relative to each chart directory from which documentation will be generated")
	command.PersistentFlags().BoolP("blank-container-defaults", "c", false, "For lists / dictionaries output blank default text instead of the whole object")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HELM_DOCS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err := viper.BindPFlags(command.PersistentFlags())

	return command, err
}

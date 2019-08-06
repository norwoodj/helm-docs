package main

import (
	"os"
	"strings"
	"sync"

	"github.com/norwoodj/helm-docs/pkg/document"
	"github.com/norwoodj/helm-docs/pkg/helm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func retrieveInfoAndPrintDocumentation(chartDirectory string, waitGroup *sync.WaitGroup, dryRun bool) {
	defer waitGroup.Done()
	chartDocumentationInfo, err := helm.ParseChartInformation(chartDirectory)

	if err != nil {
		log.Warnf("Error parsing information for chart %s, skipping: %s", chartDirectory, err)
		return
	}

	document.PrintDocumentation(chartDocumentationInfo, dryRun)

}

func helmDocs(_ *cobra.Command, _ []string) {
	initializeCli()

	ignoreDir := viper.GetString("ignore-dir")

	chartDirs, err := helm.FindChartDirectories(ignoreDir)

	if err != nil {
		log.Errorf("Error finding chart directories: %s", err)
		os.Exit(1)
	}

	log.Infof("Found Chart directories [%s]", strings.Join(chartDirs, ", "))
	dryRun := viper.GetBool("dry-run")
	waitGroup := sync.WaitGroup{}

	for _, c := range chartDirs {
		waitGroup.Add(1)

		// On dry runs all output goes to stdout, and so as to not jumble things, generate serially
		if dryRun {
			retrieveInfoAndPrintDocumentation(c, &waitGroup, dryRun)
		} else {
			go retrieveInfoAndPrintDocumentation(c, &waitGroup, dryRun)
		}
	}

	waitGroup.Wait()
}

func main() {
	command, err := newHelmDocsCommand(helmDocs)
	if err != nil {
		log.Errorf("Failed to create the CLI commander: %s", err)
		os.Exit(1)
	}

	if err := command.Execute(); err != nil {
		log.Errorf("Failed to start the CLI: %s", err)
		os.Exit(1)
	}
}

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

func retrieveInfoAndPrintDocumentation(chartDirectory string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	chartDocumentationInfo, err := helm.ParseChartInformation(chartDirectory)

	if err != nil {
		log.Warnf("Error parsing information for chart %s, skipping: %s", chartDirectory, err)
		return
	}

	document.PrintDocumentation(chartDocumentationInfo, viper.GetBool("dry-run"))

}

func helmDocs(_ *cobra.Command, _ []string) {
	initializeCli()
	chartDirs, err := helm.FindChartDirectories()

	if err != nil {
		log.Errorf("Error finding chart directories: %s", err)
		os.Exit(1)
	}

	log.Infof("Found Chart directories [%s]", strings.Join(chartDirs, ", "))
	waitGroup := sync.WaitGroup{}

	for _, c := range chartDirs {
		waitGroup.Add(1)
		go retrieveInfoAndPrintDocumentation(c, &waitGroup)
	}

	waitGroup.Wait()
}

func main() {
	command, err := newHelmDocsCommand(helmDocs)
	if err != nil {
		panic(err)
	}

	if err := command.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

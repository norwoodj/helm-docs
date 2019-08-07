package helm

import (
	"bufio"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/helm/pkg/ignore"
)

var defaultIgnore = map[string]bool{
	".git": true,
}

func FindChartDirectories() ([]string, error) {
	ignoreRules := ignore.Empty()
	ignoreFilename := viper.GetString("ignore-file")
	ignoreFile, err := os.Open(ignoreFilename)

	if err == nil {
		ignoreRules, err = ignore.Parse(bufio.NewReader(ignoreFile))

		if err != nil {
			log.Warnf("Failed to parse ignore rules from file %s", ignoreFilename)
			ignoreRules = ignore.Empty()
		}
	}

	chartDirs := make([]string, 0)
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && (ignoreRules.Ignore(path, info) || defaultIgnore[path]) {
			log.Debugf("Ignoring directory %s", path)
			return filepath.SkipDir
		}

		if filepath.Base(path) == "Chart.yaml" {
			chartDirs = append(chartDirs, filepath.Dir(path))
		}

		return nil
	})

	return chartDirs, err
}

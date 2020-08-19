package helm

import (
	"os"
	"path/filepath"

	"github.com/norwoodj/helm-docs/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func FindChartDirectories() ([]string, error) {
	ignoreFilename := viper.GetString("ignore-file")
	ignoreContext := util.NewIgnoreContext(ignoreFilename)
	chartDirs := make([]string, 0)
	searchPath := viper.GetString("search-path")

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		absolutePath, _ := filepath.Abs(path)

		if info.IsDir() && ignoreContext.ShouldIgnore(absolutePath, info) {
			log.Debugf("Ignoring directory %s", path)
			return filepath.SkipDir
		}

		if filepath.Base(path) == "Chart.yaml" {
			if ignoreContext.ShouldIgnore(absolutePath, info) {
				log.Debugf("Ignoring chart file %s", path)
				return nil
			}

			chartDirs = append(chartDirs, filepath.Dir(path))
		}

		return nil
	})

	return chartDirs, err
}

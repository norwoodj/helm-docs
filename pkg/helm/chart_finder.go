package helm

import (
	"os"
	"path/filepath"

	"github.com/norwoodj/helm-docs/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func FindChartDirectories(chartSearchRoot string) ([]string, error) {
	ignoreFilename := viper.GetString("ignore-file")
	ignoreContext := util.NewIgnoreContext(ignoreFilename)
	chartDirs := make([]string, 0)

	err := filepath.Walk(chartSearchRoot, func(path string, info os.FileInfo, err error) error {
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
			relativeChartDir, err := filepath.Rel(chartSearchRoot, filepath.Dir(path))

			if err != nil {
				return err
			}

			chartDirs = append(chartDirs, relativeChartDir)
		}

		return nil
	})

	return chartDirs, err
}

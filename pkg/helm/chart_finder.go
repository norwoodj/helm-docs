package helm

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nlepage/go-tarfs"
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

		basePath := filepath.Base(path)
		switch {
		case basePath == "Chart.yaml":
			if ignoreContext.ShouldIgnore(absolutePath, info) {
				log.Debugf("Ignoring chart file %s", path)
				return nil
			}
			relativeChartDir, err := filepath.Rel(chartSearchRoot, filepath.Dir(path))
			if err != nil {
				return err
			}

			chartDirs = append(chartDirs, relativeChartDir)
		case strings.HasSuffix(basePath, ".tgz"):
			isChart, err := checkArchiveIsChart(path)
			if err != nil {
				log.Warnf("Could not check archive %s: %s", path, err.Error())
				return nil
			}

			if isChart {
				if ignoreContext.ShouldIgnore(absolutePath, info) {
					log.Debugf("Ignoring chart archive %s", path)
					return nil
				}

				relativeChartDir, err := filepath.Rel(chartSearchRoot, filepath.Dir(path))
				if err != nil {
					return err
				}

				chartDirs = append(chartDirs, relativeChartDir)
			}
		}

		return nil
	})

	return chartDirs, err
}

func checkArchiveIsChart(path string) (bool, error) {
	tf, err := os.Open(path)
	if err != nil {
		defer tf.Close()
		return false, fmt.Errorf("could not open archive %s: %w", path, err)
	}
	defer tf.Close()

	tfs, err := tarfs.New(tf)
	if err != nil {
		return false, fmt.Errorf("could not open archive %s: %w", path, err)
	}

	dirs, err := fs.ReadDir(tfs, ".")
	if errors.Is(err, fs.ErrNotExist) {
		// Chart.yaml does not exist in the chart archive
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("could not check archive %s: %w", path, err)
	}

	// Chart archives should contain a single directory AFAIK
	if len(dirs) != 1 || !dirs[0].IsDir() {
		return false, nil
	}

	// dirs[0].Name() will be the chart name
	_, err = fs.Stat(tfs, filepath.Join(dirs[0].Name(), "Chart.yaml"))
	if errors.Is(err, fs.ErrNotExist) {
		// Chart.yaml does not exist in the chart archive
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("could not check archive %s: %w", path, err)
	}

	return true, nil
}

package main

import (
    "os"
    "path/filepath"
    "strings"
)

var ignoreDirectories = map[string]bool{
    ".git": true,
    "templates": true,
}

func findChartDirectories() []string {
    var chartDirs []string

    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() && ignoreDirectories[path] {
            return filepath.SkipDir
        }

        if strings.HasSuffix(path, "Chart.yaml") {
            chartDirs = append(chartDirs, filepath.Dir(path))
        }

        return nil
    })

    return chartDirs
}

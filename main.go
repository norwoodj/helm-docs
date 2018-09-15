package main

import (
    "log"
    "strings"
    "sync"
)


func main() {
    args := parseCommandLine()
    chartDirs := findChartDirectories()
    log.Printf("Found Chart directories [%s]", strings.Join(chartDirs, ", "))
    waitGroup := sync.WaitGroup{}

    for _, c := range chartDirs {
        waitGroup.Add(1)
        go printDocumentation(c, args.debug, args.dryRun, &waitGroup)
    }

    waitGroup.Wait()
}

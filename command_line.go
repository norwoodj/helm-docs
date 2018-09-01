package main

import (
    "fmt"
    "os"
    flag "github.com/spf13/pflag"
)

type HelmDocArgs struct {
    help bool
    debug bool
    dryRun bool
}


const DEBUG_HELP = "Print debug output"
const DRY_RUN_HELP = "Don't actually render any markdown files with docs, juts print to stdout"
const HELP_HELP = "Print this help menu, then exist"


func printHelp() {
    fmt.Println("Usage:")
    fmt.Println("  helm-doc [options]")
    fmt.Println()
    fmt.Println("helm-doc reads the values file of a chart in the root of your repository and creates a table of the values")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println(fmt.Sprintf("  --debug         %s", DEBUG_HELP))
    fmt.Println(fmt.Sprintf("  --dry-run       %s", DRY_RUN_HELP))
    fmt.Println(fmt.Sprintf("  --help          %s", HELP_HELP))
}

func parseCommandLine() HelmDocArgs {
    var args HelmDocArgs

    flag.BoolVar(&args.debug, "debug", false, DEBUG_HELP)
    flag.BoolVar(&args.dryRun, "dry-run", false, DRY_RUN_HELP)
    flag.BoolVar(&args.help, "help", false, HELP_HELP)
    flag.Parse()

    if args.help {
        printHelp()
        os.Exit(0)
    }

    return args
}

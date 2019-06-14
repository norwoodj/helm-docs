package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

type HelmDocArgs struct {
	help   bool
	debug  bool
	dryRun bool
}

const debugHelp = "Print debug output"
const dryRunHelp = "Don't actually render any markdown files with docs, juts print to stdout"
const helpHelp = "Print this help menu, then exist"

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  helm-doc [options]")
	fmt.Println()
	fmt.Println("helm-doc reads the values file of a chart in the root of your repository and creates a table of the values")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println(fmt.Sprintf("  --debug         %s", debugHelp))
	fmt.Println(fmt.Sprintf("  --dry-run       %s", dryRunHelp))
	fmt.Println(fmt.Sprintf("  --help          %s", helpHelp))
}

func parseCommandLine() HelmDocArgs {
	var args HelmDocArgs

	flag.BoolVar(&args.debug, "debug", false, debugHelp)
	flag.BoolVar(&args.dryRun, "dry-run", false, dryRunHelp)
	flag.BoolVar(&args.help, "help", false, helpHelp)
	flag.Parse()

	if args.help {
		printHelp()
		os.Exit(0)
	}

	return args
}

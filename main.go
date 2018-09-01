package main


func main() {
    args := parseCommandLine()
    printDocumentation(args.debug, args.dryRun)
}

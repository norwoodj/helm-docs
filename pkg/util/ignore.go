package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/ignore"
)

var defaultIgnore = map[string]bool{
	".git": true,
}

type IgnoreContext struct {
	rules       *ignore.Rules
	relativeDir string
}

func parseIgnoreFilePathToRules(filename string) (*ignore.Rules, error) {
	ignoreFile, err := os.Open(filename)

	if os.IsNotExist(err) {
		log.Debugf("No ignore file found at %s, using empty ignore rules", filename)
		return ignore.Empty(), nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open ignore file at %s: %s", filename, err)
	}

	ignoreRules, err := ignore.Parse(bufio.NewReader(ignoreFile))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ignore rules from file %s: %s", filename, err)
	}

	log.Debugf("Found ignore file at %s, using those ignore rules", filename)
	return ignoreRules, nil
}

func NewIgnoreContext(ignoreFilename string) IgnoreContext {
	gitRepositoryRoot, err := FindGitRepositoryRoot()

	// If we got an error reading the repository root, then let's try for a ignore file in this directory
	if err != nil {
		ignoreRules, err := parseIgnoreFilePathToRules(ignoreFilename)

		if err != nil {
			log.Warnf("Using empty ignore rules due to error: %s", err)
			return IgnoreContext{rules: ignore.Empty()}
		}

		absoluteWorkingDir, _ := filepath.Abs(".")
		return IgnoreContext{rules: ignoreRules, relativeDir: absoluteWorkingDir}
	}

	// Otherwise, let's look for a ignore file at the repository root and parse it, storing that files are ignored relative
	// to that directory
	ignoreRules, err := parseIgnoreFilePathToRules(filepath.Join(gitRepositoryRoot, ignoreFilename))

	if err != nil {
		log.Warnf("Using empty ignore rules due to error: %s", err)
		return IgnoreContext{rules: ignore.Empty()}
	}

	return IgnoreContext{rules: ignoreRules, relativeDir: gitRepositoryRoot}
}

func (i IgnoreContext) ShouldIgnore(path string, fi os.FileInfo) bool {
	pathRelativeToIgnoreFile, err := filepath.Rel(i.relativeDir, path)

	if err != nil {
		return false
	}

	return i.rules.Ignore(pathRelativeToIgnoreFile, fi) || defaultIgnore[pathRelativeToIgnoreFile]
}

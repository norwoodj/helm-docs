package util

import (
	"os/exec"
	"strings"
)

func FindGitRepositoryRoot() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(path)), nil
}

package util

import "path"

func IsRelativePath(filePath string) bool {
	return (filePath[0] == '.') && path.Base(filePath) != filePath
}

func IsBaseFilename(filePath string) bool {
	return path.Base(filePath) == filePath
}

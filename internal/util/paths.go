package util

import (
	"os"
	"path/filepath"
)

// WorkingDirectoryPath returns the absolute path to the provided path
func WorkingDirectoryPath(p string) string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, p)
}

// GetAbsPath returns the absolute path to the provided path
func AbsolutePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return abs
}

package util

import (
    "path/filepath"
)

// Get absolute to provided path
func GetAbsPath(path string) (string) {
    abs, err := filepath.Abs(path)
    if err != nil {
        panic(err)
    }

    return abs
}

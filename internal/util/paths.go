package util

import (
    "path/filepath"
)

func GetAbsPath(path string) (string) {
    abs, err := filepath.Abs(path)
    if err != nil {
        panic(err)
    }

    return abs
}

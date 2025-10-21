package main

import (
	"errors"
	"io/fs"
	"os"
)

// returns true if file exists
func exists(f string) bool {
	if _, err := os.Stat(f); errors.Is(err, fs.ErrNotExist) {
		Warnf("File not found %s", f)
		return false
	}
	return true
}

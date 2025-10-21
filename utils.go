package main

import (
	"errors"
	"io/fs"
	"os"
)

// returns true if file exists
func exists(f string) bool {
	Warnf("checking path %s", f)
	if _, err := os.Stat(f); errors.Is(err, fs.ErrNotExist) {
		Errorf("File not found %s", f)
		return false
	}
	// Debug("found", f)
	return true
}

package main

import (
	"os"
	"path/filepath"
)

func init() {
	systemHostsPath = filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
}

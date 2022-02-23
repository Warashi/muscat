package cmd

import (
	"log"
	"os"
	"path/filepath"
)

func mustGetSocketPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("os.UserHomeDir: %v", err)
	}
	dir := filepath.Join(home, ".muscat")
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("os.MkdirAll: %v", err)
	}
	return filepath.Join(dir, "socket")
}

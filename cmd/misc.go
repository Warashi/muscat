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
	return filepath.Join(home, ".muscat.socket")
}

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Warashi/muscat/v2/config"
)

func getConfigPath() (string, error) {
	if configFilePath := os.Getenv("MUSCAT_CONFIG"); configFilePath != "" {
		return configFilePath, nil
	}
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "muscat/config.json"), nil
	}
	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, ".config/muscat/config.json"), nil
	}
	return "", fmt.Errorf("failed to determine config path")
}

func getConfig(ctx context.Context) (*config.Config, error) {
	p, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	return config.LoadFromPath(ctx, p)
}

func parseConfigAddr(addr string) (string, string) {
	if strings.HasPrefix(":", addr) {
		return "unix", addr
	}
	if strings.HasPrefix(addr, "localhost:") {
		return "unix", addr
	}
	return "tcp", addr
}

func mustGetListenArgs(ctx context.Context) (string, string) {
	cfg, err := getConfig(ctx)
	if err == nil && cfg.Addr != "" {
		return parseConfigAddr(cfg.Addr)
	}
	return "unix", mustGetSocketPath()
}

func mustGetSocketPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("os.UserHomeDir: %v", err)
	}
	return filepath.Join(home, ".muscat.socket")
}

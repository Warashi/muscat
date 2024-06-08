package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Addr string
}

func LoadFromPath(ctx context.Context, p string) (*Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &c, nil
}

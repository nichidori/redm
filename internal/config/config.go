package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	URL    string `json:"url"`
	APIKey string `json:"api_key"`
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory: %w", err)
	}
	return filepath.Join(home, ".redm", "config.json"), nil
}

func Load() (*Config, error) {
	p, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("not logged in: %w", err)
		}
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	if cfg.URL == "" || cfg.APIKey == "" {
		return nil, errors.New("config file is incomplete")
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	p, err := ConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	return nil
}

func Delete() error {
	p, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("remove config: %w", err)
	}
	return nil
}

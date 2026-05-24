package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Server    string `json:"server"`
	Token     string `json:"token"`
	Email     string `json:"email"`
	Workspace string `json:"workspace"`
}

func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Local", "linkstate", "config.json")
	}
	return filepath.Join(home, ".linkstate", "config.json")
}

func Load() (*Config, error) {
	cfg := &Config{Server: "http://localhost:8080"}
	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			cfg.Workspace = defaultWorkspace()
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.Server == "" {
		cfg.Server = "http://localhost:8080"
	}
	if cfg.Workspace == "" {
		cfg.Workspace = defaultWorkspace()
	}
	return cfg, nil
}

func Save(cfg *Config) error {
	dir := filepath.Dir(Path())
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0600)
}

func defaultWorkspace() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, "linkstate", "workspace")
}

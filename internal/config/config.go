package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Server           string `json:"server"`
	Token            string `json:"token"`
	Email            string `json:"email"`
	Path             string `json:"path"`
	LastTrackSlug    string `json:"last_track_slug,omitempty"`
	LastLessonSlug   string `json:"last_lesson_slug,omitempty"`
	LastLessonTitle  string `json:"last_lesson_title,omitempty"`
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
	cfg := &Config{Server: "http://localhost", Path: defaultPath()}
	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) { return cfg, nil }
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil { return nil, err }
	if cfg.Server == "" { cfg.Server = "http://localhost" }
	if cfg.Path == "" { cfg.Path = defaultPath() }
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

func defaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil { return "." }
	return filepath.Join(home, "linkstate")
}
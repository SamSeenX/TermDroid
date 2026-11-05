package bridge

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SavedPIN string `json:"saved_pin,omitempty"`
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "termdroid", "config.json")
}

func LoadConfig() Config {
	path := getConfigPath()
	if path == "" {
		return Config{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}
	}
	return cfg
}

func SaveConfig(cfg Config) error {
	path := getConfigPath()
	if path == "" {
		return os.ErrNotExist
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

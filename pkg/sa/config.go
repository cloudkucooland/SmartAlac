package sa

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func DefaultConfigPath(filename string) string {
	base, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "SmartAlac", filename)
}

func DefaultDataDir() string {
	base, err := os.UserCacheDir() // We'll use Cache for temp work files in BME, or DataDir?
	if err != nil {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "SmartAlac")
}

func DefaultCachePath() string {
	base, err := os.UserCacheDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "SmartAlac", "smartalac.db")
}

type SharedConfig struct {
	FinalDir     string `json:"final_dir"`
	AcoustIDKey  string `json:"acoustid_key"`
	DiscogsToken string `json:"discogs_token"`
}

func LoadSharedConfig() *SharedConfig {
	path := DefaultConfigPath("config.json")
	c := &SharedConfig{
		FinalDir: "/home/music/alac",
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return c
	}

	json.Unmarshal(data, c)
	return c
}

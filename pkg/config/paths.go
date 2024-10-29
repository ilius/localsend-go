package config

import (
	"log/slog"
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	if os.Getenv("CONFIG_FILE") != "" {
		return filepath.Dir(Path())
	}
	return platformConfigDir()
}

func Path() string {
	_path := os.Getenv("CONFIG_FILE")
	if _path != "" {
		absPath, err := filepath.Abs(_path)
		if err == nil {
			return absPath
		} else {
			slog.Error("bad CONFIG_FILE", "CONFIG_FILE", _path, "err", err)
		}
	}
	return filepath.Join(GetConfigDir(), fileName)
}

//go:build windows
// +build windows

package config

import (
	"os"
	"path/filepath"
)

func platformConfigDir() string {
	// HOMEDRIVE := os.Getenv("HOMEDRIVE")
	// HOMEPATH := os.Getenv("HOMEPATH")
	// homeDir := filepath.Join(HOMEDRIVE, HOMEPATH)
	// user := os.Getenv("USERNAME")
	// tmpDir := os.Getenv("TEMP")
	appData := os.Getenv("APPDATA")
	confDir := filepath.Join(appData, "localsend-go")
	return confDir
}

/*
func GetCacheDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		// Windows Vista or older
		appData := os.Getenv("APPDATA")
		var err error
		localAppData, err = filepath.Abs(filepath.Join(appData, "..", "Local"))
		if err != nil {
			slog.Error("error", "err", err)
			return ""
		}
	}
	return filepath.Join(localAppData, "localsend-go", "Cache")
}
*/

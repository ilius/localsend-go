//go:build darwin
// +build darwin

package config

import (
	"os"
	"path/filepath"
)

func platformConfigDir() string {
	return filepath.Join(
		os.Getenv(S_HOME),
		"Library/Preferences/localsend-go",
	)
}

//func GetCacheDir() string {
//	return filepath.Join(os.Getenv(S_HOME), "Library", "Caches", "localsend-go")
//}

//go:build !windows && !darwin
// +build !windows,!darwin

package config

import (
	"os"
	"path/filepath"
)

func platformConfigDir() string {
	parent := os.Getenv("XDG_CONFIG_HOME")
	if parent == "" {
		parent = filepath.Join(os.Getenv(S_HOME), ".config")
	}
	return filepath.Join(parent, "localsend-go")
}

//func GetCacheDir() string {
//	return filepath.Join(os.Getenv(S_HOME), ".cache", "localsend-go")
//}

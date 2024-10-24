package config

import (
	"embed"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const fileName = "config.toml"

//go:embed config.toml
var EmbeddedConfig embed.FS

type Config struct {
	NameOfDevice string `toml:"name"`
	Receive      struct {
		SaveUserID  int `toml:"saveUserID"`
		SaveGroupID int `toml:"saveGroupID"`
	} `toml:"receive"`
	Functions struct {
		HttpFileServer  bool `toml:"http_file_server"`
		LocalSendServer bool `toml:"local_send_server"`
	} `toml:"functions"`
}

var ConfigData Config

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

func GetConfigDir() string {
	if os.Getenv("CONFIG_FILE") != "" {
		return filepath.Dir(Path())
	}
	return platformConfigDir()
}

func Init() {
	var bytes []byte
	var err error

	configPath := filepath.Join(GetConfigDir(), fileName)

	slog.Info("Reading user config file", "configPath", configPath)

	// Try to read configuration files from external file system
	bytes, err = os.ReadFile(configPath)
	if err != nil {
		slog.Info("Failed to read external configuration file, using built-in configuration")
		// If reading the external file fails, read from the embedded file system
		bytes, err = EmbeddedConfig.ReadFile("config.toml")
		if err != nil {
			log.Fatalf("Error reading embedded config file: %v", err)
		}
	}

	err = toml.Unmarshal(bytes, &ConfigData)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	slog.Info("Loaded user config file", "configData", ConfigData)
}

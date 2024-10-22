package config

import (
	"embed"
	"log"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

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

func init() {
	var bytes []byte
	var err error

	// Try to read configuration files from external file system
	bytes, err = os.ReadFile("internal/config/config.toml")
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
}

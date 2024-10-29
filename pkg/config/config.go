package config

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilius/localsend-go/pkg/alias"
	"github.com/ilius/localsend-go/pkg/go-clipboard"
	"github.com/ilius/localsend-go/pkg/toml"
)

const fileName = "config.toml"

//go:embed config.toml
var embedFS embed.FS

type Config struct {
	NameOfDevice string `toml:"name"`
	NameLanguage string `toml:"name_language"`
	Receive      struct {
		Directory          string `toml:"directory"`
		MaxFileSize        int    `toml:"max_file_size"`
		SaveUserID         int    `toml:"save_user_id"`
		SaveGroupID        int    `toml:"save_group_id"`
		Clipboard          bool   `toml:"clipboard"`
		ExitAfterFileCount int    `toml:"exit_after_file_count"`
	} `toml:"receive"`
	Send struct {
		Directory string `toml:"directory"`
	} `toml:"send"`
	Functions struct {
		HttpFileServer  bool `toml:"http_file_server"`
		LocalSendServer bool `toml:"local_send_server"`
	} `toml:"functions"`
	Logging struct {
		NoColor bool   `toml:"no_color"`
		Level   string `toml:"level"`
	} `toml:"logging"`
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

func GetConfigDir() string {
	if os.Getenv("CONFIG_FILE") != "" {
		return filepath.Dir(Path())
	}
	return platformConfigDir()
}

func Init() *Config {
	conf := &Config{}
	var bytes []byte
	var err error

	configPath := filepath.Join(GetConfigDir(), fileName)

	// read from the embedded file system
	{
		bytes, err = embedFS.ReadFile("config.toml")
		if err != nil {
			panic(fmt.Sprintf("Error reading embedded config file: %v", err))
		}
		err = toml.Unmarshal(bytes, conf)
		if err != nil {
			panic(fmt.Sprintf("Error parsing default config file: %v", err))
		}
	}

	slog.Info("Trying to read user config file", "configPath", configPath)
	bytes, err = os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("Failed to read external configuration file, using built-in configuration")
		}
	} else {
		err = toml.Unmarshal(bytes, conf)
		if err != nil {
			panic(fmt.Sprintf("Error parsing config file: %v", err))
		}
		slog.Info("Loaded user config file", "configData", conf)
	}
	if conf.NameOfDevice == "" {
		name := alias.GenerateRandomAlias(conf.NameLanguage)
		slog.Info("Using random name/alias: ", "name", name)
		conf.NameOfDevice = name
	}
	if conf.Receive.Clipboard {
		clipboard.Init()
	}
	return conf
}

package config

import (
	"embed"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var EmbeddedConfig embed.FS

type Config struct {
	NameOfDevice string `yaml:"name"`
	Functions    struct {
		HttpFileServer  bool `yaml:"http_file_server"`
		LocalSendServer bool `yaml:"local_send_server"`
	} `yaml:"functions"`
}

var ConfigData Config

func init() {
	var bytes []byte
	var err error

	// Try to read configuration files from external file system
	bytes, err = os.ReadFile("internal/config/config.yaml")
	if err != nil {
		fmt.Println("读取外部配置文件失败，使用内置配置")
		// If reading the external file fails, read from the embedded file system
		bytes, err = EmbeddedConfig.ReadFile("config.yaml")
		if err != nil {
			log.Fatalf("Error reading embedded config file: %v", err)
		}
	}

	err = yaml.Unmarshal(bytes, &ConfigData)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}

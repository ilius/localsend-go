package shared

import (
	"sync"

	"codeberg.org/ilius/localsend-go/pkg/config"
	"codeberg.org/ilius/localsend-go/pkg/models"
	"codeberg.org/ilius/localsend-go/pkg/utils"
)

// Global device record hash table and mutex, Message information

var (
	discoveredDevices      = make(map[string]*models.BroadcastMessage)
	discoveredDevicesMutex sync.Mutex
)

func AddDiscoveredDevice(ip string, msg *models.BroadcastMessage) bool {
	discoveredDevicesMutex.Lock()
	defer discoveredDevicesMutex.Unlock()
	if _, exists := discoveredDevices[ip]; !exists {
		discoveredDevices[ip] = msg
		return true
	}
	return false
}

// https://github.com/localsend/protocol?tab=readme-ov-file#71-device-type
func GetMesssage(conf *config.Config) models.BroadcastMessage {
	return models.BroadcastMessage{
		Alias:       conf.NameOfDevice,
		Version:     "2.0",
		DeviceModel: utils.OSType(),
		DeviceType:  "headless", // Indicates that it is running without GUI
		Fingerprint: "random-string",
		Port:        53317,
		Protocol:    "http",
		Download:    true,
		Announce:    true,
	}
}

package shared

import (
	"sync"

	"github.com/ilius/localsend-go/pkg/config"
	. "github.com/ilius/localsend-go/pkg/models"
	"github.com/ilius/localsend-go/pkg/utils"
)

// Global device record hash table and mutex, Message information

var (
	discoveredDevices      = make(map[string]*BroadcastMessage)
	discoveredDevicesMutex sync.Mutex
)

func AddDiscoveredDevice(ip string, msg *BroadcastMessage) bool {
	discoveredDevicesMutex.Lock()
	defer discoveredDevicesMutex.Unlock()
	if _, exists := discoveredDevices[ip]; !exists {
		discoveredDevices[ip] = msg
		return true
	}
	return false
}

// https://github.com/localsend/protocol?tab=readme-ov-file#71-device-type
func GetMesssage() BroadcastMessage {
	return BroadcastMessage{
		Alias:       config.ConfigData.NameOfDevice,
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

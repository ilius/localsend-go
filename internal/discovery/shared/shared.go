package shared

import (
	"sync"

	"localsend_cli/internal/config"
	. "localsend_cli/internal/models"
	"localsend_cli/internal/utils"
)

// Global device record hash table and mutex, Message information

var (
	DiscoveredDevices = make(map[string]BroadcastMessage)
	Mu                sync.Mutex
)

// https://github.com/localsend/protocol?tab=readme-ov-file#71-device-type
var Messsage BroadcastMessage = BroadcastMessage{
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

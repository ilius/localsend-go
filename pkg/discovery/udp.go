package discovery

import (
	"encoding/json"
	"net"
	"time"

	"github.com/ilius/localsend-go/pkg/discovery/shared"
	"github.com/ilius/localsend-go/pkg/models"
)

// startBroadcast sends a broadcast message
func (d *discoveryImp) startBroadcast() {
	d.log.Info("Starting broadcast", "multicast_ip", d.conf.MulticastIP)

	// Set the multicast address and port
	multicastAddr := &net.UDPAddr{
		IP:   net.ParseIP(d.conf.MulticastIP),
		Port: 53317,
	}

	msg := shared.GetMesssage(d.conf)
	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	// Create a local address and bind it to all interfaces
	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		d.log.Error("Error creating UDP connection:", "err", err)
		return
	}
	defer conn.Close()
	for {
		num_bytes, err := conn.WriteToUDP(data, multicastAddr)
		if err != nil {
			d.log.Error("Failed to send message:", "err", err)
			panic(err)
		}
		d.log.Debug("Writen to multicastAddr", "num_bytes", num_bytes)
		// log
		d.log.Debug("UDP Broadcast message sent!")
		time.Sleep(5 * time.Second) // Send a broadcast message every 5 seconds
	}
}

// listenForBroadcasts listens for UDP broadcast messages
func (d *discoveryImp) listenForBroadcasts() {
	d.log.Info("Listening for broadcasts", "multicast_ip", d.conf.MulticastIP)

	// Set the multicast address and port
	multicastAddr := &net.UDPAddr{
		IP:   net.ParseIP(d.conf.MulticastIP),
		Port: 53317,
	}

	// Create a UDP multicast listening connection
	conn, err := net.ListenMulticastUDP("udp", nil, multicastAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Set the read buffer size
	conn.SetReadBuffer(1024)

	for {
		buf := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		var message models.BroadcastMessage
		err = json.Unmarshal(buf[:n], &message)
		if err != nil {
			d.log.Error("Failed to unmarshal broadcast message:", "err", err)
			continue
		}

		if shared.AddDiscoveredDevice(remoteAddr.IP.String(), &message) {
			d.log.Info("Discovered device", "alias", message.Alias, "deviceModel", message.DeviceModel, "ip", remoteAddr.IP.String())
		}
	}
}

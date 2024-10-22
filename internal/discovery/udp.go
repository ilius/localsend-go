package discovery

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"localsend_cli/internal/discovery/shared"
	. "localsend_cli/internal/models"
)

// StartBroadcast sends a broadcast message
func StartBroadcast() {
	// Set the multicast address and port
	multicastAddr := &net.UDPAddr{
		IP:   net.ParseIP("224.0.0.167"),
		Port: 53317,
	}

	data, err := json.Marshal(shared.Messsage)
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
		fmt.Println("Error creating UDP connection:", err)
		return
	}
	defer conn.Close()
	for {
		_, err := conn.WriteToUDP(data, multicastAddr)
		if err != nil {
			fmt.Println("Failed to send message:", err)
			panic(err)
		}
		// fmt.Println(num, "bytes write to multicastAddr")
		// log
		// fmt.Println("UDP Broadcast message sent!")
		time.Sleep(5 * time.Second) // Send a broadcast message every 5 seconds
	}
}

// ListenForBroadcasts listens for UDP broadcast messages
func ListenForBroadcasts() {
	fmt.Println("Listening for broadcasts...")

	// Set the multicast address and port
	multicastAddr := &net.UDPAddr{
		IP:   net.ParseIP("224.0.0.167"),
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
		var message BroadcastMessage
		err = json.Unmarshal(buf[:n], &message)
		if err != nil {
			fmt.Println("Failed to unmarshal broadcast message:", err)
			continue
		}

		shared.Mu.Lock()
		if _, exists := shared.DiscoveredDevices[remoteAddr.IP.String()]; !exists {
			shared.DiscoveredDevices[remoteAddr.IP.String()] = message
			fmt.Printf("Discovered device: %s (%s) at %s\n", message.Alias, message.DeviceModel, remoteAddr.IP.String())
		}
		shared.Mu.Unlock()
	}
}

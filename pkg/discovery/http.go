package discovery

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"codeberg.org/ilius/localsend-go/pkg/discovery/shared"
	"codeberg.org/ilius/localsend-go/pkg/models"

	probing "github.com/prometheus-community/pro-bing"
)

// getLocalIP Get the local IP address
func getLocalIP() ([]net.IP, error) {
	ips := make([]net.IP, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			v, ok := addr.(*net.IPNet)
			if ok {
				if v.IP.To4() != nil && !v.IP.IsLoopback() {
					ips = append(ips, v.IP)
				}
			}
		}
	}
	return ips, nil
}

// pingScan uses ICMP ping to scan all active devices on the LAN
func (d *discoveryImp) pingScan() ([]string, error) {
	var ips []string
	ipGroup, err := getLocalIP()
	d.log.Debug("pingScan", "ip", ips)
	if err != nil {
		return nil, err
	}
	for _, i := range ipGroup {
		ip := i.Mask(net.IPv4Mask(255, 255, 255, 0)) // Assume the subnet mask is 24
		ip4 := ip.To4()
		if ip4 == nil {
			return nil, fmt.Errorf("invalid IPv4 address")
		}

		var wg sync.WaitGroup
		var mu sync.Mutex

		for i := 1; i < 255; i++ {
			ip4[3] = byte(i)
			targetIP := ip4.String()

			wg.Add(1)
			go func(ip string) {
				defer wg.Done()
				pinger, err := probing.NewPinger(ip)
				if err != nil {
					d.log.Error("Failed to create pinger:", "err", err)
					return
				}
				pinger.SetPrivileged(true)
				pinger.Count = 1
				pinger.Timeout = time.Second * 1

				pinger.OnRecv = func(pkt *probing.Packet) {
					mu.Lock()
					ips = append(ips, ip)
					mu.Unlock()
				}
				err = pinger.Run()
				if err != nil {
					// Ignore ping failures
					d.log.Debug("Failed to run pinger:", "err", err)
					return
				}
			}(targetIP)
		}

		wg.Wait()
	}
	d.log.Debug("pingScan (end)", "ip", ips)
	return ips, nil
}

// startHTTPBroadcast sends HTTP requests to all IPs in the LAN
func (d *discoveryImp) startHTTPBroadcast() {
	msg := shared.GetMesssage(d.conf)
	for {
		data, err := json.Marshal(msg)
		d.log.Debug(string(data))
		if err != nil {
			panic(err)
		}

		ips, err := d.pingScan()
		if err != nil {
			d.log.Error("Failed to discover devices via ping scan:", "err", err)
			return
		}

		var wg sync.WaitGroup
		for _, ip := range ips {
			wg.Add(1)
			go func(ip string) {
				defer wg.Done()
				ctx := context.Background()
				d.sendHTTPRequest(ctx, ip, data)
			}(ip)
		}

		wg.Wait()
		d.log.Debug("HTTP broadcast messages sent!")
		time.Sleep(5 * time.Second) // Send HTTP broadcast message every 5 seconds
	}
}

// sendHTTPRequest sends HTTP requests
func (d *discoveryImp) sendHTTPRequest(ctx context.Context, ip string, data []byte) {
	url := fmt.Sprintf("https://%s:53317/api/localsend/v2/register", ip)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		d.log.Info("Failed to create HTTP request:", "err", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		d.log.Error("Failed to read HTTP response body:", "err", err)
		return
	}
	var response models.BroadcastMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		d.log.Error("Failed to parse HTTP response", "ip", ip, "err", err)
		return
	}
	if shared.AddDiscoveredDevice(ip, &response) {
		d.log.Info("Discovered device", "alias", response.Alias, "deviceModel", response.DeviceModel, "ip", ip)
	}
	// slog.Info("Discovered device", "alias", response.Alias, "deviceModel", response.DeviceModel, "ip", ip)
}

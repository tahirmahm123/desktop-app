//go:build nowifi
// +build nowifi
package wifiNotifier

import "github.com/tahirmahm123/vpn-desktop-app/daemon/logger"

// GetAvailableSSIDs returns the list of the names of available Wi-Fi networks
func GetAvailableSSIDs() []string {
	return nil
}

// GetCurrentSSID returns current WiFi SSID
func GetCurrentSSID() string {
	return ""
}

// GetCurrentNetworkIsInsecure returns current security mode
func GetCurrentNetworkIsInsecure() bool {
	return false
}

// SetWifiNotifier initializes a handler method 'OnWifiChanged'
func SetWifiNotifier(cb func(string)) error {
	logger.Debug("WiFi functionality disabled in this build")
	return nil
}

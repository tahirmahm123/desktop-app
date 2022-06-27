
package service

import (
	"time"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/wifiNotifier"
)

type wifiInfo struct {
	ssid       string
	isInsecure bool
}

var lastWiFiInfo *wifiInfo
var timerDelayedNotify *time.Timer

const delayBeforeWiFiChangeNotify = time.Second * 1

func (s *Service) initWiFiFunctionality() (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("initWiFiFunctionality PANIC (recovered): ", r)
		}
	}()

	return wifiNotifier.SetWifiNotifier(s.onWiFiChanged)
}

func (s *Service) onWiFiChanged(ssid string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("onWiFiChanged PANIC (recovered): ", r)
		}
	}()

	isInsecure := wifiNotifier.GetCurrentNetworkIsInsecure()

	lastWiFiInfo = &wifiInfo{
		ssid,
		isInsecure}

	// do delay before processing wifi change
	// (same wifi change event can occur several times in short period of time)
	if timerDelayedNotify != nil {
		timerDelayedNotify.Stop()
		timerDelayedNotify = nil
	}
	timerDelayedNotify = time.AfterFunc(delayBeforeWiFiChangeNotify, func() {
		if lastWiFiInfo == nil || lastWiFiInfo.ssid != ssid || lastWiFiInfo.isInsecure != isInsecure {
			return // do nothing (new wifi info available)
		}

		// notify clients about WiFi change
		s._evtReceiver.OnWiFiChanged(ssid, isInsecure)
	})
}

// GetWiFiCurrentState returns info about currently connected wifi
func (s *Service) GetWiFiCurrentState() (ssid string, isInsecureNetwork bool) {
	return wifiNotifier.GetCurrentSSID(), wifiNotifier.GetCurrentNetworkIsInsecure()
}

// GetWiFiAvailableNetworks returns list of available WIFI networks
func (s *Service) GetWiFiAvailableNetworks() []string {
	return wifiNotifier.GetAvailableSSIDs()
}

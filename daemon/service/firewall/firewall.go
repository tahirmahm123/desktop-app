//
//  Daemon for VPN Client Desktop
//  https://github.com/tahirmahm123/vpn-desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the Daemon for VPN Desktop.
//
//  The Daemon for VPN Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The Daemon for VPN Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the Daemon for VPN Desktop. If not, see <https://www.gnu.org/licenses/>.
//

package firewall

import (
	"fmt"
	"net"
	"sync"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("frwl")
}

var (
	connectedClientInterfaceIP   net.IP
	connectedClientInterfaceIPv6 net.IP
	connectedClientPort          int
	connectedHostIP              net.IP
	connectedHostPort            int
	connectedIsTCP               bool
	mutex                        sync.Mutex
	isClientPaused               bool
)

// Initialize is doing initialization stuff
// Must be called on application start
func Initialize() error {
	return implInitialize()
}

// SetEnabled - change firewall state
func SetEnabled(enable bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	if enable {
		log.Info("Enabling...")
	} else {
		log.Info("Disabling...")
	}

	err := implSetEnabled(enable)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to change firewall state : %w", err)
	}

	if enable {
		// To fulfill such flow (example): FWEnable -> Connected -> FWDisable -> FWEnable
		// Here we should notify that client is still connected
		// We must not do it in Paused state!
		clientAddr := connectedClientInterfaceIP
		clientAddrIPv6 := connectedClientInterfaceIPv6
		if clientAddr != nil && !isClientPaused {
			e := implClientConnected(clientAddr, clientAddrIPv6, connectedClientPort, connectedHostIP, connectedHostPort, connectedIsTCP)
			if e != nil {
				log.Error(e)
			}
		}
	}
	return err
}

// SetPersistant - set persistant firewall state and enable it if necessary
func SetPersistant(persistant bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Info(fmt.Sprintf("Persistent:%t", persistant))

	err := implSetPersistant(persistant)
	if err != nil {
		log.Error(err)
	}
	return err
}

// GetEnabled - get firewall state
func GetEnabled() (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()
	log.Info("Getting status...")

	ret, err := implGetEnabled()
	if err != nil {
		log.Error(err)
	} else {
		log.Info("    ", ret)
	}

	return ret, err
}

// ClientPaused saves info about paused state of vpn
func ClientPaused() {
	isClientPaused = true
}

// ClientResumed saves info about resumed state of vpn
func ClientResumed() {
	isClientPaused = false
}

// ClientConnected - allow communication for local vpn/client IP address
func ClientConnected(clientLocalIPAddress net.IP, clientLocalIPv6Address net.IP, clientPort int, serverIP net.IP, serverPort int, isTCP bool) error {
	mutex.Lock()
	defer mutex.Unlock()
	ClientResumed()

	log.Info("Client connected: ", clientLocalIPAddress)

	connectedClientInterfaceIP = clientLocalIPAddress
	connectedClientInterfaceIPv6 = clientLocalIPv6Address
	connectedClientPort = clientPort
	connectedHostIP = serverIP
	connectedHostPort = serverPort
	connectedIsTCP = isTCP

	err := implClientConnected(clientLocalIPAddress, clientLocalIPv6Address, clientPort, serverIP, serverPort, isTCP)
	if err != nil {
		log.Error(err)
	}
	return err
}

// ClientDisconnected - Remove all hosts exceptions
func ClientDisconnected() error {
	mutex.Lock()
	defer mutex.Unlock()
	ClientResumed()

	// Remove client interface from exceptions
	if connectedClientInterfaceIP != nil {
		connectedClientInterfaceIP = nil
		connectedClientInterfaceIPv6 = nil
		log.Info("Client disconnected")
		err := implClientDisconnected()
		if err != nil {
			log.Error(err)
		}
		return err
	}
	return nil
}

// AddHostsToExceptions - allow comminication with this hosts
// Note: if isPersistent == false -> all added hosts will be removed from exceptions after client disconnection (after call 'ClientDisconnected()')
// Arguments:
//	* IPs			-	list of IP addresses to ba allowed
//	* onlyForICMP	-	(applicable only for Linux) try add rule to allow only ICMP protocol for this IP
//	* isPersistent	-	keep rule enabled even if VPN disconnected
func AddHostsToExceptions(IPs []net.IP, onlyForICMP bool, isPersistent bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	err := implAddHostsToExceptions(IPs, onlyForICMP, isPersistent)
	if err != nil {
		log.Error("Failed to add hosts to exceptions:", err)
	}

	return err
}

func RemoveHostsFromExceptions(IPs []net.IP, onlyForICMP bool, isPersistent bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	err := implRemoveHostsFromExceptions(IPs, onlyForICMP, isPersistent)
	if err != nil {
		log.Error("Failed to remove hosts from exceptions:", err)
	}

	return err
}

// AllowLAN - allow/forbid LAN communication
func AllowLAN(allowLan bool, allowLanMulticast bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Info(fmt.Sprintf("allowLan:%t allowMulticast:%t", allowLan, allowLanMulticast))

	err := implAllowLAN(allowLan, allowLanMulticast)
	if err != nil {
		log.Error(err)
	}
	return err
}

// SetManualDNS - configure firewall to allow DNS which is out of VPN tunnel
// Applicable to Windows implementation (to allow custom DNS from local network)
func SetManualDNS(addr net.IP) error {
	mutex.Lock()
	defer mutex.Unlock()

	err := implSetManualDNS(addr)
	if err != nil {
		log.Error(err)
	}
	return err
}

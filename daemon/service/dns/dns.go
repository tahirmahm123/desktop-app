//
//  Daemon for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the Daemon for IVPN Client Desktop.
//
//  The Daemon for IVPN Client Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The Daemon for IVPN Client Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the Daemon for IVPN Client Desktop. If not, see <https://www.gnu.org/licenses/>.
//

package dns

import (
	"fmt"
	"net"
	"strings"

	"github.com/ivpn/desktop-app/daemon/logger"
)

var log *logger.Logger
var lastManualDNS DnsSettings

func init() {
	log = logger.NewLogger("dns")
}

type DnsEncryption int

const (
	EncryptionNone         DnsEncryption = 0
	EncryptionDnsOverTls   DnsEncryption = 1
	EncryptionDnsOverHttps DnsEncryption = 2
)

type DnsSettings struct {
	DnsHost     string // DNS host IP address
	Encryption  DnsEncryption
	DohTemplate string // DoH/DoT template URI (for Encryption = DnsOverHttps or Encryption = DnsOverTls)
}

func (d DnsSettings) Equal(x DnsSettings) bool {
	if d.Encryption != x.Encryption ||
		d.DohTemplate != x.DohTemplate ||
		d.DnsHost != x.DnsHost {
		return false
	}
	return true
}

func (d DnsSettings) IsIPv6() (bool, error) {
	ip := d.Ip()
	if ip == nil {
		return false, fmt.Errorf("unable to determine IP protocol version for the DnsSettings object (object is not initialized)")
	}
	return ip.To4() == nil, nil
}

func (d DnsSettings) Ip() net.IP {
	return net.ParseIP(d.DnsHost)
}

func (d DnsSettings) IsEmpty() bool {
	if strings.TrimSpace(d.DnsHost) == "" {
		return true
	}
	ip := d.Ip()
	if ip == nil || ip.Equal(net.IPv4zero) || ip.Equal(net.IPv4bcast) || ip.Equal(net.IPv6zero) {
		return true
	}
	return false
}

func (d DnsSettings) InfoString() string {
	if d.IsEmpty() {
		return "<none>"
	}
	host := strings.TrimSpace(d.DnsHost)
	template := strings.TrimSpace(d.DohTemplate)

	switch d.Encryption {
	case EncryptionDnsOverTls:
		return host + " (DoT " + template + ")"
	case EncryptionDnsOverHttps:
		return host + " (DoH " + template + ")"
	case EncryptionNone:
		return host
	default:
		return host + " (UNKNOWN ENCRYPTION)"
	}
}

// Initialize is doing initialization stuff
// Must be called on application start
func Initialize() error {
	return implInitialize()
}

// Pause pauses DNS (restore original DNS)
func Pause() error {
	return implPause()
}

// Resume resuming DNS (set DNS back which was before Pause)
func Resume(defaultDNS DnsSettings) error {
	return implResume(defaultDNS)
}

func EncryptionAbilities() (dnsOverHttps, dnsOverTls bool, err error) {
	return implGetDnsEncryptionAbilities()
}

// SetDefault set DNS configuration treated as default (non-manual) configuration
// 'dnsCfg' parameter - DNS configuration
// 'localInterfaceIP' (obligatory only for Windows implementation) - local IP of VPN interface
func SetDefault(dnsCfg DnsSettings, localInterfaceIP net.IP) error {
	ret := implSetManual(dnsCfg, localInterfaceIP)
	if ret == nil {
		lastManualDNS = DnsSettings{}
	}
	return ret
}

// SetManual - set manual DNS.
// 'dnsCfg' parameter - DNS configuration
// 'localInterfaceIP' (obligatory only for Windows implementation) - local IP of VPN interface
func SetManual(dnsCfg DnsSettings, localInterfaceIP net.IP) error {
	ret := implSetManual(dnsCfg, localInterfaceIP)
	if ret == nil {
		lastManualDNS = dnsCfg
	}
	return ret
}

// DeleteManual - reset manual DNS configuration to default (DHCP)
// 'localInterfaceIP' (obligatory only for Windows implementation) - local IP of VPN interface
func DeleteManual(localInterfaceIP net.IP) error {
	ret := implDeleteManual(localInterfaceIP)
	if ret == nil {
		lastManualDNS = DnsSettings{}
	}
	return ret
}

// GetLastManualDNS - returns information about current manual DNS
func GetLastManualDNS() DnsSettings {
	// TODO: get real DNS configuration of the OS
	return lastManualDNS
}

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

package dns

import (
	"fmt"
	"net"
	"strings"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
)

var log *logger.Logger
var lastManualDNS DnsSettings

func init() {
	log = logger.NewLogger("dns")
}

type DnsError struct {
	Err error
}

func (e *DnsError) Error() string {
	if e.Err == nil {
		return "DNS error"
	}
	return "DNS error: " + e.Err.Error()
}
func (e *DnsError) Unwrap() error { return e.Err }

func wrapErrorIfFailed(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*DnsError); ok {
		return err
	}
	return &DnsError{Err: err}
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
	return wrapErrorIfFailed(implInitialize())
}

// Pause pauses DNS (restore original DNS)
func Pause() error {
	return wrapErrorIfFailed(implPause())
}

// Resume resuming DNS (set DNS back which was before Pause)
func Resume(defaultDNS DnsSettings) error {
	return wrapErrorIfFailed(implResume(defaultDNS))
}

func EncryptionAbilities() (dnsOverHttps, dnsOverTls bool, err error) {
	dnsOverHttps, dnsOverTls, err = implGetDnsEncryptionAbilities()
	return dnsOverHttps, dnsOverTls, wrapErrorIfFailed(err)
}

// SetDefault set DNS configuration treated as default (non-manual) configuration
// 'dnsCfg' parameter - DNS configuration
// 'localInterfaceIP' (obligatory only for Windows implementation) - local IP of VPN interface
func SetDefault(dnsCfg DnsSettings, localInterfaceIP net.IP) error {
	ret := SetManual(dnsCfg, localInterfaceIP)
	if ret == nil {
		lastManualDNS = DnsSettings{}
	}
	return wrapErrorIfFailed(ret)
}

// SetManual - set manual DNS.
// 'dnsCfg' parameter - DNS configuration
// 'localInterfaceIP' (obligatory only for Windows implementation) - local IP of VPN interface
func SetManual(dnsCfg DnsSettings, localInterfaceIP net.IP) error {
	ret := implSetManual(dnsCfg, localInterfaceIP)
	if ret == nil {
		lastManualDNS = dnsCfg
	}
	return wrapErrorIfFailed(ret)
}

// DeleteManual - reset manual DNS configuration to default (DHCP)
// 'localInterfaceIP' (obligatory only for Windows implementation) - local IP of VPN interface
func DeleteManual(localInterfaceIP net.IP) error {
	ret := implDeleteManual(localInterfaceIP)
	if ret == nil {
		lastManualDNS = DnsSettings{}
	}
	return wrapErrorIfFailed(ret)
}

// GetLastManualDNS - returns information about current manual DNS
func GetLastManualDNS() DnsSettings {
	// TODO: get real DNS configuration of the OS
	return lastManualDNS
}

func GetPredefinedDnsConfigurations() ([]DnsSettings, error) {
	settings, err := implGetPredefinedDnsConfigurations()
	return settings, wrapErrorIfFailed(err)
}


package openvpn

import (
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/dns"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/vpn"
)

type platformSpecificProperties struct {
	manualDNS dns.DnsSettings
}

func (o *OpenVPN) implInit() error             { return nil }
func (o *OpenVPN) implIsCanUseParamsV24() bool { return true }

func (o *OpenVPN) implOnConnected() error {
	// on Windows it is not possible to change network interface properties until it not enabled
	// apply DNS value when VPN connected (TAP interface enabled)
	if !o.psProps.manualDNS.IsEmpty() {
		return dns.SetManual(o.psProps.manualDNS, o.clientIP)
	}

	// There could be manual-dns value saved from last connection in adapter properties. We must ensure that it erased.
	return dns.DeleteManual(o.clientIP)
}

func (o *OpenVPN) implOnDisconnected() error {
	return o.implOnResetManualDNS()
}

func (o *OpenVPN) implOnPause() error {
	// not in use in Windows implementation
	return nil
}

func (o *OpenVPN) implOnResume() error {
	// not in use in Windows implementation
	return nil
}

func (o *OpenVPN) implOnSetManualDNS(dnsCfg dns.DnsSettings) error {
	o.psProps.manualDNS = dnsCfg

	if o.state != vpn.CONNECTED {
		// on Windows it is not possible to change network interface properties until it not enabled
		// apply DNS value when VPN connected (TAP interface enabled)
	} else {
		return dns.SetManual(o.psProps.manualDNS, o.clientIP)
	}
	return nil
}

func (o *OpenVPN) implOnResetManualDNS() error {
	if !o.psProps.manualDNS.IsEmpty() {
		o.psProps.manualDNS = dns.DnsSettings{}
		return dns.DeleteManual(o.clientIP)
	}
	return nil
}

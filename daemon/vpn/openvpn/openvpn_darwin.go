
package openvpn

import (
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/dns"
)

type platformSpecificProperties struct {
	// no specific properties for macOS implementation
}

func (o *OpenVPN) implInit() error             { return nil }
func (o *OpenVPN) implIsCanUseParamsV24() bool { return true }

func (o *OpenVPN) implOnConnected() error {
	// not in use in macOS implementation
	return nil
}

func (o *OpenVPN) implOnDisconnected() error {
	// not in use in macOS implementation
	return nil
}

func (o *OpenVPN) implOnPause() error {
	return dns.Pause()
}

func (o *OpenVPN) implOnResume() error {
	return dns.Resume(dns.DnsSettings{})
}

func (o *OpenVPN) implOnSetManualDNS(dnsCfg dns.DnsSettings) error {
	return dns.SetManual(dnsCfg, nil)
}

func (o *OpenVPN) implOnResetManualDNS() error {
	return dns.DeleteManual(nil)
}

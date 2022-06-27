
package openvpn

import (
	"fmt"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/dns"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/platform/filerights"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/vpn"
)

type platformSpecificProperties struct {
	// no specific properties for Linux implementation
	isCanUseParamsV24 bool
}

func (o *OpenVPN) implInit() error {
	o.psProps.isCanUseParamsV24 = true

	if err := filerights.CheckFileAccessRightsExecutable(o.binaryPath); err != nil {
		return fmt.Errorf("error checking OpenVPN binary file: %w", err)
	}

	// Check OpenVPN minimum version
	minVer := []int{2, 3}
	verNums := GetOpenVPNVersion(o.binaryPath)
	log.Info("OpenVPN version:", verNums)
	for i := range minVer {
		if len(verNums) <= i {
			continue
		}
		if verNums[i] < minVer[i] {
			return fmt.Errorf("OpenVPN version '%v' not supported (minimum required version '%v')", verNums, minVer)
		}
	}
	if len(verNums) >= 2 && verNums[0] == 2 && verNums[1] < 4 {
		o.psProps.isCanUseParamsV24 = false
	}
	return nil
}

func (o *OpenVPN) implIsCanUseParamsV24() bool {
	return o.psProps.isCanUseParamsV24
}

func (o *OpenVPN) implOnConnected() error {
	// TODO: not implemented
	return nil
}

func (o *OpenVPN) implOnDisconnected() error {
	// TODO: not implemented
	return nil
}

func (o *OpenVPN) implOnPause() error {
	return dns.Pause()
}

func (o *OpenVPN) implOnResume() error {
	return dns.Resume(o.getDefaultDNS())
}

func (o *OpenVPN) implOnSetManualDNS(dnsCfg dns.DnsSettings) error {
	return dns.SetManual(dnsCfg, nil)
}

func (o *OpenVPN) implOnResetManualDNS() error {
	if o.IsPaused() == false {
		// restore default DNS pushed by OpenVPN server
		defaultDns := o.getDefaultDNS()
		if !defaultDns.IsEmpty() {
			return dns.SetManual(defaultDns, nil)
		}
	}

	return dns.DeleteManual(nil)
}

// getDefaultDNS returns default DNS pushed by OpenVPN server
func (o *OpenVPN) getDefaultDNS() dns.DnsSettings {
	mi := o.managementInterface
	if mi != nil && mi.isConnected && o.state != vpn.DISCONNECTED && o.state != vpn.EXITING {
		return dns.DnsSettings{DnsHost: mi.pushReplyDNS.String(), Encryption: dns.EncryptionNone}
	}
	return dns.DnsSettings{}
}

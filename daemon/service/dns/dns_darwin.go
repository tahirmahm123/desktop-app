
package dns

import (
	"fmt"
	"net"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/platform"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/shell"
)

// implInitialize doing initialization stuff (called on application start)
func implInitialize() error {
	return nil
}

func implPause() error {
	err := shell.Exec(log, platform.DNSScript(), "-pause")
	if err != nil {
		return fmt.Errorf("DNS pause: Failed to change DNS: %w", err)
	}
	return nil
}

// defaultDNS - not in use for darwin platfrom
func implResume(defaultDNS DnsSettings) error {
	err := shell.Exec(log, platform.DNSScript(), "-resume")
	if err != nil {
		return fmt.Errorf("DNS resume: Failed to change DNS: %w", err)
	}

	return nil
}

func implGetDnsEncryptionAbilities() (dnsOverHttps, dnsOverTls bool, err error) {
	return false, false, nil
}

// Set manual DNS.
// 'localInterfaceIP' - not in use for macOS implementation
func implSetManual(dnsCfg DnsSettings, localInterfaceIP net.IP) error {
	if dnsCfg.Encryption != EncryptionNone {
		return fmt.Errorf("DNS encryption is not supported on this platform")
	}

	err := shell.Exec(log, platform.DNSScript(), "-set_alternate_dns", dnsCfg.Ip().String())
	if err != nil {
		return fmt.Errorf("set manual DNS: Failed to change DNS: %w", err)
	}

	return nil
}

// DeleteManual - reset manual DNS configuration to default (DHCP)
// 'localInterfaceIP' (obligatory only for Windows implementation) - local IP of VPN interface
func implDeleteManual(localInterfaceIP net.IP) error {
	err := shell.Exec(log, platform.DNSScript(), "-delete_alternate_dns")
	if err != nil {
		return fmt.Errorf("reset manual DNS: Failed to change DNS: %w", err)
	}

	return nil
}

func implGetPredefinedDnsConfigurations() ([]DnsSettings, error) {
	return []DnsSettings{}, nil
}

// IsPrimaryInterfaceFound (macOS specific implementation) returns 'true' when networking is available (primary interface is available)
// When no networking available (WiFi off ?) - returns 'false'
// <this method in use by macOS:WireGuard implementation>
func IsPrimaryInterfaceFound() bool {
	err := shell.Exec(log, platform.DNSScript(), "-is_main_interface_detected")
	if err != nil {
		return false
	}
	return true
}



//go:build darwin && !debug
// +build darwin,!debug
package platform

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/platform/filerights"
)

func doOsInitForBuild() (warnings []string, errors []error) {
	// macOS-specific variable initialization
	firewallScript = "/Applications/VPN.app/Contents/Resources/etc/firewall.sh"
	dnsScript = "/Applications/VPN.app/Contents/Resources/etc/dns.sh"

	// common variables initialization
	settingsDir := "/Library/Application Support/VPN"
	settingsFile = path.Join(settingsDir, "settings.json")
	serversFile = path.Join(settingsDir, "servers.json")
	openvpnConfigFile = path.Join(settingsDir, "openvpn.cfg")
	openvpnProxyAuthFile = path.Join(settingsDir, "proxyauth.txt")
	wgConfigFilePath = path.Join(settingsDir, "wireguard.conf")

	openVpnBinaryPath = "/Applications/VPN.app/Contents/MacOS/openvpn"
	openvpnCaKeyFile = "/Applications/VPN.app/Contents/Resources/etc/ca.crt"
	openvpnTaKeyFile = "/Applications/VPN.app/Contents/Resources/etc/ta.key"
	openvpnUpScript = "/Applications/VPN.app/Contents/Resources/etc/dns.sh -up"
	openvpnDownScript = "/Applications/VPN.app/Contents/Resources/etc/dns.sh -down"

	obfsproxyStartScript = "/Applications/VPN.app/Contents/Resources/obfsproxy/obfs4proxy"

	wgBinaryPath = "/Applications/VPN.app/Contents/MacOS/WireGuard/wireguard-go"
	wgToolBinaryPath = "/Applications/VPN.app/Contents/MacOS/WireGuard/wg"

	return nil, nil
}

func doInitOperations() (w string, e error) {
	serversFile := ServersFile()
	if _, err := os.Stat(serversFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File '%s' does not exists. Copying from bundle...\n", serversFile)
			// Servers file is not exists on required place
			// Probably, it is first start after clean install
			// Copying it from a bundle
			os.MkdirAll(filepath.Base(serversFile), os.ModePerm)
			if _, err = copyFile("/Applications/VPN.app/Contents/Resources/etc/servers.json", serversFile); err != nil {
				return err.Error(), nil
			}
			return "", nil
		}

		return err.Error(), nil
	}
	return "", nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	destination.Chmod(filerights.DefaultFilePermissionsForConfig())
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

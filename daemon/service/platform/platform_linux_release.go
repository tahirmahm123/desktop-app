

//go:build linux && !debug
// +build linux,!debug
package platform

import (
	"path"
)

func doOsInitForBuild() (warnings []string, errors []error) {
	installDir := "/opt/vpn"

	firewallScript = path.Join(installDir, "etc/firewall.sh")
	splitTunScript = path.Join(installDir, "etc/splittun.sh")
	openvpnCaKeyFile = path.Join(installDir, "etc/ca.crt")
	openvpnTaKeyFile = path.Join(installDir, "etc/ta.key")
	openvpnUpScript = path.Join(installDir, "etc/client.up")
	openvpnDownScript = path.Join(installDir, "etc/client.down")
	serversFile = path.Join(installDir, "etc/servers.json")

	obfsproxyStartScript = path.Join(installDir, "obfsproxy/obfs4proxy")

	wgBinaryPath = path.Join(installDir, "wireguard-tools/wg-quick")
	wgToolBinaryPath = path.Join(installDir, "wireguard-tools/wg")

	settingsFile = path.Join(tmpDir, "settings.json")
	openvpnConfigFile = path.Join(tmpDir, "openvpn.cfg")
	openvpnProxyAuthFile = path.Join(tmpDir, "proxyauth.txt")
	wgConfigFilePath = path.Join(tmpDir, "wgvpn.conf")

	return nil, nil
}



//go:build linux && debug
// +build linux,debug
package platform

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func doOsInitForBuild() (warnings []string, errors []error) {
	installDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("Failed to obtain folder of current binary: %s", err.Error()))
	}

	if len(os.Args) > 2 {
		firstArg := strings.Split(os.Args[1], "=")
		if len(firstArg) == 2 && firstArg[0] == "-debug_install_dir" {
			installDir = firstArg[1]
		}
	}

	// When running tests, the installDir is detected as a dir where test located
	// we need to point installDir to project root
	// Therefore, we cutting rest after "desktop-app/daemon"
	rootDir := "desktop-app/daemon"
	if idx := strings.LastIndex(installDir, rootDir); idx > 0 {
		installDir = installDir[:idx+len(rootDir)]
	}

	installDir = path.Join(installDir, "References/Linux")

	firewallScript = path.Join(installDir, "etc/firewall.sh")
	splitTunScript = path.Join(installDir, "etc/splittun.sh")
	openvpnCaKeyFile = path.Join(installDir, "etc/ca.crt")
	openvpnTaKeyFile = path.Join(installDir, "etc/ta.key")
	openvpnUpScript = path.Join(installDir, "etc/client.up")
	openvpnDownScript = path.Join(installDir, "etc/client.down")
	serversFile = path.Join(installDir, "etc/servers.json")

	obfsproxyStartScript = path.Join(installDir, "_deps/obfs4proxy_inst/obfs4proxy")

	wgBinaryPath = path.Join(installDir, "_deps/wireguard-tools_inst/wg-quick")
	wgToolBinaryPath = path.Join(installDir, "_deps/wireguard-tools_inst/wg")

	settingsFile = path.Join(tmpDir, "settings.json")
	openvpnConfigFile = path.Join(tmpDir, "openvpn.cfg")
	openvpnProxyAuthFile = path.Join(tmpDir, "proxyauth.txt")
	wgConfigFilePath = path.Join(tmpDir, "wgvpn.conf")

	return nil, nil
}

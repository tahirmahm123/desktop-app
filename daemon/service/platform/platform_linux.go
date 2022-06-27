
package platform

import (
	"path"
)

var (
	firewallScript string
	splitTunScript string
	logDir         string = "/opt/tahirmahm123/vpn-log"
	tmpDir         string = "/opt/tahirmahm123/vpn-mutable"
)

// initialize all constant values (e.g. servicePortFile) which can be used in external projects (VPN
func doInitConstants() {
	fwInitialValueAllowApiServers = false
	servicePortFile = path.Join(tmpDir, "port.txt")

	logFile = path.Join(logDir, "VPNlog")
	openvpnLogFile = path.Join(logDir, "openvpn.log")

	openvpnUserParamsFile = path.Join(tmpDir, "ovpn_extra_params.txt")
}

func doOsInit() (warnings []string, errors []error) {
	openVpnBinaryPath = path.Join("/usr/sbin", "openvpn")
	routeCommand = "/sbin/ip route"

	warnings, errors = doOsInitForBuild()

	if errors == nil {
		errors = make([]error, 0)
	}

	if err := checkFileAccessRightsExecutable("firewallScript", firewallScript); err != nil {
		errors = append(errors, err)
	}
	if err := checkFileAccessRightsExecutable("splitTunScript", splitTunScript); err != nil {
		errors = append(errors, err)
	}

	return warnings, errors
}

func doInitOperations() (w string, e error) { return "", nil }

// FirewallScript returns path to firewal script
func FirewallScript() string {
	return firewallScript
}

// SplitTunScript returns path to script which control split-tunneling functionality
func SplitTunScript() string {
	return splitTunScript
}

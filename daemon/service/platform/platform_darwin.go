
package platform

import (
	"path"
)

var (
	firewallScript string
	dnsScript      string
)

// initialize all constant values (e.g. servicePortFile) which can be used in external projects (VPN
func doInitConstants() {
	fwInitialValueAllowApiServers = false
	servicePortFile = "/Library/Application Support/VPNxt"
	openvpnUserParamsFile = "/Library/Application Support/VPNN/ovpn_extra_params.txt"

	logDir := "/Library/Logs/"
	logFile = path.Join(logDir, "VPNlog")
	openvpnLogFile = path.Join(logDir, "openvpn.log")
}

func doOsInit() (warnings []string, errors []error) {
	routeCommand = "/sbin/route"

	warnings, errors = doOsInitForBuild()

	if errors == nil {
		errors = make([]error, 0)
	}

	if err := checkFileAccessRightsExecutable("firewallScript", firewallScript); err != nil {
		errors = append(errors, err)
	}
	if err := checkFileAccessRightsExecutable("dnsScript", dnsScript); err != nil {
		errors = append(errors, err)
	}

	return warnings, errors
}

// FirewallScript returns path to firewal script
func FirewallScript() string {
	return firewallScript
}

// DNSScript returns path to DNS script
func DNSScript() string {
	return dnsScript
}

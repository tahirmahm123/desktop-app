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

//go:build windows && !debug
// +build windows,!debug

package platform

import (
	"fmt"
	"path"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// initialize all constant values (e.g. servicePortFile) which can be used in external projects (VPN
func doInitConstantsForBuild() {
}

func doOsInitForBuild() {
	installDir := getInstallDir()
	wfpDllPath = path.Join(installDir, "VPNll Native x64.dll")
	nativeHelpersDllPath = path.Join(installDir, "VPN Helpers Native x64.dll")
	splitTunDriverPath = path.Join(installDir, "SplitTunnelDriver", "x86_64", "vpn-split-tunnel.sys")
	if !Is64Bit() {
		wfpDllPath = path.Join(installDir, "VPNll Native.dll")
		nativeHelpersDllPath = path.Join(installDir, "VPN Helpers Native.dll")
	}
}

func getInstallDir() string {
	ret := ""

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\VPN`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	defer k.Close()

	if err == nil {
		ret, _, err = k.GetStringValue("")
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
	}

	if len(ret) == 0 {
		fmt.Println("WARNING: There is no info about VPN install folder in the registry. Is VPN CliVPNled?")
		return ""
	}

	return strings.ReplaceAll(ret, `\`, `/`)
}

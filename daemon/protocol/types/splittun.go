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

package types

import (
	"github.com/tahirmahm123/vpn-desktop-app/daemon/oshelpers"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/splittun"
)

// GetInstalledApps (request) requests information about installed applications on the system
type GetInstalledApps struct {
	CommandBase
	// (optional) Platform-depended: extra parameters (in JSON)
	// For Windows:
	//			WindowsEnvAppdata 	string
	// 				Applicable only for Windows: APPDATA environment variable
	// 				Needed to know path of current user's (not root) StartMenu folder location
	// For Linux:
	//			EnvVar_XDG_CURRENT_DESKTOP string
	//			EnvVar_XDG_DATA_DIRS       string
	//			EnvVar_HOME                string
	//			IconsTheme                 string
	ExtraArgsJSON string
}

// InstalledAppsResp (response) contains information about installed applications on the system
type InstalledAppsResp struct {
	CommandBase
	Apps []oshelpers.AppInfo
}

// GetAppIcon (request) requests shell icon for binary file (application)
// Note: ensure if SplitTunnelStatus.IsCanGetAppIconForBinary is active
type GetAppIcon struct {
	CommandBase
	AppBinaryPath string
}

// AppIconResp (response) contains information about shell icon for binary file (application)
type AppIconResp struct {
	CommandBase
	AppBinaryPath string
	AppIcon       string // base64 png image
}

// SplitTunnelSet (request) sets the split-tunnelling configuration
type SplitTunnelSetConfig struct {
	CommandBase
	IsEnabled bool // is ST enabled
	Reset     bool // disable ST and erase all ST config
}

// GetSplitTunnelStatus (request) requests the Split-Tunnelling configuration
type SplitTunnelGetStatus struct {
	CommandBase
}

// SplitTunnelStatus (response) returns the split-tunnelling configuration
type SplitTunnelStatus struct {
	CommandBase
	// is ST enabled
	IsEnabled                   bool
	IsFunctionalityNotAvailable bool
	// This parameter informs availability of the functionality to get icon for particular binary
	// (true - if commands GetAppIcon/AppIconResp  applicable for this platform)
	IsCanGetAppIconForBinary bool
	// Information about applications added to ST configuration
	// (applicable for Windows)
	SplitTunnelApps []string
	// Information about active applications running in Split-Tunnel environment
	// (applicable for Linux)
	RunningApps []splittun.RunningApp
}

// SplitTunnelAddApp (request) add application to SplitTunneling
// Expected response:
// 		Windows	- types.EmptyResp (success)
//  	Linux	- SplitTunnelAddAppCmdResp -> contains shell command which have to be executed in user space environment
//
// Description of Split Tunneling commands sequence to run the application:
//	[client]					[daemon]
//	SplitTunnelAddApp		->
//							<-	windows:	types.EmptyResp (success)
//							<-	linux:		types.SplitTunnelAddAppCmdResp (some operations required on client side)
//	<windows: done>
// 	<execute shell command: types.SplitTunnelAddAppCmdResp.CmdToExecute and get PID>
//  SplitTunnelAddedPidInfo	->
// 							<-	types.EmptyResp (success)
type SplitTunnelAddApp struct {
	CommandBase
	// Windows: full path to the app binary
	// Linux: command to be executed in ST environment (e.g. binary + arguments)
	Exec string
}

// SplitTunnelAddAppCmdResp (response) contains shell command which have to be executed in user space environment
// (not in use for Windows platform)
type SplitTunnelAddAppCmdResp struct {
	CommandBase
	// Command will be executed in ST environment
	// (identical to SplitTunnelAddApp.Exec)
	Exec string
	// Shell command which have to be executed in user space environment
	CmdToExecute string

	IsAlreadyRunning        bool
	IsAlreadyRunningMessage string
}

// SplitTunnelAddedPidInfo (request) informs the daemon about started process in ST environment
// (not in use for Windows platform)
type SplitTunnelAddedPidInfo struct {
	CommandBase
	Pid int
	// Command will be executed in ST environment (e.g. binary + arguments)
	// (identical to SplitTunnelAddApp.Exec and SplitTunnelAddAppCmdResp.Exec)
	Exec string
	// Shell command used to perform this operation
	CmdToExecute string
}

type SplitTunnelRemoveApp struct {
	CommandBase
	// (applicable for Linux) PID of the running process in ST environment
	Pid int
	// (applicable for Windows) full path to the app binary to be excluded from ST
	Exec string
}

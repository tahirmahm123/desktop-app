//
//  Daemon for VPN Client Desktop
//  https://github.com/tahirmahm123/vpn-desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the Daemon for VPN Client Desktop.
//
//  The Daemon for VPN Client Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The Daemon for VPN Client Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the Daemon for VPN Client Desktop. If not, see <https://www.gnu.org/licenses/>.
//

//go:build darwin
// +build darwin

package libvpn

/*
#include <libvpn.h>
*/
import (
	"C"
)

import (
	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
)

// TODO: reimplement accessing libvpn using syscall.NewLazyDLL+NewProc (avoid using CGO)

// Unload - unload (uninitialize\close) 'libvpn' shared library
func Unload() {
	C.UnLoadLibrary()
}

// StartXpcListener starts listener for helper
func StartXpcListener(tcpPort int, secret uint64) {

	ret := C.start_xpc_listener(C.CString("net.vpn.client.Helper"), C.int(tcpPort), C.uint64_t(secret))
	if ret == 0 {
		return
	}

	switch ret {
	case C.ERROR_LIB_NOT_FOUND:
		logger.Panic("Unable to find dynamic library")
	case C.ERROR_METHOD_NOT_FOUND:
		logger.Panic("Method was not found in dynamic library")
	default:
		logger.Panic("Unexpected error: ", ret)
	}
}

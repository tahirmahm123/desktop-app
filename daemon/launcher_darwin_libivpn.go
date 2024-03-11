//go:build darwin && libivpn
// +build darwin,libivpn

package main

import (
	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/oshelpers/macos/libivpn"
)

// inform OS-specific implementation about listener port
func implStartedOnPort(openedPort int, secret uint64) {
	libivpn.StartXpcListener(openedPort, secret)
}

// OS-specific service finalizer
func implStopped() {
	// do not forget to close 'libivpn' dynamic library
	logger.Debug("Unloading libivpn...")
	libivpn.Unload()
	logger.Debug("Unloaded libivpn")
}

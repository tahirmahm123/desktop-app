//go:build darwin && libvpn
// +build darwin,libvpn

package launcher

import (
	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/oshelpers/macos/libvpn"
)

// inform OS-specific implementation about listener port
func implStartedOnPort(openedPort int, secret uint64) {
	libvpn.StartXpcListener(openedPort, secret)
}

// OS-specific service finalizer
func implStopped() {
	// do not forget to close 'libvpn' dynamic library
	logger.Debug("Unloading libvpn...")
	libvpn.Unload()
	logger.Debug("Unloaded libvpn")
}

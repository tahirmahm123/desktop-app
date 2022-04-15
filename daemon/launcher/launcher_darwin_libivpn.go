//go:build darwin && libivpn
// +build darwin,libivpn

package launcher

import (
	"github.com/ivpn/desktop-app/daemon/logger"
	"github.com/ivpn/desktop-app/daemon/oshelpers/macos/libivpn"
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

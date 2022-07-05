
package launcher

import (
	"fmt"
	"os"
	"path"
)

// Prepare to start VPN for macOS
func doPrepareToRun() error {
	// create symlink to 'vpn' cli client
	binFolder := "/usr/local/bin"           // "/usr/local/bin"
	linkpath := path.Join(binFolder, "vpn") // "/usr/local/bin/vpn"
	if _, err := os.Stat(linkpath); os.IsNotExist(err) {
		// "/usr/local/bin"
		if _, err := os.Stat(binFolder); os.IsNotExist(err) {
			log.Info(fmt.Sprintf("Folder '%s' not exists. Creating it...", binFolder))
			if err = os.Mkdir(binFolder, 0775); err != nil {
				log.Error(fmt.Sprintf("Failed to create folder '%s': ", binFolder), err)
			}
		}
	}
	return nil
}

// inform OS-specific implementation about listener port
func doStartedOnPort(openedPort int, secret uint64) {
	implStartedOnPort(openedPort, secret)
}

// OS-specific service finalizer
func doStopped() {
	implStopped()
}

func isNeedToSavePortInFile() bool {
	return true
}

// checkIsAdmin - check is application running with root privileges
func doCheckIsAdmin() bool {
	uid := os.Geteuid()
	if uid != 0 {
		return false
	}

	return true
}


package service

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) implPingServersStarting(hosts []net.IP) error {
	// nothing to do for Windows implementation
	// firewall configured to allow all connectivity for service
	return nil
}
func (s *Service) implPingServersStopped(hosts []net.IP) error {
	// nothing to do for Windows implementation
	// firewall configured to allow all connectivity for service
	return nil
}

func (s *Service) implSplitTunnelling_AddApp(binaryFile string) (requiredCmdToExec string, isAlreadyRunning bool, err error) {
	binaryFile = strings.TrimSpace(binaryFile)
	if len(binaryFile) <= 0 {
		return "", false, nil
	}

	prefs := s._preferences
	// current binary folder path
	var exeDir string
	if ex, err := os.Executable(); err == nil && len(ex) > 0 {
		exeDir = filepath.Dir(ex)
	}

	// Ensure no binaries from VPNe is included into apps list to Split-Tunnel
	if strings.HasPrefix(binaryFile, exeDir) {
		return "", false, fmt.Errorf("Split-Tunnelling for VPNes is forbidden (%s)", binaryFile)
	}
	// Ensure file is exists
	if _, err := os.Stat(binaryFile); os.IsNotExist(err) {
		return "", false, err
	}

	binaryPathLowCase := strings.ToLower(binaryFile)
	for _, a := range prefs.SplitTunnelApps {
		if strings.ToLower(a) == binaryPathLowCase {
			// the binary is already in configuration
			return "", false, nil
		}
	}

	prefs.SplitTunnelApps = append(prefs.SplitTunnelApps, binaryFile)
	s.setPreferences(prefs)

	return "", false, nil
}

func (s *Service) implSplitTunnelling_RemoveApp(pid int, binaryPath string) (err error) {
	binaryPath = strings.TrimSpace(binaryPath)
	if len(binaryPath) <= 0 {
		return nil
	}

	prefs := s._preferences
	newStApps := make([]string, 0, len(prefs.SplitTunnelApps))
	binaryPathLowCase := strings.ToLower(binaryPath)

	for _, a := range prefs.SplitTunnelApps {
		if strings.ToLower(a) == binaryPathLowCase {
			continue
		}
		newStApps = append(newStApps, a)
	}

	prefs.SplitTunnelApps = newStApps
	s.setPreferences(prefs)

	return nil
}

func (s *Service) implSplitTunnelling_AddedPidInfo(pid int, exec string, cmdToExecute string) error {
	return fmt.Errorf("function not applicable for this platform")
}

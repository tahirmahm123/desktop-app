
package service

import (
	"fmt"
	"net"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/firewall"
)

func (s *Service) implPingServersStarting(hosts []net.IP) error {
	const onlyForICMP = true
	const isPersistent = false
	return firewall.AddHostsToExceptions(hosts, onlyForICMP, isPersistent)
}
func (s *Service) implPingServersStopped(hosts []net.IP) error {
	const onlyForICMP = true
	const isPersistent = false
	return firewall.RemoveHostsFromExceptions(hosts, onlyForICMP, isPersistent)
}

func (s *Service) implSplitTunnelling_AddApp(binaryFile string) (requiredCmdToExec string, isAlreadyRunning bool, err error) {
	// Split Tunneling is not implemented for macOS
	return "", false, nil
}
func (s *Service) implSplitTunnelling_RemoveApp(pid int, binaryPath string) (err error) {
	// Split Tunneling is not implemented for macOS
	return nil
}
func (s *Service) implSplitTunnelling_AddedPidInfo(pid int, exec string, cmdToExecute string) error {
	return fmt.Errorf("function not applicable for this platform")
}

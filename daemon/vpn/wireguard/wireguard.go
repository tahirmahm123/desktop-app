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

package wireguard

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/helpers"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/netinfo"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/dns"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/vpn"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("wg")
}

// ConnectionParams contains all information to make new connection
type ConnectionParams struct {
	clientLocalIP    net.IP
	clientPrivateKey string
	hostPort         int
	hostIP           net.IP
	hostPublicKey    string
	hostLocalIP      net.IP
	ipv6Prefix       string

	// multihopExitSrvID (geteway ID) just in use to keep clients notified about connected MH exit server
	// in same manner as for OpenVPN connection.
	// Example: "gateway":"zz.wg.vpn.net" => "zz"
	multihopExitSrvID string
}

func (cp *ConnectionParams) GetIPv6ClientLocalIP() net.IP {
	if len(cp.ipv6Prefix) <= 0 {
		return nil
	}
	return net.ParseIP(cp.ipv6Prefix + cp.clientLocalIP.String())
}
func (cp *ConnectionParams) GetIPv6HostLocalIP() net.IP {
	if len(cp.ipv6Prefix) <= 0 {
		return nil
	}
	return net.ParseIP(cp.ipv6Prefix + cp.hostLocalIP.String())
}

// SetCredentials update WG credentials
func (cp *ConnectionParams) SetCredentials(privateKey string, localIP net.IP) {
	cp.clientPrivateKey = privateKey
	cp.clientLocalIP = localIP
}

// CreateConnectionParams initializing connection parameters object
func CreateConnectionParams(
	multihopExitSrvID string,
	hostPort int,
	hostIP net.IP,
	hostPublicKey string,
	hostLocalIP net.IP,
	ipv6Prefix string) ConnectionParams {

	return ConnectionParams{
		multihopExitSrvID: multihopExitSrvID,
		hostPort:          hostPort,
		hostIP:            hostIP,
		hostPublicKey:     hostPublicKey,
		hostLocalIP:       hostLocalIP,
		ipv6Prefix:        ipv6Prefix}
}

// WireGuard structure represents all data of wireguard connection
type WireGuard struct {
	binaryPath     string
	toolBinaryPath string
	configFilePath string
	connectParams  ConnectionParams
	localPort      int

	// Must be implemented (AND USED) in correspond file for concrete platform. Must contain platform-specified properties (or can be empty struct)
	internals internalVariables
}

// NewWireGuardObject creates new wireguard structure
func NewWireGuardObject(wgBinaryPath string, wgToolBinaryPath string, wgConfigFilePath string, connectionParams ConnectionParams) (*WireGuard, error) {
	if connectionParams.clientLocalIP == nil || len(connectionParams.clientPrivateKey) == 0 {
		return nil, fmt.Errorf("WireGuard local credentials not defined")
	}

	return &WireGuard{
		binaryPath:     wgBinaryPath,
		toolBinaryPath: wgToolBinaryPath,
		configFilePath: wgConfigFilePath,
		connectParams:  connectionParams}, nil
}

// DestinationIP -  Get destination IP (VPN host server or proxy server IP address)
// This information if required, for example, to allow this address in firewall
func (wg *WireGuard) DestinationIP() net.IP {
	return wg.connectParams.hostIP
}

// Type just returns VPN type
func (wg *WireGuard) Type() vpn.Type { return vpn.WireGuard }

// Init performs basic initializations before connection
// It is useful, for example:
//	- for WireGuard(Windows) - to ensure that WG service is fully uninstalled
//	- for OpenVPN(Linux) - to ensure that OpenVPN has correct version
func (wg *WireGuard) Init() error {
	return wg.init()
}

// Connect - SYNCHRONOUSLY execute openvpn process (wait until it finished)
func (wg *WireGuard) Connect(stateChan chan<- vpn.StateInfo) error {

	disconnectDescription := ""

	stateChan <- vpn.NewStateInfo(vpn.CONNECTING, "")
	defer func() {
		stateChan <- vpn.NewStateInfo(vpn.DISCONNECTED, disconnectDescription)
	}()

	err := wg.connect(stateChan)

	if err != nil {
		disconnectDescription = err.Error()
	}

	return err
}

// Disconnect stops the connection
func (wg *WireGuard) Disconnect() error {
	return wg.disconnect()
}

// IsPaused checking if we are in paused state
func (wg *WireGuard) IsPaused() bool {
	return wg.isPaused()
}

// Pause doing required operation for Pause (temporary restoring default DNS)
func (wg *WireGuard) Pause() error {
	return wg.pause()
}

// Resume doing required operation for Resume (restores DNS configuration before Pause)
func (wg *WireGuard) Resume() error {
	return wg.resume()
}

// SetManualDNS changes DNS to manual IP
func (wg *WireGuard) SetManualDNS(dnsCfg dns.DnsSettings) error {
	return wg.setManualDNS(dnsCfg)
}

// ResetManualDNS restores DNS
func (wg *WireGuard) ResetManualDNS() error {
	return wg.resetManualDNS()
}

func (wg *WireGuard) generateAndSaveConfigFile(cfgFilePath string) error {
	cfg, err := wg.generateConfig()
	if err != nil {
		return fmt.Errorf("failed to generate WireGuard configuration: %w", err)
	}

	// write configuration into temporary file
	configText := strings.Join(cfg, "\n")

	err = ioutil.WriteFile(cfgFilePath, []byte(configText), 0600)
	if err != nil {
		return fmt.Errorf("failed to save WireGuard configuration into a file: %w", err)
	}

	log.Info("WireGuard  configuration:",
		"\n=====================\n",
		strings.ReplaceAll(configText, wg.connectParams.clientPrivateKey, "***"),
		"\n=====================\n")

	return nil
}

func (wg *WireGuard) generateConfig() ([]string, error) {
	localPort, err := netinfo.GetFreeUDPPort()
	if err != nil {
		return nil, fmt.Errorf("unable to obtain free local port: %w", err)
	}

	wg.localPort = localPort

	// prevent user-defined data injection: ensure that nothing except the base64 public key will be stored in the configuration
	if !helpers.ValidateBase64(wg.connectParams.hostPublicKey) {
		return nil, fmt.Errorf("WG public key is not base64 string")
	}
	if !helpers.ValidateBase64(wg.connectParams.clientPrivateKey) {
		return nil, fmt.Errorf("WG private key is not base64 string")
	}

	interfaceCfg := []string{
		"[Interface]",
		"PrivateKey = " + wg.connectParams.clientPrivateKey,
		"ListenPort = " + strconv.Itoa(wg.localPort)}

	peerCfg := []string{
		"[Peer]",
		"PublicKey = " + wg.connectParams.hostPublicKey,
		"Endpoint = " + wg.connectParams.hostIP.String() + ":" + strconv.Itoa(wg.connectParams.hostPort),
		"PersistentKeepalive = 25"}

	// add some OS-specific configurations (if necessary)
	iCfg, pCgf := wg.getOSSpecificConfigParams()
	interfaceCfg = append(interfaceCfg, iCfg...)
	peerCfg = append(peerCfg, pCgf...)

	return append(interfaceCfg, peerCfg...), nil
}

func (wg *WireGuard) notifyConnectedStat(stateChan chan<- vpn.StateInfo) {
	const isTCP = false
	const isCanPause = true

	si := vpn.NewStateInfoConnected(
		isTCP,
		wg.connectParams.clientLocalIP,
		wg.connectParams.GetIPv6ClientLocalIP(),
		wg.localPort,
		wg.connectParams.hostIP,
		wg.connectParams.hostPort,
		isCanPause)

	si.ExitServerID = wg.connectParams.multihopExitSrvID

	stateChan <- si
}

func (wg *WireGuard) OnRoutingChanged() error {
	return wg.onRoutingChanged()
}

func (wg *WireGuard) IsIPv6InTunnel() bool {
	return len(wg.connectParams.GetIPv6ClientLocalIP()) > 0
}

//
//  Daemon for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the Daemon for IVPN Client Desktop.
//
//  The Daemon for IVPN Client Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The Daemon for IVPN Client Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the Daemon for IVPN Client Desktop. If not, see <https://www.gnu.org/licenses/>.
//

package service

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ivpn/desktop-app/daemon/api"
	"github.com/ivpn/desktop-app/daemon/api/types"
	"github.com/ivpn/desktop-app/daemon/helpers"
	"github.com/ivpn/desktop-app/daemon/logger"
	"github.com/ivpn/desktop-app/daemon/netinfo"
	"github.com/ivpn/desktop-app/daemon/oshelpers"
	protocolTypes "github.com/ivpn/desktop-app/daemon/protocol/types"
	"github.com/ivpn/desktop-app/daemon/service/dns"
	"github.com/ivpn/desktop-app/daemon/service/firewall"
	"github.com/ivpn/desktop-app/daemon/service/platform"
	"github.com/ivpn/desktop-app/daemon/service/platform/filerights"
	"github.com/ivpn/desktop-app/daemon/service/preferences"
	"github.com/ivpn/desktop-app/daemon/service/srverrors"
	"github.com/ivpn/desktop-app/daemon/splittun"
	"github.com/ivpn/desktop-app/daemon/vpn"
	"github.com/ivpn/desktop-app/daemon/vpn/openvpn"
	"github.com/ivpn/desktop-app/daemon/vpn/wireguard"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("servc")
}

// RequiredState VPN state which service is going to reach
type RequiredState int

// Requested VPN states
const (
	Disconnect     RequiredState = 0
	Connect        RequiredState = 1
	KeepConnection RequiredState = 2
)

const (
	// SessionCheckInterval - the interval for periodical check session status
	SessionCheckInterval time.Duration = time.Hour * 6
)

// Service - IVPN service
type Service struct {
	_evtReceiver       IServiceEventsReceiver
	_api               *api.API
	_serversUpdater    IServersUpdater
	_netChangeDetector INetChangeDetector
	_wgKeysMgr         IWgKeysManager
	_vpn               vpn.Process
	_preferences       preferences.Preferences
	_connectMutex      sync.Mutex

	// Additional information about current VPN connection
	// Use GetVpnSessionInfo()/SetVpnSessionInfo() to access this data
	_vpnSessionInfo      VpnSessionInfo
	_vpnSessionInfoMutex sync.Mutex

	// manual DNS value (if not defined - nil)
	_manualDNS dns.DnsSettings

	// Required VPN state which service is going to reach (disconnect->keep connection->connect)
	// When KeepConnection - reconnects immediately after disconnection
	_requiredVpnState RequiredState

	// Note: Disconnect() function will wait until VPN fully disconnects
	_done chan struct{}

	_isServersPingInProgress bool

	// nil - when session checker stopped
	// to stop -> write to channel (it is synchronous channel)
	_sessionCheckerStopChn chan struct{}
}

// VpnSessionInfo - Additional information about current VPN connection
type VpnSessionInfo struct {
	// The outbound IP addresses on the moment BEFORE the VPN connection
	OutboundIPv4 net.IP
	OutboundIPv6 net.IP
	// local VPN addresses (outbound IPs)
	VpnLocalIPv4 net.IP
	VpnLocalIPv6 net.IP
}

// CreateService - service constructor
func CreateService(evtReceiver IServiceEventsReceiver, api *api.API, updater IServersUpdater, netChDetector INetChangeDetector, wgKeysMgr IWgKeysManager) (*Service, error) {
	if updater == nil {
		return &Service{}, fmt.Errorf("ServersUpdater is not defined")
	}

	serv := &Service{
		_preferences:       *preferences.Create(),
		_evtReceiver:       evtReceiver,
		_api:               api,
		_serversUpdater:    updater,
		_netChangeDetector: netChDetector,
		_wgKeysMgr:         wgKeysMgr}

	// register the current service as a 'Connectivity checker' for API object
	serv._api.SetConnectivityChecker(serv)

	if err := serv.init(); err != nil {
		return nil, fmt.Errorf("service initialization error : %w", err)
	}

	return serv, nil
}

func (s *Service) init() error {
	if err := s._preferences.LoadPreferences(); err != nil {
		log.Error("Failed to load service preferences: ", err)

		log.Warning("Saving default values for preferences")
		s._preferences.SavePreferences()
	}

	if err := dns.Initialize(); err != nil {
		log.Error(fmt.Sprintf("failed to initialize DNS : %s", err))
	}

	if err := firewall.Initialize(); err != nil {
		return fmt.Errorf("service initialization error : %w", err)
	}

	if err := splittun.Initialize(); err != nil {
		log.Warning(fmt.Errorf("Split-Tunnelling initialization error : %w", err))
	} else {
		// apply Split Tunneling configuration
		s.splitTunnelling_ApplyConfig()
	}

	// Logging mus be already initialized (by launcher). Do nothing here.
	// Init logger (if not initialized before)
	//logger.Enable(s._preferences.IsLogging)

	// Init firewall
	if err := firewall.AllowLAN(s._preferences.IsFwAllowLAN, s._preferences.IsFwAllowLANMulticast); err != nil {
		log.Error("Failed to initialize firewall with AllowLAN preference value: ", err)
	}

	if s._preferences.IsFwPersistant {
		log.Info("Enabling firewal (persistant configuration)")
		if err := firewall.SetPersistant(true); err != nil {
			log.Error("Failed to enable firewall: ", err)
		}
	}

	// start WireGuard keys rotation
	if err := s._wgKeysMgr.Init(s); err != nil {
		log.Error("Failed to initialize WG keys rotation:", err)
	} else {
		if err := s._wgKeysMgr.StartKeysRotation(); err != nil {
			log.Error("Failed to start WG keys rotation:", err)
		}
	}

	if err := s.initWiFiFunctionality(); err != nil {
		log.Error("Failed to init WiFi functionality:", err)
	}

	// Check session status (start as go-routine to do not block service initialization)
	go s.RequestSessionStatus()
	// Start session status checker
	s.startSessionChecker()

	s.updateAPIAddrInFWExceptions()
	// servers updated notifier
	go func() {
		session := s.Preferences().Session
		if session.IsLoggedIn() {
			s._serversUpdater.startService(session.Session)
			defer func() {
				if r := recover(); r != nil {
					log.Error("PANIC in Servers update notifier!: ", r)
					if err, ok := r.(error); ok {
						log.ErrorTrace(err)
					}
				}
			}()
			log.Info("Servers update notifier started")
			for {
				// wait for 'servers updated' event
				<-s._serversUpdater.UpdateNotifierChannel()
				// notify clients
				servers, err := s._serversUpdater.GetServers()
				if err != nil {
					s._evtReceiver.OnServersUpdated(servers)
				}
			}
		}
	}()

	return nil
}

func (s *Service) IsConnectivityBlocked() (isBlocked bool, reasonDescription string, err error) {
	preferences := s._preferences
	if !preferences.IsFwAllowApiServers &&
		preferences.Session.IsLoggedIn() &&
		(!s.Connected() || s.IsPaused()) {
		enabled, err := s.FirewallEnabled()
		if err == nil && enabled {
			return true, "Access to IVPN servers is blocked (check IVPN Firewall settings)", nil
		}
	}
	return false, "", nil
}

func (s *Service) GetVpnSessionInfo() VpnSessionInfo {
	s._vpnSessionInfoMutex.Lock()
	defer s._vpnSessionInfoMutex.Unlock()
	return s._vpnSessionInfo
}

func (s *Service) SetVpnSessionInfo(i VpnSessionInfo) {
	s._vpnSessionInfoMutex.Lock()
	defer s._vpnSessionInfoMutex.Unlock()
	s._vpnSessionInfo = i
}

func (s *Service) updateAPIAddrInFWExceptions() {
	/*svrs, _ := s.ServersList()

	ivpnAPIAddr := svrs.Config.API.IPAddresses

	if len(ivpnAPIAddr) <= 0 {
		return
	}

	apiAddrs := make([]net.IP, 0, len(ivpnAPIAddr))
	for _, ipStr := range ivpnAPIAddr {
		apiIP := net.ParseIP(ipStr)
		if apiIP != nil {
			apiAddrs = append(apiAddrs, apiIP)
		}
	}

	if len(apiAddrs) > 0 {
		const onlyForICMP = false
		const isPersistent = true
		prefs := s.Preferences()
		if prefs.IsFwAllowApiServers {
			firewall.AddHostsToExceptions(apiAddrs, onlyForICMP, isPersistent)
		} else {
			firewall.RemoveHostsFromExceptions(apiAddrs, onlyForICMP, isPersistent)
		}
	}*/
}

// OnControlConnectionClosed - Perform reqired operations when protocol (controll channel with UI application) was closed
// (for example, we must disable firewall (if it not persistant))
// Must be called by protocol object
// Return parameters:
// - isServiceMustBeClosed: true informing that service have to be closed ("Stop IVPN Agent when application is not running" feature)
// - err: error
func (s *Service) OnControlConnectionClosed() (isServiceMustBeClosed bool, err error) {
	isServiceMustBeClosed = s._preferences.IsStopOnClientDisconnect
	// disable firewall if it not persistant
	if !s._preferences.IsFwPersistant {
		log.Info("Control connection was closed. Disabling firewall.")
		err = s.SetKillSwitchState(false)
	}
	return isServiceMustBeClosed, err
}

// ServersList - get VPN servers info
//func (s *Service) ServersList() (*types.ServersInfoResponse, error) {
//	return s._serversUpdater.GetServers()
//}

// APIRequest do custom request to API
func (s *Service) APIRequest(apiAlias string, ipTypeRequired protocolTypes.RequiredIPProtocol) (responseData []byte, err error) {

	if ipTypeRequired == protocolTypes.IPv6 {
		// IPV6-LOC-200 - IVPN Apps should request only IPv4 location information when connected  to the gateway, which doesn’t support IPv6
		vpn := s._vpn
		if vpn != nil && !vpn.IsPaused() && !vpn.IsIPv6InTunnel() {
			return nil, fmt.Errorf("no IPv6 support inside tunnel for current connection")
		}
	}

	return s._api.DoRequestByAlias(apiAlias, ipTypeRequired)
}

// GetDisabledFunctions returns info about functions which are disabled
// Some functionality can be not accessible
// It can happen, for example, if some external binaries not installed
// (e.g. obfsproxy or WireGuard on Linux)
func (s *Service) GetDisabledFunctions() (wgErr, ovpnErr, obfspErr, splitTunErr error) {
	if err := filerights.CheckFileAccessRightsExecutable(platform.OpenVpnBinaryPath()); err != nil {
		ovpnErr = fmt.Errorf("OpenVPN binary: %w", err)
	}

	if err := filerights.CheckFileAccessRightsExecutable(platform.ObfsproxyStartScript()); err != nil {
		obfspErr = fmt.Errorf("obfsproxy binary: %w", err)
	}

	if err := filerights.CheckFileAccessRightsExecutable(platform.WgBinaryPath()); err != nil {
		wgErr = fmt.Errorf("WireGuard binary: %w", err)
	} else {
		if err := filerights.CheckFileAccessRightsExecutable(platform.WgToolBinaryPath()); err != nil {
			wgErr = fmt.Errorf("WireGuard tools binary: %w", err)
		}
	}

	// returns non-nil error object if Split-Tunneling functionality not available
	splitTunErr = splittun.GetFuncNotAvailableError()

	if errors.Is(ovpnErr, os.ErrNotExist) {
		ovpnErr = fmt.Errorf("%w. Please install OpenVPN", ovpnErr)
	}
	if errors.Is(obfspErr, os.ErrNotExist) {
		obfspErr = fmt.Errorf("%w. Please install obfsproxy binary", obfspErr)
	}
	if errors.Is(wgErr, os.ErrNotExist) {
		wgErr = fmt.Errorf("%w. Please install WireGuard", wgErr)
	}

	return wgErr, ovpnErr, obfspErr, splitTunErr
}

// ConnectOpenVPN start OpenVPN connection
func (s *Service) ConnectOpenVPN(connectionParams openvpn.ConnectionParams, manualDNS dns.DnsSettings, firewallOn bool, firewallDuringConnection bool, stateChan chan<- vpn.StateInfo) error {

	createVpnObjfunc := func() (vpn.Process, error) {
		prefs := s.Preferences()

		// checking if functionality accessible
		_, ovpnErr, obfspErr, _ := s.GetDisabledFunctions()
		if ovpnErr != nil {
			return nil, ovpnErr
		}
		if prefs.IsObfsproxy && obfspErr != nil {
			return nil, obfspErr
		}

		connectionParams.SetCredentials(prefs.Session.OpenVPNUser, prefs.Session.OpenVPNPass)

		openVpnExtraParameters := ""
		// read user-defined extra parameters for OpenVPN configuration (if exists)
		extraParamsFile := platform.OpenvpnUserParamsFile()

		if helpers.FileExists(extraParamsFile) {
			if err := filerights.CheckFileAccessRightsConfig(extraParamsFile); err != nil {
				log.Info("NOTE! User-defined OpenVPN parameters are ignored! %w", err)
				os.Remove(extraParamsFile)
			} else {
				// read file line by line
				openVpnExtraParameters = func() string {
					var allParams strings.Builder

					file, err := os.Open(extraParamsFile)
					if err != nil {
						log.Error(err)
						return ""
					}
					defer file.Close()

					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						line := scanner.Text()
						line = strings.TrimSpace(line)
						if len(line) <= 0 {
							continue
						}
						if strings.HasPrefix(line, "#") {
							continue // comment
						}
						if strings.HasPrefix(line, ";") {
							continue // comment
						}
						allParams.WriteString(line + "\n")
					}

					if err := scanner.Err(); err != nil {
						log.Error("Failed to parse '%s': %s", extraParamsFile, err)
						return ""
					}
					return allParams.String()
				}()

				if len(openVpnExtraParameters) > 0 {
					log.Info(fmt.Sprintf("WARNING! User-defined OpenVPN parameters loaded from file '%s'!", extraParamsFile))
				}
			}
		}

		// creating OpenVPN object
		vpnObj, err := openvpn.NewOpenVpnObject(
			platform.OpenVpnBinaryPath(),
			platform.OpenvpnConfigFile(),
			platform.OpenvpnLogFile(),
			prefs.IsObfsproxy,
			openVpnExtraParameters,
			connectionParams)

		if err != nil {
			return nil, fmt.Errorf("failed to create new openVPN object: %w", err)
		}
		return vpnObj, nil
	}

	return s.keepConnection(createVpnObjfunc, manualDNS, firewallOn, firewallDuringConnection, stateChan)
}

// ConnectWireGuard start WireGuard connection
func (s *Service) ConnectWireGuard(connectionParams wireguard.ConnectionParams, manualDNS dns.DnsSettings, firewallOn bool, firewallDuringConnection bool, stateChan chan<- vpn.StateInfo) error {
	// stop active connection (if exists)
	if err := s.Disconnect(); err != nil {
		return fmt.Errorf("failed to connect. Unable to stop active connection: %w", err)
	}

	// checking if functionality accessible
	wgErr, _, _, _ := s.GetDisabledFunctions()
	if wgErr != nil {
		return wgErr
	}

	// Update WG keys, if necessary
	err := s.WireGuardGenerateKeys(true)
	if err != nil {
		// If new WG keys regeneration failed but we still have active keys - keep connecting
		// (this could happen, for example, when FW is enabled and we even not tried to make API request)
		// Return error only if the keys had to be regenerated more than 3 days ago.
		_, activePublicKey, _, _, lastUpdate, interval := s.WireGuardGetKeys()

		if len(activePublicKey) > 0 && lastUpdate.Add(interval).Add(time.Hour*24*3).After(time.Now()) {
			// continue connection
			log.Warning(fmt.Errorf("WG KEY generation failed (%w). But we keep connecting (will try to regenerate it next 3 days)", err))
		} else {
			return err
		}
	}

	createVpnObjfunc := func() (vpn.Process, error) {
		session := s.Preferences().Session

		if !session.IsWGCredentialsOk() {
			return nil, fmt.Errorf("WireGuard credentials are not defined (please, regenerate WG credentials or re-login)")
		}

		localip := net.ParseIP(session.WGLocalIP)
		if localip == nil {
			return nil, fmt.Errorf("error updating WG connection preferences (failed parsing local IP for WG connection)")
		}
		connectionParams.SetCredentials(session.WGPrivateKey, localip)

		vpnObj, err := wireguard.NewWireGuardObject(
			platform.WgBinaryPath(),
			platform.WgToolBinaryPath(),
			platform.WGConfigFilePath(),
			connectionParams)

		if err != nil {
			return nil, fmt.Errorf("failed to create new WireGuard object: %w", err)
		}
		return vpnObj, nil
	}

	return s.keepConnection(createVpnObjfunc, manualDNS, firewallOn, firewallDuringConnection, stateChan)
}

func (s *Service) keepConnection(createVpnObj func() (vpn.Process, error), manualDNS dns.DnsSettings, firewallOn bool, firewallDuringConnection bool, stateChan chan<- vpn.StateInfo) error {
	prefs := s.Preferences()
	if !prefs.Session.IsLoggedIn() {
		return srverrors.ErrorNotLoggedIn{}
	}

	s._manualDNS = manualDNS

	// Not necessary to keep connection until we are not connected
	// So just 'Connect' required for now
	s._requiredVpnState = Connect

	// no delay before first reconnection
	delayBeforeReconnect := 0 * time.Second

	stateChan <- vpn.NewStateInfo(vpn.CONNECTING, "Connecting")
	for {
		// create new VPN object
		vpnObj, err := createVpnObj()
		if err != nil {
			return fmt.Errorf("failed to create VPN object: %w", err)
		}

		lastConnectionTryTime := time.Now()

		// start connection
		connErr := s.connect(vpnObj, s._manualDNS, firewallOn, firewallDuringConnection, stateChan)
		if connErr != nil {
			log.Error(fmt.Sprintf("Connection error: %s", connErr))
			if s._requiredVpnState == Connect {
				// throw error only on first try to connect
				// if we were already connected (_requiredVpnState==KeepConnection) - ignore error and try to reconnect
				return connErr
			}
		}

		// retry, if reconnection requested
		if s._requiredVpnState == KeepConnection {
			// notifying clients about reconnection
			stateChan <- vpn.NewStateInfo(vpn.RECONNECTING, "Reconnecting due to disconnection")

			// no delay before reconnection (if last connection was long time ago)
			if time.Now().After(lastConnectionTryTime.Add(time.Second * 30)) {
				delayBeforeReconnect = 0
			}
			// no delay before reconnection if reconnection was requested by VPN object
			if connErr != nil {
				var reconnectReqErr *vpn.ReconnectionRequiredError
				if errors.As(connErr, &reconnectReqErr) {
					log.Info("VPN object requested re-connection")
					delayBeforeReconnect = 0
				}
			}

			if delayBeforeReconnect > 0 {
				log.Info(fmt.Sprintf("Reconnecting (pause %s)...", delayBeforeReconnect))
				// do delay before next reconnection
				pauseTill := time.Now().Add(delayBeforeReconnect)
				for time.Now().Before(pauseTill) && s._requiredVpnState != Disconnect {
					time.Sleep(time.Millisecond * 10)
				}
			} else {
				log.Info("Reconnecting...")
			}

			if s._requiredVpnState == KeepConnection {
				// consecutive reconnections has delay 5 seconds
				delayBeforeReconnect = time.Second * 5
				continue
			}
		}

		// stop loop
		break
	}

	return nil
}

// Connect connect vpn.
// Param 'firewallOn' - enable firewall before connection (if true - the parameter 'firewallDuringConnection' will be ignored).
// Param 'firewallDuringConnection' - enable firewall before connection and disable after disconnection (has effect only if Firewall not enabled before)
func (s *Service) connect(vpnProc vpn.Process, manualDNS dns.DnsSettings, firewallOn bool, firewallDuringConnection bool, stateChan chan<- vpn.StateInfo) error {
	var connectRoutinesWaiter sync.WaitGroup

	// stop active connection (if exists)
	if err := s.disconnect(); err != nil {
		return fmt.Errorf("failed to connect. Unable to stop active connection: %w", err)
	}

	// check session status each disconnection (asynchronously, in separate goroutine)
	defer func() { go s.RequestSessionStatus() }()

	s._connectMutex.Lock()
	defer s._connectMutex.Unlock()

	s._done = make(chan struct{}, 1)
	defer func() {
		// notify: connection stopped
		done := s._done
		s._done = nil
		if done != nil {
			done <- struct{}{}
			// Closing channel
			// Note: reading from empty and closed channel will not lead to deadlock (immediately returns zero value)
			close(done)
		}
	}()

	var err error

	log.Info("Connecting...")

	// checking default outbound IPs
	// It is necessary for Split-Tunnelling configuration
	var sInfo VpnSessionInfo
	sInfo.OutboundIPv4, err = netinfo.GetOutboundIP(false)
	if err != nil {
		log.Warning(fmt.Errorf("failed to detect outbound IPv4 address: %w", err))
	}
	sInfo.OutboundIPv6, err = netinfo.GetOutboundIP(true)
	if err != nil {
		log.Warning(fmt.Errorf("failed to detect outbound IPv6 address: %w", err))
	}
	s.SetVpnSessionInfo(sInfo)

	// save vpn object
	s._vpn = vpnProc

	internalStateChan := make(chan vpn.StateInfo, 1)
	stopChannel := make(chan bool, 1)

	fwInitState := false
	// finalize everything
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic on VPN connection: ", r)
			if err, ok := r.(error); ok {
				log.ErrorTrace(err)
			}
		}

		// Ensure that routing-change detector is stopped (we do not need it when VPN disconnected)
		s._netChangeDetector.Stop()

		// notify firewall that client is disconnected
		err := firewall.ClientDisconnected()
		if err != nil {
			log.Error("(stopping) error on notifying FW about disconnected client:", err)
		}

		// when we were requested to enable firewall for this connection
		// And initial FW state was disabled - we have to disable it back
		if firewallDuringConnection && !fwInitState {
			if err = s.SetKillSwitchState(false); err != nil {
				log.Error("(stopping) failed to disable firewall:", err)
			}
		}

		// notify routines to stop
		close(stopChannel)

		// resetting manual DNS (if it is necessary)
		err = vpnProc.ResetManualDNS()
		if err != nil {
			log.Error("(stopping) error resetting manual DNS: ", err)
		}

		connectRoutinesWaiter.Wait()

		// Forget VPN object
		s._vpn = nil

		// Notify Split-Tunneling module about disconnected VPN status
		s.splitTunnelling_ApplyConfig()

		log.Info("VPN process stopped")
	}()

	// Signaling when the default routing is NOT over the 'interfaceToProtect' anymore
	routingChangeChan := make(chan struct{}, 1)
	// Signaling when there were some routing changes but 'interfaceToProtect' is still is the default route
	routingUpdateChan := make(chan struct{}, 1)

	destinationHostIP := vpnProc.DestinationIP()

	// goroutine: process + forward VPN state change
	connectRoutinesWaiter.Add(1)
	go func() {
		log.Info("VPN state forwarder started")
		defer func() {
			log.Info("VPN state forwarder stopped")
			connectRoutinesWaiter.Done()
		}()

		var state vpn.StateInfo
		for isRuning := true; isRuning; {
			select {
			case state = <-internalStateChan:

				// store info about current time
				state.Time = time.Now().Unix()
				// store info about VPN connection type
				state.VpnType = vpnProc.Type()

				// forward state to 'stateChan'
				stateChan <- state

				log.Info(fmt.Sprintf("State: %v", state))

				// internally process VPN state change
				switch state.State {

				case vpn.RECONNECTING:
					// Disable routing-change detector when reconnecting
					s._netChangeDetector.Stop()

					// Add host IP to firewall exceptions
					// Some OS-specific implementations (e.g. macOS) can remove server host from firewall rules after connection established
					// We have to allow it's IP to be able to reconnect
					const onlyForICMP = false
					const isPersistent = false
					err := firewall.AddHostsToExceptions([]net.IP{destinationHostIP}, onlyForICMP, isPersistent)
					if err != nil {
						log.Error("Unable to add host to firewall exceptions:", err.Error())
					}

				case vpn.CONNECTED:
					// since we are connected - keep connection (reconnect if unexpected disconnection)
					if s._requiredVpnState == Connect {
						s._requiredVpnState = KeepConnection
					}

					// start routing change detection
					if netInterface, err := netinfo.InterfaceByIPAddr(state.ClientIP); err != nil {
						log.Error(fmt.Sprintf("Unable to initialize routing change detection. Failed to get interface '%s'", state.ClientIP.String()))
					} else {

						log.Info("Starting route change detection")
						s._netChangeDetector.Start(routingChangeChan, routingUpdateChan, netInterface)
					}

					// Inform firewall about client local IP
					firewall.ClientConnected(
						state.ClientIP, state.ClientIPv6,
						state.ClientPort,
						state.ServerIP, state.ServerPort,
						state.IsTCP)

					// save ClientIP/ClientIPv6 into vpn-session-info
					sInfo := s.GetVpnSessionInfo()
					sInfo.VpnLocalIPv4 = state.ClientIP
					sInfo.VpnLocalIPv6 = state.ClientIPv6
					s.SetVpnSessionInfo(sInfo)

					// Notify Split-Tunneling module about connected VPN status
					s.splitTunnelling_ApplyConfig()
				default:
				}

			case <-stopChannel: // triggered when the stopChannel is closed
				isRuning = false
			}
		}
	}()

	// receiving routing change notifications
	connectRoutinesWaiter.Add(1)
	go func() {
		log.Info("Route change receiver started")
		defer func() {
			log.Info("Route change receiver stopped")
			connectRoutinesWaiter.Done()
		}()

		for isRuning := true; isRuning; {
			select {
			case <-routingChangeChan: // routing changed (the default routing is NOT over the 'interfaceToProtect' anymore)
				if s._vpn.IsPaused() {
					log.Info("Route change ignored due to Paused state.")
				} else {
					// Disconnect (client will request then reconnection, because of unexpected disconnection)
					// reconnect in separate routine (do not block current thread)
					go func() {
						defer func() {
							if r := recover(); r != nil {
								log.Error("PANIC: ", r)
							}
						}()

						log.Info("Route change detected. Reconnecting...")
						s.reconnect()
					}()

					isRuning = false
				}
			case <-routingUpdateChan: // there were some routing changes but 'interfaceToProtect' is still is the default route
				s._vpn.OnRoutingChanged()
			case <-stopChannel: // triggered when the stopChannel is closed
				isRuning = false
			}
		}
	}()

	log.Info("Initializing...")
	if err := vpnProc.Init(); err != nil {
		return fmt.Errorf("failed to initialize VPN object: %w", err)
	}

	log.Info("Initializing firewall")
	// firewallOn - enable firewall before connection (if true - the parameter 'firewallDuringConnection' will be ignored)
	// firewallDuringConnection - enable firewall before connection and disable after disconnection (has effect only if Firewall not enabled before)
	if firewallOn {
		fw, err := firewall.GetEnabled()
		if err != nil {
			log.Error("Failed to check firewall state:", err.Error())
			return err
		}
		if !fw {
			if err := s.SetKillSwitchState(true); err != nil {
				log.Error("Failed to enable firewall:", err.Error())
				return err
			}
		}
	} else if firewallDuringConnection {
		// in case to enable FW for this connection parameter:
		// - check initial FW state
		// - if it disabled - enable it (will be disabled on disconnect)
		fw, err := firewall.GetEnabled()
		if err != nil {
			log.Error("Failed to check firewall state:", err.Error())
			return err
		}
		fwInitState = fw
		if !fwInitState {
			if err := s.SetKillSwitchState(true); err != nil {
				log.Error("Failed to enable firewall:", err.Error())
				return err
			}
		}
	}

	// Add host IP to firewall exceptions
	const onlyForICMP = false
	const isPersistent = false
	err = firewall.AddHostsToExceptions([]net.IP{destinationHostIP}, onlyForICMP, isPersistent)
	if err != nil {
		log.Error("Failed to start. Unable to add hosts to firewall exceptions:", err.Error())
		return err
	}

	log.Info("Initializing DNS")
	// set manual DNS
	if manualDNS.IsEmpty() {
		err = s.ResetManualDNS()
	} else {
		err = s.SetManualDNS(manualDNS)
	}
	if err != nil {
		err = fmt.Errorf("failed to set DNS: %w", err)
		log.Error(err.Error())
		return err
	}

	log.Info("Starting VPN process")
	// connect: start VPN process and wait until it finishes
	err = vpnProc.Connect(internalStateChan)
	if err != nil {
		err = fmt.Errorf("connection error: %w", err)
		log.Error(err.Error())
		return err
	}

	return nil
}

func (s *Service) reconnect() {
	// Just call disconnect
	// The reconnection will be performed automatically in method 'keepConnection(...)'
	// (according to s._requiredVpnState value == KeepConnection)
	s.disconnect()
}

// Disconnect disconnect vpn
func (s *Service) Disconnect() error {
	s._requiredVpnState = Disconnect
	if err := s.Resume(); err != nil {
		log.Error("Resume failed:", err)
	}
	return s.disconnect()
}

func (s *Service) disconnect() error {
	vpn := s._vpn
	if vpn == nil {
		return nil
	}

	done := s._done
	if s._requiredVpnState == KeepConnection {
		log.Info("Disconnecting (going to reconnect)...")
	} else {
		log.Info("Disconnecting...")
	}

	// stop detections for routing changes
	s._netChangeDetector.Stop()

	// stop VPN
	if err := vpn.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect VPN: %w", err)
	}

	// wait for stop
	if done != nil {
		<-done
	}

	return nil
}

// Connected returns 'true' if VPN connected
func (s *Service) Connected() bool {
	return s._vpn != nil
}

// ConnectedType returns connected VPN type (only if VPN connected!)
func (s *Service) ConnectedType() (isConnected bool, connectedVpnType vpn.Type) {
	vpnObj := s._vpn
	if vpnObj == nil {
		return false, 0
	}
	return true, vpnObj.Type()
}

// FirewallEnabled returns firewall state (enabled\disabled)
// (in use, for example, by WireGuard keys manager, to know is it have sense to make API requests.)
func (s *Service) FirewallEnabled() (bool, error) {
	return firewall.GetEnabled()
}

// Pause pause vpn connection
func (s *Service) Pause() error {
	vpn := s._vpn
	if vpn == nil {
		return nil
	}

	log.Info("Pausing...")
	firewall.ClientPaused()
	return vpn.Pause()
}

// Resume resume vpn connection
func (s *Service) Resume() error {
	vpn := s._vpn
	if vpn == nil {
		return nil
	}
	if !vpn.IsPaused() {
		return nil
	}

	log.Info("Resuming...")
	firewall.ClientResumed()
	return vpn.Resume()
}

// IsPaused returns 'true' if current vpn connection is in paused state
func (s *Service) IsPaused() bool {
	vpn := s._vpn
	if vpn == nil {
		return false
	}
	return vpn.IsPaused()
}

// SetManualDNS set dns
func (s *Service) SetManualDNS(dnsCfg dns.DnsSettings) error {
	vpn := s._vpn
	if vpn == nil {
		return nil
	}

	s._manualDNS = dnsCfg

	if dnsCfg.Encryption == dns.EncryptionNone {
		// for DoH/DoT - no sense to allow DNS port (53)
		if err := firewall.SetManualDNS(s._manualDNS.Ip()); err != nil {
			return fmt.Errorf("failed to set manual DNS: %w", err)
		}
	} else {
		firewall.SetManualDNS(nil)
	}

	return vpn.SetManualDNS(s._manualDNS)
}

// ResetManualDNS set dns to default
func (s *Service) ResetManualDNS() error {
	vpn := s._vpn
	if vpn == nil {
		return nil
	}

	s._manualDNS = dns.DnsSettings{}
	if err := firewall.SetManualDNS(nil); err != nil {
		return fmt.Errorf("failed to reset manual DNS: %w", err)
	}

	return vpn.ResetManualDNS()
}

// if 'currentLocation' defined - the output hosts list will be sorted by distance to current location
func (s *Service) getHostsToPing(currentLocation *types.GeoLookupResponse) ([]net.IP, error) {
	// get servers info
	servers, err := s._serversUpdater.GetServers()
	if err != nil {
		return nil, fmt.Errorf("unable to get servers list: %w", err)
	}

	type hostInfo struct {
		Latitude  float32
		Longitude float32
		host      net.IP
	}

	hosts := make([]hostInfo, 0, uint16(len(servers.Servers)))

	// OpenVPN servers
	for _, s := range servers.Servers {
		if len(s.Hosts) <= 0 {
			continue
		}

		ip := net.ParseIP(s.Hosts[0].IP)
		if ip != nil {
			hosts = append(hosts, hostInfo{host: ip})
		}
	}
	ret := make([]net.IP, 0, len(hosts))
	for _, h := range hosts {
		ret = append(ret, h.host)
	}
	return ret, nil
}

//////////////////////////////////////////////////////////
// KillSwitch
//////////////////////////////////////////////////////////

// SetKillSwitchState enable\disable killswitch
func (s *Service) SetKillSwitchState(isEnabled bool) error {

	if !isEnabled && s._preferences.IsFwPersistant {
		return fmt.Errorf("unable to disable Firewall in 'Persistent' state. Please, disable 'Always-on firewall' first")
	}

	err := firewall.SetEnabled(isEnabled)
	if err == nil {
		s._evtReceiver.OnKillSwitchStateChanged()
	}
	return err
}

// KillSwitchState returns killswitch state
func (s *Service) KillSwitchState() (isEnabled, isPersistant, isAllowLAN, isAllowLanMulticast, isAllowApiServers bool, err error) {
	prefs := s._preferences
	enabled, err := firewall.GetEnabled()
	return enabled, prefs.IsFwPersistant, prefs.IsFwAllowLAN, prefs.IsFwAllowLANMulticast, prefs.IsFwAllowApiServers, err
}

// SetKillSwitchIsPersistent change kill-switch value
func (s *Service) SetKillSwitchIsPersistent(isPersistant bool) error {
	prefs := s._preferences
	prefs.IsFwPersistant = isPersistant
	s.setPreferences(prefs)

	err := firewall.SetPersistant(isPersistant)
	if err == nil {
		s._evtReceiver.OnKillSwitchStateChanged()
	}
	return err
}

// SetKillSwitchAllowLAN change kill-switch value
func (s *Service) SetKillSwitchAllowLAN(isAllowLan bool) error {
	return s.setKillSwitchAllowLAN(isAllowLan, s._preferences.IsFwAllowLANMulticast)
}

// SetKillSwitchAllowLANMulticast change kill-switch value
func (s *Service) SetKillSwitchAllowLANMulticast(isAllowLanMulticast bool) error {
	return s.setKillSwitchAllowLAN(s._preferences.IsFwAllowLAN, isAllowLanMulticast)
}

func (s *Service) setKillSwitchAllowLAN(isAllowLan bool, isAllowLanMulticast bool) error {
	prefs := s._preferences
	prefs.IsFwAllowLAN = isAllowLan
	prefs.IsFwAllowLANMulticast = isAllowLanMulticast
	s.setPreferences(prefs)

	err := firewall.AllowLAN(prefs.IsFwAllowLAN, prefs.IsFwAllowLANMulticast)
	if err == nil {
		s._evtReceiver.OnKillSwitchStateChanged()
	}
	return err
}

func (s *Service) SetKillSwitchAllowAPIServers(isAllowAPIServers bool) error {
	if !isAllowAPIServers {
		// Do not allow to disable access to IVPN API server if user logged-out
		// Otherwise, we will not have possibility to login
		session := s.Preferences().Session
		if !session.IsLoggedIn() {
			return srverrors.ErrorNotLoggedIn{}
		}
	}

	prefs := s._preferences
	prefs.IsFwAllowApiServers = isAllowAPIServers
	s.setPreferences(prefs)
	s._evtReceiver.OnKillSwitchStateChanged()
	s.updateAPIAddrInFWExceptions()
	return nil
}

//////////////////////////////////////////////////////////
// PREFERENCES
//////////////////////////////////////////////////////////

// SetPreference set preference value
func (s *Service) SetPreference(key string, val string) error {
	prefs := s._preferences

	switch key {
	case "enable_logging":
		if val, err := strconv.ParseBool(val); err == nil {
			prefs.IsLogging = val
			logger.Enable(val)
		}
	case "is_stop_server_on_client_disconnect":
		if val, err := strconv.ParseBool(val); err == nil {
			prefs.IsStopOnClientDisconnect = val
		}
	case "enable_obfsproxy":
		if val, err := strconv.ParseBool(val); err == nil {
			prefs.IsObfsproxy = val
		}
	case "firewall_is_persistent":
		log.Debug("Skipping 'firewall_is_persistent' value. IVPNKillSwitchSetIsPersistentRequest should be used")
	default:
		log.Warning(fmt.Sprintf("Preference key '%s' not supported", key))
	}

	s.setPreferences(prefs)
	log.Info(fmt.Sprintf("preferences %s='%s'", key, val))

	return nil
}

// Preferences returns preferences
func (s *Service) Preferences() preferences.Preferences {
	return s._preferences
}

func (s *Service) ResetPreferences() error {
	s._preferences = *preferences.Create()

	// erase ST config
	s.SplitTunnelling_SetConfig(false, true)
	return nil
}

//////////////////////////////////////////////////////////
// SPLIT TUNNEL
//////////////////////////////////////////////////////////

func (s *Service) GetInstalledApps(extraArgsJSON string) ([]oshelpers.AppInfo, error) {
	return oshelpers.GetInstalledApps(extraArgsJSON)
}

func (s *Service) GetBinaryIcon(binaryPath string) (string, error) {
	return oshelpers.GetBinaryIconBase64(binaryPath)
}

func (s *Service) SplitTunnelling_GetStatus() (protocolTypes.SplitTunnelStatus, error) {
	var prefs = s.Preferences()
	runningProcesses, err := splittun.GetRunningApps()
	if err != nil {
		runningProcesses = []splittun.RunningApp{}
	}

	ret := protocolTypes.SplitTunnelStatus{
		IsFunctionalityNotAvailable: splittun.GetFuncNotAvailableError() != nil,
		IsEnabled:                   prefs.IsSplitTunnel,
		IsCanGetAppIconForBinary:    oshelpers.IsCanGetAppIconForBinary(),
		SplitTunnelApps:             prefs.SplitTunnelApps,
		RunningApps:                 runningProcesses}

	return ret, nil
}

func (s *Service) SplitTunnelling_SetConfig(isEnabled bool, reset bool) error {
	if reset || splittun.GetFuncNotAvailableError() != nil {
		return s.splitTunnelling_Reset()
	}

	prefs := s._preferences
	prefs.IsSplitTunnel = isEnabled
	s.setPreferences(prefs)

	return s.splitTunnelling_ApplyConfig()
}
func (s *Service) splitTunnelling_Reset() error {
	prefs := s._preferences
	prefs.IsSplitTunnel = false
	prefs.SplitTunnelApps = make([]string, 0)
	s.setPreferences(prefs)

	splittun.Reset()

	return s.splitTunnelling_ApplyConfig()
}
func (s *Service) splitTunnelling_ApplyConfig() error {
	// notify changed ST configuration status (even if functionality not available)
	defer s._evtReceiver.OnSplitTunnelStatusChanged()

	if splittun.GetFuncNotAvailableError() != nil {
		// Split-Tunneling not accessable (not able to connect to a driver or not implemented for current platform)
		return nil
	}

	prefs := s.Preferences()
	sInf := s.GetVpnSessionInfo()

	addressesCfg := splittun.ConfigAddresses{
		IPv4Tunnel: sInf.VpnLocalIPv4,
		IPv4Public: sInf.OutboundIPv4,
		IPv6Tunnel: sInf.VpnLocalIPv6,
		IPv6Public: sInf.OutboundIPv6}

	return splittun.ApplyConfig(prefs.IsSplitTunnel, s.Connected(), addressesCfg, prefs.SplitTunnelApps)
}

func (s *Service) SplitTunnelling_AddApp(exec string) (cmdToExecute string, isAlreadyRunning bool, err error) {
	if !s._preferences.IsSplitTunnel {
		return "", false, fmt.Errorf("unable to run application in Split Tunneling environment: Split Tunneling is disabled")
	}
	// apply ST configuration after function ends
	defer s.splitTunnelling_ApplyConfig()
	return s.implSplitTunnelling_AddApp(exec)
}

func (s *Service) SplitTunnelling_RemoveApp(pid int, exec string) (err error) {
	// apply ST configuration after function ends
	defer s.splitTunnelling_ApplyConfig()
	return s.implSplitTunnelling_RemoveApp(pid, exec)
}

// SplitTunnelling_AddedPidInfo Inform the daemon about started process in ST environment
// Parameters:
// pid 			- process PID
// exec 		- Command executed in ST environment (e.g. binary + arguments)
// 				  (identical to SplitTunnelAddApp.Exec and SplitTunnelAddAppCmdResp.Exec)
// cmdToExecute - Shell command used to perform this operation
func (s *Service) SplitTunnelling_AddedPidInfo(pid int, exec string, cmdToExecute string) error {
	// notify changed ST configuration status
	defer s._evtReceiver.OnSplitTunnelStatusChanged()
	return s.implSplitTunnelling_AddedPidInfo(pid, exec, cmdToExecute)
}

//////////////////////////////////////////////////////////
// SESSIONS
//////////////////////////////////////////////////////////

func (s *Service) setCredentials(accountID, session, vpnUser, vpnPass, wgPublicKey, wgPrivateKey, wgLocalIP string, wgKeyGenerated int64) error {
	// save session info
	s._preferences.SetSession(accountID,
		session,
		vpnUser,
		vpnPass,
		wgPublicKey,
		wgPrivateKey,
		wgLocalIP)

	// manually set info about WG keys timestamp
	if wgKeyGenerated > 0 {
		s._preferences.Session.WGKeyGenerated = time.Unix(wgKeyGenerated, 0)
		s._preferences.SavePreferences()
	}

	// notify clients about session update
	s._evtReceiver.OnServiceSessionChanged()

	// success
	s.startSessionChecker()
	return nil
}

// SessionNew creates new session
func (s *Service) SessionNew(username string, password string) (
	apiCode int,
	apiErrorMsg string,
	accountInfo preferences.AccountStatus,
	rawResponse string,
	err error) {

	// Temporary allow API server access (If Firewall is enabled)
	// Otherwise, there will not be any possibility to Login (because all connectivity is blocked)
	fwIsEnabled, _, _, _, fwIsAllowApiServers, _ := s.KillSwitchState()
	if fwIsEnabled && !fwIsAllowApiServers {
		s.SetKillSwitchAllowAPIServers(true)
	}
	defer func() {
		if fwIsEnabled && !fwIsAllowApiServers {
			// restore state for 'AllowAPIServers' configuration (previously, was enabled)
			s.SetKillSwitchAllowAPIServers(false)
		}
	}()

	// delete current session (if exists)
	isCanDeleteSessionLocally := true
	if err := s.SessionDelete(isCanDeleteSessionLocally); err != nil {
		log.Error("Creating new session -> Failed to delete active session: ", err)
	}

	// generate new keys for WireGuard
	publicKey, privateKey, err := wireguard.GenerateKeys(platform.WgToolBinaryPath())
	if err != nil {
		log.Warning("Failed to generate wireguard keys for new session: %s", err)
	}

	log.Info("Logging in...")
	defer func() {
		if err != nil {
			log.Info("Logging in - FAILED: ", err)
		} else {
			log.Info("Logging in - SUCCESS")

		}
	}()
	successResp, errorLimitResp, apiErr, rawRespStr, err := s._api.SessionNew(username, password)
	rawResponse = rawRespStr

	apiCode = 0
	if apiErr != nil {
		apiCode = apiErr.Status
	}

	if err != nil {
		// if SessionsLimit response
		if errorLimitResp != nil {
			accountInfo = s.createAccountStatus(errorLimitResp.SessionLimitData)
			return apiCode, apiErr.Message, accountInfo, rawResponse, err
		}

		// in case of other API error
		if apiErr != nil {
			return apiCode, apiErr.Message, accountInfo, rawResponse, err
		}
		// not API error
		return apiCode, "", accountInfo, rawResponse, err
	}

	if successResp == nil {
		return apiCode, "", accountInfo, rawResponse, fmt.Errorf("unexpected error when creating a new session")
	}

	// get account status info
	accountInfo = s.createAccountStatus(successResp.ServiceStatus)

	s.setCredentials(username,
		successResp.Token,
		successResp.VpnUsername,
		successResp.VpnPassword,
		publicKey,
		privateKey, "", 0)

	return apiCode, "", accountInfo, rawResponse, nil
}

func (s *Service) ServersList() (
	response *types.ServersInfoResponse,
	rawResponse string,
	err error) {
	session := s.Preferences().Session
	if !session.IsLoggedIn() {
		return nil, "", srverrors.ErrorNotLoggedIn{}
	}
	resp, list, err := s._api.ServersList(session.Session)
	if err != nil {
		return nil, "", err
	}
	err = writeServersToCache(resp)
	if err != nil {
		fmt.Println("unable to write servers data to file.")
	}
	return resp, list, err
}

// SessionDelete removes session info
func (s *Service) SessionDelete(isCanDeleteSessionLocally bool) error {
	sessionNeedToDeleteOnBackend := true
	return s.logOut(sessionNeedToDeleteOnBackend, isCanDeleteSessionLocally)
}

// logOut performs log out from current session
// 1) if 'sessionNeedToDeleteOnBackend' == false: the app not trying to make API request
//	  the session info just erasing locally
//    (this is useful for the situations when we already know that session is not available on backend anymore)
// 2) if 'sessionNeedToDeleteOnBackend' == true (and 'isCanDeleteSessionLocally' == false): app is trying to make API request to logout correctly
//	  in case if API request failed the function returns error (session keeps not logged out)
// 3) if 'isCanDeleteSessionLocally' == true (and 'sessionNeedToDeleteOnBackend' == true): app is trying to make API request to logout correctly
//	  in case if API request failed we just erasing session info locally (no errors returned)

func (s *Service) logOut(sessionNeedToDeleteOnBackend bool, isCanDeleteSessionLocally bool) error {
	// Disconnect (if connected)
	s.Disconnect()

	// stop session checker (use goroutine to avoid deadlocks)
	go s.stopSessionChecker()

	// stop WG keys rotation
	s._wgKeysMgr.StopKeysRotation()

	if sessionNeedToDeleteOnBackend {

		// Temporary allow API server access (If Firewall is enabled)
		// Otherwise, there will not be any possibility to Login (because all connectivity is blocked)
		fwIsEnabled, _, _, _, fwIsAllowApiServers, _ := s.KillSwitchState()
		if fwIsEnabled && !fwIsAllowApiServers {
			s.SetKillSwitchAllowAPIServers(true)
		}
		defer func() {
			if fwIsEnabled && !fwIsAllowApiServers {
				// restore state for 'AllowAPIServers' configuration (previously, was enabled)
				s.SetKillSwitchAllowAPIServers(false)
			}
		}()

		session := s.Preferences().Session
		if session.IsLoggedIn() {
			log.Info("Logging out")
			err := s._api.SessionDelete(session.Session)
			if err != nil {
				log.Info("Logging out error:", err)
				if !isCanDeleteSessionLocally {
					return err // do not allow to logout if failed to delete session on backend
				}
			} else {
				log.Info("Logging out: done")
			}
		}
	}

	s._preferences.SetSession("", "", "", "", "", "", "")
	log.Info("Logged out locally")

	// notify clients about session update
	s._evtReceiver.OnServiceSessionChanged()
	return nil
}

func (s *Service) OnSessionNotFound() {
	// Logging out now
	log.Info("Session not found. Logging out.")
	needToDeleteOnBackend := false
	canLogoutOnlyLocally := true
	s.logOut(needToDeleteOnBackend, canLogoutOnlyLocally)
}

// RequestSessionStatus receives session status
func (s *Service) RequestSessionStatus() (
	apiCode int,
	apiErrorMsg string,
	sessionToken string,
	accountInfo preferences.AccountStatus,
	err error) {

	session := s.Preferences().Session
	if !session.IsLoggedIn() {
		return apiCode, "", "", accountInfo, srverrors.ErrorNotLoggedIn{}
	}

	log.Info("Requesting session status...")
	stat, apiErr, err := s._api.SessionStatus(session.Session)
	log.Info(fmt.Sprintf("Session status request: done code: %d", apiCode))

	currSession := s.Preferences().Session
	if currSession.Session != session.Session {
		// It could happen that logout\login was performed during the session check
		// Ignoring result if there is already a new session
		log.Info("Ignoring requested session status result. Local session already changed.")
		return apiCode, "", "", accountInfo, srverrors.ErrorNotLoggedIn{}
	}

	apiCode = 0
	if apiErr != nil {
		apiCode = apiErr.Status
		// Session not found - can happens when user forced to logout from another device
		if apiCode == types.SessionNotFound {
			s.OnSessionNotFound()
		}

		// notify clients that account not active
		if apiCode == types.AccountNotActive {
			// notify about account status
			s._evtReceiver.OnAccountStatus(session.Session, accountInfo)
			return apiCode, apiErr.Message, session.Session, preferences.AccountStatus{Active: false}, err
		}
	}

	if err != nil {
		// in case of other API error
		if apiErr != nil {
			return apiCode, apiErr.Message, "", accountInfo, err
		}

		// not API error
		return apiCode, "", "", accountInfo, err
	}

	if stat == nil {
		return apiCode, "", "", accountInfo, fmt.Errorf("unexpected error when creating requesting session status")
	}

	// get account status info
	accountInfo = s.createAccountStatus(*stat)
	// notify about account status
	s._evtReceiver.OnAccountStatus(session.Session, accountInfo)

	// success
	return apiCode, "", session.Session, accountInfo, nil
}

func (s *Service) createAccountStatus(apiResp types.ServiceStatusAPIResp) preferences.AccountStatus {
	return preferences.AccountStatus{
		Active:      apiResp.Active,
		ActiveUntil: apiResp.ActiveUntil,
		CurrentPlan: apiResp.CurrentPlan,
		IsFreeTrial: apiResp.IsFreeTrial}
}

func (s *Service) startSessionChecker() {
	// ensure that session checker is not running
	s.stopSessionChecker()

	session := s.Preferences().Session
	if !session.IsLoggedIn() {
		return
	}

	s._sessionCheckerStopChn = make(chan struct{})
	go func() {
		log.Info("Session checker started")
		defer log.Info("Session checker stopped")

		stopChn := s._sessionCheckerStopChn
		for {
			// wait for timeout or stop request
			select {
			case <-stopChn:
				return
			case <-time.After(SessionCheckInterval):
			}

			// check status
			s.RequestSessionStatus()

			// if not logged-in - no sense to check status anymore
			session := s.Preferences().Session
			if !session.IsLoggedIn() {
				return
			}
		}
	}()
}

func (s *Service) stopSessionChecker() {
	stopChan := s._sessionCheckerStopChn
	s._sessionCheckerStopChn = nil
	if stopChan != nil {
		stopChan <- struct{}{}
	}
}

//////////////////////////////////////////////////////////
// WireGuard keys
//////////////////////////////////////////////////////////

// WireGuardSaveNewKeys saves WG keys
func (s *Service) WireGuardSaveNewKeys(wgPublicKey string, wgPrivateKey string, wgLocalIP string) {
	s._preferences.UpdateWgCredentials(wgPublicKey, wgPrivateKey, wgLocalIP)

	// notify clients about session (wg keys) update
	s._evtReceiver.OnServiceSessionChanged()

	go func() {
		// reconnect in separate routine (do not block current thread)
		vpnObj := s._vpn
		if vpnObj == nil {
			return
		}
		if vpnObj.Type() != vpn.WireGuard {
			return
		}
		if !s.Connected() {
			return
		}
		log.Info("Reconnecting WireGuard connection with new credentials...")
		s.reconnect()
	}()
}

// WireGuardSetKeysRotationInterval change WG key rotation interval
func (s *Service) WireGuardSetKeysRotationInterval(interval int64) {
	s._preferences.Session.WGKeysRegenInerval = time.Second * time.Duration(interval)
	s._preferences.SavePreferences()

	// restart WG keys rotation
	if err := s._wgKeysMgr.StartKeysRotation(); err != nil {
		log.Error(err)
	}

	// notify clients about session (wg keys) update
	s._evtReceiver.OnServiceSessionChanged()
}

// WireGuardGetKeys get WG keys
func (s *Service) WireGuardGetKeys() (session, wgPublicKey, wgPrivateKey, wgLocalIP string, generatedTime time.Time, updateInterval time.Duration) {
	p := s._preferences

	return p.Session.Session,
		p.Session.WGPublicKey,
		p.Session.WGPrivateKey,
		p.Session.WGLocalIP,
		p.Session.WGKeyGenerated,
		p.Session.WGKeysRegenInerval
}

// WireGuardGenerateKeys - generate new wireguard keys
func (s *Service) WireGuardGenerateKeys(updateIfNecessary bool) error {
	if !s._preferences.Session.IsLoggedIn() {
		return srverrors.ErrorNotLoggedIn{}
	}

	// Update WG keys, if necessary
	var err error
	if updateIfNecessary {
		_, err = s._wgKeysMgr.UpdateKeysIfNecessary()
	} else {
		err = s._wgKeysMgr.GenerateKeys()
	}
	if err != nil {
		return fmt.Errorf("failed to regenerate WireGuard keys: %w", err)
	}

	return nil
}

//////////////////////////////////////////////////////////
// Internal methods
//////////////////////////////////////////////////////////

func (s *Service) setPreferences(p preferences.Preferences) {
	if !reflect.DeepEqual(s._preferences, p) {
		//if s._preferences != p {
		s._preferences = p
		s._preferences.SavePreferences()
	}
}


package firewall

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/netinfo"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/firewall/winlib"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/platform"
)

var (
	providerKey = syscall.GUID{Data1: 0xfed0afd4, Data2: 0x98d4, Data3: 0x4233, Data4: [8]byte{0xa4, 0xf3, 0x8b, 0x7c, 0x02, 0x44, 0x50, 0x01}}
	sublayerKey = syscall.GUID{Data1: 0xfed0afd4, Data2: 0x98d4, Data3: 0x4233, Data4: [8]byte{0xa4, 0xf3, 0x8b, 0x7c, 0x02, 0x44, 0x50, 0x02}}

	v4Layers = []syscall.GUID{winlib.FwpmLayerAleAuthConnectV4, winlib.FwpmLayerAleAuthRecvAcceptV4}
	v6Layers = []syscall.GUID{winlib.FwpmLayerAleAuthConnectV6, winlib.FwpmLayerAleAuthRecvAcceptV6}

	manager                winlib.Manager
	clientLocalIPFilterIDs []uint64
	customDNS              net.IP

	isPersistant        bool
	isAllowLAN          bool
	isAllowLANMulticast bool
)

const (
	providerDName = "VPNwitch"
	sublayerDName = "VPNwitch Sub-Layer"
	filterDName   = "VPNwitch filter"
)

// implInitialize doing initialization stuff (called on application start)
func implInitialize() error {
	if err := winlib.Initialize(platform.WindowsWFPDllPath()); err != nil {
		return err
	}

	pInfo, err := manager.GetProviderInfo(providerKey)
	if err != nil {
		return err
	}

	// save initial persistant state into package-variable
	isPersistant = pInfo.IsPersistent

	return nil
}

func implGetEnabled() (bool, error) {
	pInfo, err := manager.GetProviderInfo(providerKey)
	if err != nil {
		return false, fmt.Errorf("failed to get provider info: %w", err)
	}
	return pInfo.IsInstalled, nil
}

func implSetEnabled(isEnabled bool) (retErr error) {
	// start transaction
	if err := manager.TransactionStart(); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	// do not forget to stop transaction
	defer func() {
		if r := recover(); r == nil {
			manager.TransactionCommit() // commit WFP transaction
		} else {
			manager.TransactionAbort() // abort WFPtransaction

			log.Error("PANIC (recovered): ", r)
			if e, ok := r.(error); ok {
				retErr = e
			} else {
				retErr = errors.New(fmt.Sprint(r))
			}
		}
	}()

	if isEnabled {
		return doEnable()
	}
	return doDisable()
}

func implSetPersistant(persistant bool) (retErr error) {
	// save persistent state
	isPersistant = persistant

	pinfo, err := manager.GetProviderInfo(providerKey)
	if err != nil {
		return fmt.Errorf("failed to get provider info: %w", err)
	}

	if pinfo.IsInstalled == true {
		if pinfo.IsPersistent == isPersistant {
			log.Info(fmt.Sprintf("Already enabled (persistent=%t).", isPersistant))
			return nil
		}

		log.Info(fmt.Sprintf("Re-enabling with persistent flag = %t", isPersistant))
		return reEnable()
	}

	return doEnable()
}

// ClientConnected - allow communication for local vpn/client IP address
func implClientConnected(clientLocalIPAddress net.IP, clientLocalIPv6Address net.IP, clientPort int, serverIP net.IP, serverPort int, isTCP bool) (retErr error) {
	// start / commit transaction
	if err := manager.TransactionStart(); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if retErr == nil {
			manager.TransactionCommit()
		} else {
			// abort transaction if there was an error
			manager.TransactionAbort()
		}
	}()

	err := doRemoveClientIPFilters()
	if err != nil {
		log.Error("Failed to remove previously defined client IP filters: ", err)
	}
	return doAddClientIPFilters(clientLocalIPAddress, clientLocalIPv6Address)
}

// ClientDisconnected - Disable communication for local vpn/client IP address
func implClientDisconnected() (retErr error) {
	// start / commit transaction
	if err := manager.TransactionStart(); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if retErr == nil {
			manager.TransactionCommit()
		} else {
			// abort transaction if there was an error
			manager.TransactionAbort()
		}
	}()

	return doRemoveClientIPFilters()
}

func implAddHostsToExceptions(IPs []net.IP, onlyForICMP bool, isPersistent bool) error {
	// nothing to do for windows implementation
	return nil
}

func implRemoveHostsFromExceptions(IPs []net.IP, onlyForICMP bool, isPersistent bool) error {
	// nothing to do for windows implementation
	return nil
}

// AllowLAN - allow/forbid LAN communication
func implAllowLAN(allowLan bool, allowLanMulticast bool) error {

	if isAllowLAN == allowLan && isAllowLANMulticast == allowLanMulticast {
		return nil
	}

	isAllowLAN = allowLan
	isAllowLANMulticast = allowLanMulticast

	enabled, err := implGetEnabled()
	if err != nil {
		return fmt.Errorf("failed to get info if firewall is on: %w", err)
	}
	if enabled == false {
		return nil
	}

	return reEnable()
}

// SetManualDNS - configure firewall to allow DNS which is out of VPN tunnel
// Applicable to Windows implementation (to allow custom DNS from local network)
func implSetManualDNS(addr net.IP) error {
	if addr.Equal(customDNS) {
		return nil
	}

	customDNS = addr

	enabled, err := implGetEnabled()
	if err != nil {
		return fmt.Errorf("failed to get info if firewall is on: %w", err)
	}
	if enabled == false {
		return nil
	}

	return reEnable()
}

func reEnable() (retErr error) {
	// start / commit transaction
	if err := manager.TransactionStart(); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if retErr == nil {
			manager.TransactionCommit()
		} else {
			// abort transaction if there was an error
			manager.TransactionAbort()
		}
	}()

	err := doDisable()
	if err != nil {
		return fmt.Errorf("failed to disable firewall: %w", err)
	}

	err = doEnable()
	if err != nil {
		return fmt.Errorf("failed to enable firewall: %w", err)
	}

	return doAddClientIPFilters(connectedClientInterfaceIP, connectedClientInterfaceIPv6)
}

func doEnable() (retErr error) {
	enabled, err := implGetEnabled()
	if err != nil {
		return fmt.Errorf("failed to get info if firewall is on: %w", err)
	}
	if enabled == true {
		return nil
	}

	addressesV6, err := netinfo.GetAllLocalV6Addresses()
	if err != nil {
		return fmt.Errorf("failed to get all local IPv6 addresses: %w", err)
	}
	addressesV4, err := netinfo.GetAllLocalV4Addresses()
	if err != nil {
		return fmt.Errorf("failed to get all local IPv4 addresses: %w", err)
	}

	provider := winlib.CreateProvider(providerKey, providerDName, "", isPersistant)
	sublayer := winlib.CreateSubLayer(sublayerKey, providerKey,
		sublayerDName, "",
		0xFFF0, // The weight of current layer should be smaller than 0xFFFF (The layer of split-tunneling driver using weight 0xFFFF)
		isPersistant)

	// add provider
	pinfo, err := manager.GetProviderInfo(providerKey)
	if err != nil {
		return fmt.Errorf("failed to get provider info: %w", err)
	}
	if !pinfo.IsInstalled {
		if err = manager.AddProvider(provider); err != nil {
			return fmt.Errorf("failed to add provider : %w", err)
		}
	}

	// add sublayer
	installed, err := manager.IsSubLayerInstalled(sublayerKey)
	if err != nil {
		return fmt.Errorf("failed to check sublayer is installed: %w", err)
	}
	if !installed {
		if err = manager.AddSubLayer(sublayer); err != nil {
			return fmt.Errorf("failed to add sublayer: %w", err)
		}
	}

	// IPv6 filters
	for _, layer := range v6Layers {
		// block all
		_, err := manager.AddFilter(winlib.NewFilterBlockAll(providerKey, layer, sublayerKey, filterDName, "", true, isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'block all IPv6': %w", err)
		}

		// block DNS
		_, err = manager.AddFilter(winlib.NewFilterBlockDNS(providerKey, layer, sublayerKey, sublayerDName, "", nil, isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'block dns': %w", err)
		}

		ipv6loopback := net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}     // LOOPBACK 		::1/128
		ipv6llocal := net.IP{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} // LINKLOCAL		fe80::/10
		// ipv6slocal := net.IP{0xfe, 0xc0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} // SITELOCAL	fec0::/10
		// ipv6ulocal := net.IP{0xfd, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}    // UNIQUELOCAL	fd00::/8
		_, err = manager.AddFilter(winlib.NewFilterAllowRemoteIPV6(providerKey, layer, sublayerKey, filterDName, "", ipv6loopback, 128, isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'allow remote IP' for ipv6loopback: %w", err)
		}
		_, err = manager.AddFilter(winlib.NewFilterAllowRemoteIPV6(providerKey, layer, sublayerKey, filterDName, "", ipv6llocal, 10, isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'allow remote IP' for ipv6llocal: %w", err)
		}

		for _, ip := range addressesV6 {
			if isAllowLAN {
				prefixLen, _ := ip.Mask.Size()
				_, err = manager.AddFilter(winlib.NewFilterAllowRemoteIPV6(providerKey, layer, sublayerKey, filterDName, "", ip.IP, byte(prefixLen), isPersistant))
				if err != nil {
					return fmt.Errorf("failed to add filter 'allow lan IPv6': %w", err)
				}
			}
		}
	}

	// IPv4 filters
	for _, layer := range v4Layers {
		// block all
		_, err := manager.AddFilter(winlib.NewFilterBlockAll(providerKey, layer, sublayerKey, filterDName, "", false, isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'block all': %w", err)
		}

		// block DNS
		_, err = manager.AddFilter(winlib.NewFilterBlockDNS(providerKey, layer, sublayerKey, sublayerDName, "", customDNS, isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'block dns': %w", err)
		}

		// allow DHCP port
		_, err = manager.AddFilter(winlib.NewFilterAllowLocalPort(providerKey, layer, sublayerKey, sublayerDName, "", 68, isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'allow dhcp': %w", err)
		}

		// allow current executable
		binaryPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to obtain executable info: %w", err)
		}
		_, err = manager.AddFilter(winlib.NewFilterAllowApplication(providerKey, layer, sublayerKey, sublayerDName, "", binaryPath, isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'allow application': %w", err)
		}

		// allow OpenVPN executable
		_, err = manager.AddFilter(winlib.NewFilterAllowApplication(providerKey, layer, sublayerKey, sublayerDName, "", platform.OpenVpnBinaryPath(), isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'allow application - openvpn': %w", err)
		}
		// allow WireGuard executable
		_, err = manager.AddFilter(winlib.NewFilterAllowApplication(providerKey, layer, sublayerKey, sublayerDName, "", platform.WgBinaryPath(), isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'allow application - wireguard': %w", err)
		}
		// allow obfsproxy
		_, err = manager.AddFilter(winlib.NewFilterAllowApplication(providerKey, layer, sublayerKey, sublayerDName, "", platform.ObfsproxyStartScript(), isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'allow application - obfsproxy': %w", err)
		}

		_, err = manager.AddFilter(winlib.NewFilterAllowRemoteIP(providerKey, layer, sublayerKey, filterDName, "", net.ParseIP("127.0.0.1"), net.IPv4(255, 255, 255, 255), isPersistant))
		if err != nil {
			return fmt.Errorf("failed to add filter 'allow remote IP': %w", err)
		}

		for _, ip := range addressesV4 {
			if isAllowLAN {
				_, err = manager.AddFilter(winlib.NewFilterAllowRemoteIP(providerKey, layer, sublayerKey, filterDName, "", ip.IP, net.IP(ip.Mask), isPersistant))
				if err != nil {
					return fmt.Errorf("failed to add filter 'allow LAN': %w", err)
				}
			}
		}

		if isAllowLANMulticast {
			_, err = manager.AddFilter(winlib.NewFilterAllowRemoteIP(providerKey, layer, sublayerKey, filterDName, "",
				net.IPv4(224, 0, 0, 0), net.IPv4(240, 0, 0, 0), isPersistant))
			if err != nil {
				return fmt.Errorf("failed to add filter 'allow lan-multicast': %w", err)
			}
		}
	}

	return nil
}

func doDisable() error {
	enabled, err := implGetEnabled()
	if err != nil {
		return fmt.Errorf("failed to get info if firewall is on: %w", err)
	}
	if !enabled {
		return nil
	}

	// delete filters
	for _, l := range v6Layers {
		// delete filters and callouts registered for the provider+layer
		if err := manager.DeleteFilterByProviderKey(providerKey, l); err != nil {
			return fmt.Errorf("failed to delete filter : %w", err)
		}
	}

	for _, l := range v4Layers {
		// delete filters and callouts registered for the provider+layer
		if err := manager.DeleteFilterByProviderKey(providerKey, l); err != nil {
			return fmt.Errorf("failed to delete filter : %w", err)
		}
	}

	// delete sublayer
	installed, err := manager.IsSubLayerInstalled(sublayerKey)
	if err != nil {
		return fmt.Errorf("failed to check is sublayer installed : %w", err)
	}
	if installed {
		if err := manager.DeleteSubLayer(sublayerKey); err != nil {
			return fmt.Errorf("failed to delete sublayer : %w", err)
		}
	}

	// delete provider
	pinfo, err := manager.GetProviderInfo(providerKey)
	if err != nil {
		return fmt.Errorf("failed to get provider info : %w", err)
	}
	if pinfo.IsInstalled {
		if err := manager.DeleteProvider(providerKey); err != nil {
			return fmt.Errorf("failed to delete provider : %w", err)
		}
	}

	clientLocalIPFilterIDs = nil

	return nil
}

func doAddClientIPFilters(clientLocalIP net.IP, clientLocalIPv6 net.IP) (retErr error) {
	if clientLocalIP == nil {
		return nil
	}

	enabled, err := implGetEnabled()
	if err != nil {
		return fmt.Errorf("failed to get info if firewall is on: %w", err)
	}
	if enabled == false {
		return nil
	}

	filters := make([]uint64, 0, len(v4Layers))
	for _, layer := range v4Layers {
		f := winlib.NewFilterAllowLocalIP(providerKey, layer, sublayerKey, filterDName, "", clientLocalIP, net.IPv4(255, 255, 255, 255), false)
		id, err := manager.AddFilter(f)
		if err != nil {
			return fmt.Errorf("failed to add filter : %w", err)
		}
		filters = append(filters, id)
	}

	// IPv6: allow IPv6 communication inside tunnel
	if clientLocalIPv6 != nil {
		for _, layer := range v6Layers {
			f := winlib.NewFilterAllowLocalIPV6(providerKey, layer, sublayerKey, filterDName, "", clientLocalIPv6, byte(128), false)
			id, err := manager.AddFilter(f)
			if err != nil {
				return fmt.Errorf("failed to add IPv6 filter : %w", err)
			}
			filters = append(filters, id)
		}
	}

	clientLocalIPFilterIDs = filters

	return nil
}

func doRemoveClientIPFilters() (retErr error) {
	defer func() {
		clientLocalIPFilterIDs = nil
	}()

	enabled, err := implGetEnabled()
	if err != nil {
		return fmt.Errorf("failed to get info if firewall is on: %w", err)
	}
	if enabled == false {
		return nil
	}

	for _, filterID := range clientLocalIPFilterIDs {
		err := manager.DeleteFilterByID(filterID)
		if err != nil {
			return fmt.Errorf("failed to delete filter : %w", err)
		}
	}

	return nil
}

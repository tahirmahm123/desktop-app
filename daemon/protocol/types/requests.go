
package types

import (
	"github.com/tahirmahm123/vpn-desktop-app/daemon/api/types"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/dns"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/vpn"
)

// Hello is an initial request
type Hello struct {
	CommandBase
	// connected client version
	Version string
	Secret  uint64

	// GetServersList == true - client requests to send back info about all servers
	GetServersList bool

	// GetStatus == true - client requests current status (Vpn connection, Firewal... etc.)
	GetStatus bool

	// GetConfigParams == true - client requests config parameters (user-defined OpevVPN file location ... etc.)
	GetConfigParams bool

	// GetSplitTunnelStatus == true - client requests configuration of SplitTunnelling
	GetSplitTunnelStatus bool

	// GetWiFiCurrentState == true - client requests info about current WiFi
	GetWiFiCurrentState bool

	//	KeepDaemonAlone == false (default) - VPN disconnects when client disconnects from a daemon
	//	KeepDaemonAlone == true - do nothing when client disconnects from a daemon (if VPN is connected - do not disconnect)
	KeepDaemonAlone bool
}

// GetServers request servers list
type GetServers struct {
	CommandBase
}

// PingServers request to ping servers
type PingServers struct {
	CommandBase
	RetryCount int
	TimeOutMs  int
}

// KillSwitchSetAllowLANMulticast enable\disable LAN multicast acces for kill-switch
type KillSwitchSetAllowLANMulticast struct {
	CommandBase
	AllowLANMulticast bool

	// When true - deamon returns empty response as confirmation
	// Needed for supporting old UI clients which are don't require confirmation
	Synchronously bool
}

// KillSwitchSetAllowLAN enable\disable LAN acces for kill-switch
type KillSwitchSetAllowLAN struct {
	CommandBase
	AllowLAN bool

	// When true - deamon returns empty response as confirmation
	// Needed for supporting old UI clients which are don't require confirmation
	Synchronously bool
}

type KillSwitchSetAllowApiServers struct {
	CommandBase
	IsAllowApiServers bool
}

// KillSwitchSetEnabled request to enable\disable kill-switch
type KillSwitchSetEnabled struct {
	CommandBase
	IsEnabled bool
}

// KillSwitchGetStatus get full killswitch status
type KillSwitchGetStatus struct {
	CommandBase
}

// KillSwitchSetIsPersistent request to mark kill-switch persistant
type KillSwitchSetIsPersistent struct {
	CommandBase
	IsPersistent bool
}

// SetPreference sets daemon configuration parameter
type SetPreference struct {
	CommandBase
	Key   string
	Value string
}

// SetAlternateDns request to set custom DNS
type SetAlternateDns struct {
	CommandBase
	Dns dns.DnsSettings
}

// GetDnsPredefinedConfigs request to get list of predefined DoH/DoT configurations (if exists)
type GetDnsPredefinedConfigs struct {
	CommandBase
}

// Connect request to establish new VPN connection
type Connect struct {
	CommandBase
	// Can use IPv6 connection inside tunnel
	// The hosts which support IPv6 have higher priority,
	// but if there are no IPv6 hosts - we will use the IPv4 host.
	IPv6 bool
	// Use ONLY IPv6 hosts (ignored when IPv6!=true)
	IPv6Only  bool
	VpnType   vpn.Type
	ManualDNS dns.DnsSettings

	// Enable firewall before connection
	// (if true - the parameter 'firewallDuringConnection' will be ignored)
	FirewallOn bool
	// Enable firewall before connection and disable after disconnection
	// (has effect only if Firewall not enabled before)
	FirewallOnDuringConnection bool

	WireGuardParameters struct {
		Port struct {
			Port int
		}

		EntryVpnServer struct {
			Hosts []types.WireGuardServerHostInfo
		}

		MultihopExitServer struct {
			// ExitSrvID (geteway ID) just in use to keep clients notified about connected MH exit server
			// in same manner as for OpenVPN connection.
			// Example: "gateway":"zz.wg.vpn.net" => "zz"
			ExitSrvID string
			Hosts     []types.WireGuardServerHostInfo
		}
	}

	OpenVpnParameters struct {
		EntryVpnServer struct {
			Hosts []types.OpenVPNServerHostInfo
		}

		// MultihopExitSrvID example: "gateway":"zz.wg.vpn.net" => "zz"
		MultihopExitSrvID string
		ProxyType         string
		ProxyAddress      string
		ProxyPort         int
		ProxyUsername     string
		ProxyPassword     string

		Port struct {
			Port     int
			Protocol int
		}
	}
}

// Disconnect disconnect active VPN connection
type Disconnect struct {
	CommandBase
}

// GetVPNState request daemon to provive current VPN connection state
type GetVPNState struct {
	CommandBase
}

// SessionNew - create new session
//
// When force is set to true - all active sessions will be deleted prior to creating a new one if user reached session limit.
// Initial call to /sessin/new should always be performed with force set to false, to display special form, when sessions limit is reached.
// VPN apps have to set force to true only when customer clicks Log all other clients button.
type SessionNew struct {
	CommandBase
	Username string
	Password string
}

// SessionDelete logout from current device
type SessionDelete struct {
	CommandBase
	NeedToResetSettings   bool
	NeedToDisableFirewall bool
	// If IsCanDeleteSessionLocally==true: the account will be logged out
	// even if there is no connectivity to API server
	IsCanDeleteSessionLocally bool
}

// AccountStatus get account status
type AccountStatus struct {
	CommandBase
}

// WireGuardGenerateNewKeys - generate WG keys
type WireGuardGenerateNewKeys struct {
	CommandBase
	OnlyUpdateIfNecessary bool
}

// WireGuardSetKeysRotationInterval -  change WG keys rotation interval
type WireGuardSetKeysRotationInterval struct {
	CommandBase
	Interval int64
}

// WiFiAvailableNetworks - get list of available WIFI networks
type WiFiAvailableNetworks struct {
	CommandBase
}

// WiFiCurrentNetwork - request info about connected WIFI
type WiFiCurrentNetwork struct {
	CommandBase
}

// IPProtocol - VPN type
type RequiredIPProtocol int

const (
	IPvAny RequiredIPProtocol = 0
	IPv4   RequiredIPProtocol = 1
	IPv6   RequiredIPProtocol = 2
)

// APIRequest do custom request to API
type APIRequest struct {
	CommandBase
	APIPath            string
	IPProtocolRequired RequiredIPProtocol
}

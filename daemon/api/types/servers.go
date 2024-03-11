//
//  Daemon for IVPN Client Desktop
//  https://github.com/tahirmahm123/vpn-desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2023 IVPN Limited.
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

package types

import (
	"fmt"
	"strconv"
	"strings"
)

type OpenVPNProtocol struct {
	Certificate string `json:"certificate"`
	Ports       []struct {
		Protocol string `json:"protocol"`
		Port     int    `json:"port"`
	} `json:"ports"`
}
type DNSServers struct {
	DNS1 string `json:"dns1"`
	DNS2 string `json:"dns2"`
}
type WireGuardInstance struct {
	PublicKey string `json:"PublicKey"`
	Port      int    `json:"port"`
	LocalIP   string `json:"local_ip,omitempty"`
	DNS       string `json:"dnsServers,omitempty"`
}
type OpenVPNInstance struct {
	Protocol string `json:"proto"`
	Port     int    `json:"port"`
}
type ServerListItem struct {
	Id          int                 `json:"id"`
	Name        string              `json:"name"`
	Ip          string              `json:"ip"`
	Port        int                 `json:"port"`
	Flag        string              `json:"flag"`
	DNS1        string              `json:"dns1"`
	DNS2        string              `json:"dns2"`
	Premium     bool                `json:"premium"`
	Country     string              `json:"country"`
	CountryCode string              `json:"country_code"`
	OpenVPN     []OpenVPNInstance   `json:"openvpn"`
	WireGuard   []WireGuardInstance `json:"wg"`
	Location    struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"location"`
}

func (s *ServerListItem) Latitude() float32 {
	value, err := strconv.ParseFloat(s.Location.Latitude, 32)
	if err != nil {
		return 0
	}
	return float32(value)
}
func (s *ServerListItem) Longitude() float32 {
	value, err := strconv.ParseFloat(s.Location.Longitude, 32)
	if err != nil {
		return 0
	}
	return float32(value)
}

type ServerListCountryItem struct {
	Flag    string           `json:"flag"`
	Country string           `json:"country"`
	Hosts   []ServerListItem `json:"servers"`
}
type ServerListProtoItem struct {
	OpenVPNServers   []ServerListCountryItem `json:"openvpn"`
	WireGuardServers []ServerListCountryItem `json:"wireguard"`
}
type ServerListResponse struct {
	ServerList ServerListProtoItem `json:"servers,omitempty"`
	DnsServers DNSServers          `json:"dnsServers"`
	OpenVPN    OpenVPNProtocol     `json:"openvpn"`
	WireGuard  []int               `json:"wireguard"`
}

// -----------------------------------------------------------
type ServerGeneric interface {
	GetServerInfoBase() ServerInfoBase
	GetHostsInfoBase() []HostInfoBase
}

type HostInfoBase struct {
	Hostname     string  `json:"hostname"`
	Host         string  `json:"host"`
	DnsName      string  `json:"dns_name"`
	MultihopPort int     `json:"multihop_port"`
	Load         float32 `json:"load"`
	V2RayHost    string  `json:"v2ray"`
}

func (h HostInfoBase) GetHostInfoBase() HostInfoBase {
	return h
}

type ServerInfoBase struct {
	Gateway     string `json:"gateway"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
	City        string `json:"city"`

	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`

	ISP string `json:"isp"`
}

func (s ServerInfoBase) GetServerInfoBase() ServerInfoBase {
	return s
}

// -----------------------------------------------------------

type WireGuardServerHostInfoIPv6 struct {
	Host    string `json:"host"`
	LocalIP string `json:"local_ip"`
}

// WireGuardServerHostInfo contains info about WG server host
type WireGuardServerHostInfo struct {
	HostInfoBase
	PublicKey string                      `json:"public_key"`
	LocalIP   string                      `json:"local_ip"`
	IPv6      WireGuardServerHostInfoIPv6 `json:"ipv6"`
}

// WireGuardServerInfo contains all info about WG server
type WireGuardServerInfo struct {
	ServerInfoBase
	Hosts []WireGuardServerHostInfo `json:"hosts"`
}

func (s WireGuardServerInfo) GetHostsInfoBase() []HostInfoBase {
	ret := []HostInfoBase{}
	for _, host := range s.Hosts {
		ret = append(ret, host.HostInfoBase)
	}
	return ret
}

// -----------------------------------------------------------

type ObfsParams struct {
	Obfs3MultihopPort int    `json:"obfs3_multihop_port"`
	Obfs4MultihopPort int    `json:"obfs4_multihop_port"`
	Obfs4Key          string `json:"obfs4_key"`
}

// OpenVPNServerHostInfo contains info about OpenVPN server host
type OpenVPNServerHostInfo struct {
	HostInfoBase
	Obfs ObfsParams `json:"obfs"`
}

// OpenvpnServerInfo contains all info about OpenVPN server
type OpenvpnServerInfo struct {
	ServerInfoBase
	Hosts []OpenVPNServerHostInfo `json:"hosts"`
}

func (s OpenvpnServerInfo) GetHostsInfoBase() []HostInfoBase {
	ret := []HostInfoBase{}
	for _, host := range s.Hosts {
		ret = append(ret, host.HostInfoBase)
	}
	return ret
}

// -----------------------------------------------------------

// DNSInfo contains info about DNS server
type DNSInfo struct {
	IP string `json:"ip"`
}

// AntitrackerInfo all info about antitracker DNSs
type AntitrackerInfo struct {
	Default  DNSInfo `json:"default"`
	Hardcore DNSInfo `json:"hardcore"`
}

// -----------------------------------------------------------

type AntiTrackerPlusServer struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Normal      string `json:"Normal"`
	Hardcore    string `json:"Hardcore"`
}

type AntiTrackerPlusInfo struct {
	DnsServers []AntiTrackerPlusServer `json:"DnsServers"`
}

// -----------------------------------------------------------

// InfoAPI contains API IP adresses
type InfoAPI struct {
	IPAddresses   []string `json:"ips"`
	IPv6Addresses []string `json:"ipv6s"`
}

type PortRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type PortInfoBase struct {
	Type string `json:"type"` // "TCP" or "UDP"
	Port int    `json:"port"`
}

type PortInfo struct {
	PortInfoBase
	Range PortRange `json:"range"`
}

func (pi PortInfo) String() string {
	if pi.Port > 0 {
		return fmt.Sprintf("%s:%d", pi.Type, pi.Port)
	}
	if pi.Range.Min > 0 && pi.Range.Min < pi.Range.Max {
		return fmt.Sprintf("%s:[%d-%d]", pi.Type, pi.Range.Min, pi.Range.Max)
	}
	return ""
}

func (pi PortInfo) IsTCP() bool {
	return strings.TrimSpace(strings.ToLower(pi.Type)) == "tcp"
}

func (pi PortInfo) IsUDP() bool {
	return strings.TrimSpace(strings.ToLower(pi.Type)) == "udp"
}

func (pi PortInfo) Equal(x PortInfo) bool {
	return pi.Port == x.Port &&
		strings.TrimSpace(strings.ToLower(pi.Type)) == strings.TrimSpace(strings.ToLower(x.Type)) &&
		pi.Range.Max == x.Range.Max && pi.Range.Min == x.Range.Min
}

type ObfsPortInfo struct {
	Port int `json:"port"`
}

type EchoServer struct {
	EchoServer string `json:"echoserver"`
}

// V2Ray ports structure
type V2Ray struct {
	ID        string         `json:"id"`
	OpenVPN   []PortInfoBase `json:"openvpn"`
	WireGuard []PortInfoBase `json:"wireguard"`
}

type PortsInfo struct {
	OpenVPN   []PortInfo   `json:"openvpn"`
	WireGuard []PortInfo   `json:"wireguard"`
	Obfs3     ObfsPortInfo `json:"obfs3"`
	Obfs4     ObfsPortInfo `json:"obfs4"`
	Test      []EchoServer `json:"test"`
	V2Ray     V2Ray        `json:"v2ray"`
}

// ConfigInfo contains different configuration info (Antitracker, API ...)
type ConfigInfo struct {
	Antitracker     AntitrackerInfo     `json:"antitracker"`
	AntiTrackerPlus AntiTrackerPlusInfo `json:"antitracker_plus"`
	API             InfoAPI             `json:"api"`
	Ports           PortsInfo           `json:"ports"`
}

// ServerInfoResponse all info from servers.json
type ServerInfoResponse struct {
	WireguardServers []WireGuardServerInfo `json:"wireguard"`
	OpenvpnServers   []OpenvpnServerInfo   `json:"openvpn"`
	Config           ConfigInfo            `json:"config"`
}

func (si ServerListResponse) ServersGenericWireguard() (ret []ServerListCountryItem) {
	for _, s := range si.ServerList.WireGuardServers {
		ret = append(ret, s)
	}
	return
}

func (si ServerListResponse) ServersGenericOpenvpn() (ret []ServerListCountryItem) {
	for _, s := range si.ServerList.OpenVPNServers {
		ret = append(ret, s)
	}
	return
}

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

package types

type WireGuardServerHostInfoIPv6 struct {
	Host    string `json:"host"`
	LocalIP string `json:"local_ip"`
}

// WireGuardServerHostInfo contains info about WG server host
type WireGuardServerHostInfo struct {
	Hostname     string                      `json:"hostname"`
	Host         string                      `json:"host"`
	PublicKey    string                      `json:"public_key"`
	LocalIP      string                      `json:"local_ip"`
	IPv6         WireGuardServerHostInfoIPv6 `json:"ipv6"`
	MultihopPort int                         `json:"multihop_port"`
}

// WireGuardServerInfo contains all info about WG server
type WireGuardServerInfo struct {
	Gateway     string `json:"gateway"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
	City        string `json:"city"`

	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`

	Hosts []WireGuardServerHostInfo `json:"hosts"`
}

// OpenVPNServerHostInfo contains info about OpenVPN server host
type OpenVPNServerHostInfo struct {
	Hostname     string `json:"hostname"`
	Host         string `json:"ip"`
	MultihopPort int    `json:"multihop_port:omitempty"`
}

// OpenvpnServerInfo contains all info about OpenVPN server
type OpenvpnServerInfo struct {
	Gateway     string `json:"gateway"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
	City        string `json:"city"`

	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`

	Hosts []OpenVPNServerHostInfo `json:"hosts"`
}

// DNSInfo contains info about DNS server
type DNSInfo struct {
	IP         string `json:"ip"`
	MultihopIP string `json:"multihop-ip"`
}

// AntitrackerInfo all info about antitracker DNSs
type AntitrackerInfo struct {
	Default  DNSInfo `json:"default"`
	Hardcore DNSInfo `json:"hardcore"`
}

// InfoAPI contains API IP adresses
type InfoAPI struct {
	IPAddresses   []string `json:"ips"`
	IPv6Addresses []string `json:"ipv6s"`
}

// ConfigInfo contains different configuration info (Antitracker, API ...)
type ConfigInfo struct {
	Antitracker AntitrackerInfo `json:"antitracker"`
	API         InfoAPI         `json:"api"`
}

// ServersInfoResponse all info from servers.json
type ServersInfoResponse struct {
	OpenVPNConfig string            `json:"certificate"`
	Servers       []ServerByCountry `json:"data"`
}

type ServerByCountry struct {
	Flag    string         `json:"flag"`
	Country string         `json:"country"`
	Hosts   []ServerObject `json:"hosts"`
}

type ServerObject struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Country string `json:"country"`
}

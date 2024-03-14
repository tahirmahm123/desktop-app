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
	"strconv"
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

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

import "strconv"

// APIResponse - generic API response
type APIResponse struct {
	Status int `json:"status"` // status code
}

// APIErrorResponse generic IVPN API error
type APIErrorResponse struct {
	APIResponse
	Message string `json:"message,omitempty"` // Text description of the message
}

// KemCiphers in use for KEM: to exchange WG PresharedKey
type KemCiphers struct {
	KemCipher_Kyber1024             string `json:"kem_cipher1,omitempty"` // (Kyber-1024) in use for KEM: to exchange WG PresharedKey
	KemCipher_ClassicMcEliece348864 string `json:"kem_cipher2,omitempty"` // (Classic-McEliece-348864) in use for KEM: to exchange WG PresharedKey
}

// GeoLookupResponse geolocation info
type GeoLookupResponse struct {
	//ip_address   string
	//isp          string
	//organization string
	//country      string
	//country_code string
	//city         string

	SLatitude  string `json:"latitude"`
	SLongitude string `json:"longitude"`

	//isIvpnServer bool
}

func (s *GeoLookupResponse) Latitude() float32 {
	value, err := strconv.ParseFloat(s.SLatitude, 32)
	if err != nil {
		return 0
	}
	return float32(value)
}
func (s *GeoLookupResponse) Longitude() float32 {
	value, err := strconv.ParseFloat(s.SLongitude, 32)
	if err != nil {
		return 0
	}
	return float32(value)
}

type WGKeysUpdateResponse struct {
	Message string `json:"message"`
	LocalIP string `json:"localIP"`
}
type PinValidationResponse struct {
	Status     string `json:"status"`
	Code       string `json:"code"`
	Token      string `json:"token,omitempty"`
	ExpiryDate string `json:"expiry_date"`
	Timestamp  int64  `json:"timestamp"`
}

const (
	ValidPin   = "Valid"
	InvalidPin = "Invalid"
	ExpiredPin = "Expired"
)

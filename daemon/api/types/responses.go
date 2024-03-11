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

// APIResponse - generic API response
type APIResponse struct {
	Status int `json:"status"` // status code
}

// APIErrorResponse generic IVPN API error
type APIErrorResponse struct {
	APIResponse
	Message string `json:"message,omitempty"` // Text description of the message
}

// ServiceStatusAPIResp account info
type ServiceStatusAPIResp struct {
	Active         bool     `json:"is_active"`
	ActiveUntil    int64    `json:"active_until"`
	CurrentPlan    string   `json:"current_plan"`
	PaymentMethod  string   `json:"payment_method"`
	IsRenewable    bool     `json:"is_renewable"`
	WillAutoRebill bool     `json:"will_auto_rebill"`
	IsFreeTrial    bool     `json:"is_on_free_trial"`
	Capabilities   []string `json:"capabilities"`
	Upgradable     bool     `json:"upgradable"`
	UpgradeToPlan  string   `json:"upgrade_to_plan"`
	UpgradeToURL   string   `json:"upgrade_to_url"`
	Limit          int      `json:"limit"` // applicable for 'session limit' error
}

// KemCiphers in use for KEM: to exchange WG PresharedKey
type KemCiphers struct {
	KemCipher_Kyber1024             string `json:"kem_cipher1,omitempty"` // (Kyber-1024) in use for KEM: to exchange WG PresharedKey
	KemCipher_ClassicMcEliece348864 string `json:"kem_cipher2,omitempty"` // (Classic-McEliece-348864) in use for KEM: to exchange WG PresharedKey
}

// SessionNewResponse information about created session
type SessionNewResponse struct {
	APIErrorResponse
	Token       string `json:"token"`
	VpnUsername string `json:"vpn_username"`
	VpnPassword string `json:"vpn_password"`

	CaptchaID    string `json:"captcha_id"`
	CaptchaImage string `json:"captcha_image"`

	ServiceStatus ServiceStatusAPIResp `json:"service_status"`

	WireGuard struct {
		Status    int    `json:"status"`
		Message   string `json:"message,omitempty"`
		IPAddress string `json:"ip_address,omitempty"`
		KemCiphers
	} `json:"wireguard"`
}

// SessionNewErrorLimitResponse information about session limit error
type SessionNewErrorLimitResponse struct {
	APIErrorResponse
	SessionLimitData ServiceStatusAPIResp `json:"data"`
}

// SessionsWireGuardResponse Sessions WireGuard response
type SessionsWireGuardResponse struct {
	APIErrorResponse
	IPAddress string `json:"ip_address,omitempty"`
	KemCiphers
}

// SessionStatusResponse session status response
type SessionStatusResponse struct {
	APIErrorResponse
	ServiceStatus ServiceStatusAPIResp `json:"service_status"`
}

// GeoLookupResponse geolocation info
type GeoLookupResponse struct {
	//ip_address   string
	//isp          string
	//organization string
	//country      string
	//country_code string
	//city         string

	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`

	//isIvpnServer bool
}

type UserLoginResponse struct {
	Response UserLoginResponseDetails `json:"response"`
	Message  string                   `json:"message"`
	State    bool                     `json:"state"`
}
type UserLoginResponseDetails struct {
	Auth                bool            `json:"auth"`
	Active              bool            `json:"active"`
	Expired             bool            `json:"expired"`
	AllowLogin          bool            `json:"allowLogin"`
	ExpiryDate          string          `json:"expiry_date"`
	Token               string          `json:"ApiToken"`
	TotalSessionAllowed int             `json:"totalSessionAllowed"`
	LoggedInSessions    int             `json:"loggedInSessions"`
	IsPaid              bool            `json:"is_paid"`
	ActiveSessions      []DeviceSession `json:"activeSessions"`
	ActiveUntil         int64           `json:"timestamp"`
	Email               string          `json:"email"`
	Plan                string          `json:"plan"`
	VPNUsername         string          `json:"vpn_username"`
	VPNPassword         string          `json:"vpn_password"`
	WGLocalIP           string          `json:"localIP"`
}
type DeviceSession struct {
	TokenId       int                  `json:"tokenId"`
	Details       DeviceSessionDetails `json:"details,omitempty"`
	LastUsedAt    string               `json:"_last_used_at"`
	CurrentDevice bool                 `json:"currentSession"`
}
type DeviceSessionDetails struct {
	DeviceId   string `json:"id"`
	DeviceName string `json:"name"`
	DeviceType string `json:"type"`
}

type UserDetailsResponse struct {
	Plan                string          `json:"plan"`
	Active              bool            `json:"active"`
	Expired             bool            `json:"expired"`
	ActiveUntil         int64           `json:"timestamp"`
	TotalSessionAllowed int             `json:"totalSessionAllowed"`
	LoggedInSessions    int             `json:"loggedInSessions"`
	ActiveSessions      []DeviceSession `json:"activeSessions"`
}

type UserLogoutResponse struct {
	Message string `json:"message"`
}
type UserLogoutDeviceResponse struct {
	ActiveSessions      []DeviceSession `json:"activeSessions"`
	TotalSessionAllowed int             `json:"totalSessionsAllowed,omitempty"`
	LoggedInSessions    int             `json:"loggedInSessions,omitempty"`
}
type WGKeysUpdateResponse struct {
	Message string `json:"message"`
	LocalIP string `json:"localIP"`
}

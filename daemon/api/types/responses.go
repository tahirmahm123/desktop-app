
package types

// APIResponse - generic API response
type APIResponse struct {
	Status int `json:"status"` // status code
}

// APIErrorResponse generic VPNror
type APIErrorResponse struct {
	APIResponse
	Message string `json:"message,omitempty"` // Text description of the message
}

// ServiceStatusAPIResp account info
type ServiceStatusAPIResp struct {
	Active      bool   `json:"is_active"`
	ActiveUntil int64  `json:"expiry_date"`
	IsFreeTrial bool   `json:"trial"`
	CurrentPlan string `json:"plan"`
}

// SessionNewResponse information about created session
type SessionNewResponse struct {
	APIErrorResponse
	Token       string `json:"ApiToken"`
	VpnUsername string `json:"vpn_username"`
	VpnPassword string `json:"vpn_password"`

	ServiceStatus ServiceStatusAPIResp `json:"service_status:omitempty"`

	Authenticated bool   `json:"auth"`
	Active        bool   `json:"active"`
	IsExpired     bool   `json:"expired"`
	ExpiryDate    string `json:"expiry_date"`
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

	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`

	//isIvpnServer bool
}

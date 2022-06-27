
package preferences

import (
	"net"
	"strings"
	"time"
)

// SessionStatus contains information about current session
type SessionStatus struct {
	AccountID          string
	Session            string `json:",omitempty"`
	OpenVPNUser        string `json:",omitempty"`
	OpenVPNPass        string `json:",omitempty"`
	WGPublicKey        string
	WGPrivateKey       string `json:",omitempty"`
	WGLocalIP          string
	WGKeyGenerated     time.Time
	WGKeysRegenInerval time.Duration // syntax error in variable name. Keeping it as is for compatibility with previous versions
}

// IsLoggedIn returns 'true' when user logged-in
func (s *SessionStatus) IsLoggedIn() bool {
	return len(s.Session) != 0
}

// IsWGCredentialsOk returns 'true' when WireGuard credentials are initialized
func (s *SessionStatus) IsWGCredentialsOk() bool {
	if len(s.WGPublicKey) == 0 || len(s.WGLocalIP) == 0 || len(s.WGPrivateKey) == 0 {
		return false
	}
	return true
}

func (s *SessionStatus) updateWgCredentials(wgPublicKey string, wgPrivateKey string, wgLocalIP string) {
	if len(wgLocalIP) > 0 {
		if net.ParseIP(wgLocalIP) == nil {
			log.Error("Unable to save WG credentials (local IP has wrong format)")
			wgLocalIP = ""
		}
	}

	s.WGPublicKey = strings.TrimSpace(wgPublicKey)
	s.WGPrivateKey = strings.TrimSpace(wgPrivateKey)
	s.WGLocalIP = strings.TrimSpace(wgLocalIP)

	if len(s.WGPublicKey) > 0 && len(s.WGPrivateKey) > 0 && len(s.WGLocalIP) > 0 {
		s.WGKeyGenerated = time.Now()
	} else {
		s.WGKeyGenerated = time.Time{}
	}
}

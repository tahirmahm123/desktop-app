package preferences

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ivpn/desktop-app-daemon/logger"
	"github.com/ivpn/desktop-app-daemon/service/platform"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("sprefs")
}

// Preferences - IVPN service preferences
type Preferences struct {
	IsLogging                bool
	IsFwPersistant           bool
	IsFwAllowLAN             bool
	IsFwAllowLANMulticast    bool
	IsStopOnClientDisconnect bool
	IsObfsproxy              bool
	OpenVpnExtraParameters   string

	// last known account status
	//Account AccountStatus
	Session SessionStatus
}

// SetSession save account credentials
func (p *Preferences) SetSession(accountID string,
	session string,
	vpnUser string,
	vpnPass string,
	wgPublicKey string,
	wgPrivateKey string,
	wgLocalIP string) {

	p.setSession(accountID, session, vpnUser, vpnPass, wgPublicKey, wgPrivateKey, wgLocalIP)
	p.SavePreferences()
}

// UpdateWgCredentials save wireguard credentials
func (p *Preferences) UpdateWgCredentials(wgPublicKey string, wgPrivateKey string, wgLocalIP string) {
	p.Session.updateWgCredentials(wgPublicKey, wgPrivateKey, wgLocalIP)
	p.SavePreferences()
}

// SavePreferences saves preferences
func (p *Preferences) SavePreferences() error {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to save preferences file (json marshal error): %w", err)
	}

	return ioutil.WriteFile(platform.SettingsFile(), data, 0644)
}

// LoadPreferences loads preferences
func (p *Preferences) LoadPreferences() error {
	data, err := ioutil.ReadFile(platform.SettingsFile())

	if err != nil {
		return fmt.Errorf("failed to read preferences file: %w", err)
	}

	dataStr := string(data)
	if strings.Contains(dataStr, `"firewall_is_persistent"`) {
		// It is a first time loading preferences after IVPN Client upgrade from old version (<= v2.10.9)
		// Loading preferences with an old parameter names and types:
		type PreferencesOld struct {
			IsLogging                string `json:"enable_logging"`
			IsFwPersistant           string `json:"firewall_is_persistent"`
			IsFwAllowLAN             string `json:"firewall_allow_lan"`
			IsFwAllowLANMulticast    string `json:"firewall_allow_lan_multicast"`
			IsStopOnClientDisconnect string `json:"is_stop_server_on_client_disconnect"`
			IsObfsproxy              string `json:"enable_obfsproxy"`
			OpenVpnExtraParameters   string `json:"open_vpn_extra_parameters"`
		}
		oldStylePrefs := &PreferencesOld{}

		if err := json.Unmarshal(data, oldStylePrefs); err != nil {
			return err
		}

		p.IsLogging = oldStylePrefs.IsLogging == "1"
		p.IsFwPersistant = oldStylePrefs.IsFwPersistant == "1"
		p.IsFwAllowLAN = oldStylePrefs.IsFwAllowLAN == "1"
		p.IsFwAllowLANMulticast = oldStylePrefs.IsFwAllowLANMulticast == "1"
		p.IsStopOnClientDisconnect = oldStylePrefs.IsStopOnClientDisconnect == "1"
		p.IsObfsproxy = oldStylePrefs.IsObfsproxy == "1"
		p.OpenVpnExtraParameters = oldStylePrefs.OpenVpnExtraParameters

		return nil
	}

	err = json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	if len(p.Session.WGPublicKey) == 0 || len(p.Session.WGPrivateKey) == 0 || len(p.Session.WGLocalIP) == 0 {
		p.Session.WGKeyGenerated = time.Time{}
	}

	if p.Session.WGKeysRegenInerval <= 0 {
		p.Session.WGKeysRegenInerval = time.Hour * 24 * 7
		log.Info(fmt.Sprintf("default value for preferences: WgKeysRegenInervalDays=%v", p.Session.WGKeysRegenInerval))
		p.SavePreferences()
	}

	return nil
}

func (p *Preferences) setSession(accountID string,
	session string,
	vpnUser string,
	vpnPass string,
	wgPublicKey string,
	wgPrivateKey string,
	wgLocalIP string) {

	p.Session = SessionStatus{
		AccountID:          strings.TrimSpace(accountID),
		Session:            strings.TrimSpace(session),
		OpenVPNUser:        strings.TrimSpace(vpnUser),
		OpenVPNPass:        strings.TrimSpace(vpnPass),
		WGKeysRegenInerval: p.Session.WGKeysRegenInerval} // keep 'WGKeysRegenInerval' from previus Session object

	p.Session.updateWgCredentials(wgPublicKey, wgPrivateKey, wgLocalIP)
}
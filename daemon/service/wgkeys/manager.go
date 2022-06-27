
package wgkeys

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/api"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/api/types"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/vpn"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/vpn/wireguard"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("wgkeys")
}

//HardExpirationIntervalDays = 40;

// IWgKeysChangeReceiver WG key update handler
type IWgKeysChangeReceiver interface {
	WireGuardSaveNewKeys(wgPublicKey string, wgPrivateKey string, wgLocalIP string)
	WireGuardGetKeys() (session, wgPublicKey, wgPrivateKey, wgLocalIP string, generatedTime time.Time, updateInterval time.Duration)
	FirewallEnabled() (bool, error)
	Connected() bool
	ConnectedType() (isConnected bool, connectedVpnType vpn.Type)
	IsConnectivityBlocked() (isBlocked bool, reasonDescription string, err error)
	OnSessionNotFound()
}

// CreateKeysManager create WireGuard keys manager
func CreateKeysManager(apiObj *api.API, wgToolBinPath string) *KeysManager {
	return &KeysManager{
		stopKeysRotation: make(chan struct{}),
		wgToolBinPath:    wgToolBinPath,
		api:              apiObj}
}

// KeysManager WireGuard keys manager
type KeysManager struct {
	mutex            sync.Mutex
	service          IWgKeysChangeReceiver
	api              *api.API
	wgToolBinPath    string
	stopKeysRotation chan struct{}
}

// Init - initialize master service
func (m *KeysManager) Init(receiver IWgKeysChangeReceiver) error {
	if receiver == nil || m.service != nil {
		return fmt.Errorf("failed to initialize WG KeysManager")
	}
	m.service = receiver
	return nil
}

// StartKeysRotation start keys rotation
func (m *KeysManager) StartKeysRotation() error {
	if m.service == nil {
		return fmt.Errorf("unable to start WG keys rotation (KeysManager not initialized)")
	}

	m.StopKeysRotation()

	_, activePublicKey, _, _, lastUpdate, interval := m.service.WireGuardGetKeys()

	if interval <= 0 {
		return fmt.Errorf("unable to start WG keys rotation (update interval not defined)")
	}

	if len(activePublicKey) == 0 {
		log.Info("Active public WG key is not defined. WG key rotation disabled.")
		return nil
	}

	go func() {
		log.Info(fmt.Sprintf("Keys rotation started (interval:%v)", interval))
		defer log.Info("Keys rotation stopped")

		const maxCheckInterval = time.Minute * 5

		needStop := false

		isLastUpdateFailed := false
		isLastUpdateFailedCnt := 0

		for !needStop {
			waitInterval := maxCheckInterval

			if isLastUpdateFailed {
				// If the last update failed - do next try after some delay
				// (delay is incrteasing every retry: 5, 10, 15 ... 60 min)
				waitInterval = maxCheckInterval * time.Duration(isLastUpdateFailedCnt)
				if waitInterval > time.Hour {
					waitInterval = time.Hour
				}
				lastUpdate = time.Now()
			} else {
				_, _, _, _, lastUpdate, interval = m.service.WireGuardGetKeys()
				waitInterval = time.Until(lastUpdate.Add(interval))
			}

			// update immediately, if it is a time
			if lastUpdate.Add(interval).Before(time.Now()) {
				waitInterval = time.Second
			}

			if waitInterval > maxCheckInterval && !isLastUpdateFailed {
				// We can not trust "time.After()" that it will be triggered in exact time.
				// If the computer fall to sleep on a long time, after wake up the "time.After()"
				// will trigger after [sleep time]+[time].
				// Therefore we defining maximum allowed interval to check necessity on keys generation
				waitInterval = maxCheckInterval
			}

			select {
			case <-time.After(waitInterval):
				_, err := m.UpdateKeysIfNecessary()
				if err != nil {
					isLastUpdateFailed = true
					isLastUpdateFailedCnt += 1
				} else {
					isLastUpdateFailed = false
					isLastUpdateFailedCnt = 0
				}

			case <-m.stopKeysRotation:
				needStop = true
			}
		}
	}()

	return nil
}

// StopKeysRotation stop keys rotation
func (m *KeysManager) StopKeysRotation() {
	select {
	case m.stopKeysRotation <- struct{}{}:
	default:
	}
}

// GenerateKeys generate keys
func (m *KeysManager) GenerateKeys() error {
	isUpdated, err := m.generateKeys(false)
	if err == nil && !isUpdated {
		err = fmt.Errorf("WG keys were not updated")
	}
	return err
}

// UpdateKeysIfNecessary generate or update keys
// 1) If no active WG keys defined - new keys will be generated + key rotation will be started
// 2) If active WG key defined - key will be updated only if it is a time to do it
func (m *KeysManager) UpdateKeysIfNecessary() (isUpdated bool, retErr error) {
	return m.generateKeys(true)
}

func (m *KeysManager) generateKeys(onlyUpdateIfNecessary bool) (isUpdated bool, retErr error) {
	defer func() {
		if retErr != nil {
			log.Error("Failed to update WG keys: ", retErr)
		}
	}()

	if m.service == nil {
		return false, fmt.Errorf("WG KeysManager not initialized")
	}

	// Check update configuration
	// (not blocked by mutex because in order to return immediately if nothing to do)
	_, activePublicKey, _, _, lastUpdate, interval := m.service.WireGuardGetKeys()

	// function to check if update required
	isNecessaryUpdate := func() (bool, error) {
		if !onlyUpdateIfNecessary {
			return true, nil
		}
		if interval <= 0 {
			// update interval must be defined
			return false, fmt.Errorf("unable to 'GenerateOrUpdateKeys' (update interval is not defined)")
		}
		if len(activePublicKey) > 0 {
			// If active WG key defined - key will be updated only if it is a time to do it
			if lastUpdate.Add(interval).After(time.Now()) {
				// it is not a time to regenerate keys
				return false, nil
			}
		}
		return true, nil
	}

	if haveToUpdate, err := isNecessaryUpdate(); !haveToUpdate || err != nil {
		return false, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check update configuration second time (locked by mutex)
	session, activePublicKey, _, _, lastUpdate, interval := m.service.WireGuardGetKeys()
	if haveToUpdate, err := isNecessaryUpdate(); !haveToUpdate || err != nil {
		return false, err
	}

	isRotationStopped := false
	if len(activePublicKey) == 0 {
		isRotationStopped = true
	}

	log.Info("Updating WG keys...")

	if isBlocked, reasonDescription, err := m.service.IsConnectivityBlocked(); err == nil && isBlocked {
		// Connectivity with API servers is blocked. No sense to make API requests
		return false, fmt.Errorf(`%s`, reasonDescription)
	}

	pub, priv, err := wireguard.GenerateKeys(m.wgToolBinPath)
	if err != nil {
		return false, err
	}

	isVPNConnected, connectedVpnType := m.service.ConnectedType()

	if !isVPNConnected || connectedVpnType != vpn.WireGuard {
		// use 'activePublicKey' ONLY if WireGuard is connected
		activePublicKey = ""
	}

	// trying to update WG keys with notifying API about current active public key (if it exists)
	localIP, err := m.api.WireGuardKeySet(session, pub, activePublicKey)
	if err != nil {
		if len(activePublicKey) == 0 {
			// IMPORTANT! As soon as server receive request with empty 'activePublicKey' - it clears all keys
			// Therefore, we have to ensure that local keys are not using anymore (we have to clear them independently from we received response or not)
			m.service.WireGuardSaveNewKeys("", "", "")
		}
		log.Info("WG keys not updated: ", err)

		var e types.APIError
		if errors.As(err, &e) {
			if e.ErrorCode == types.SessionNotFound {
				m.service.OnSessionNotFound()
				return false, fmt.Errorf("WG keys not updated (session not found)")
			}
		}
		return false, fmt.Errorf("WG keys not updated. Please check your internet connection")
	}

	log.Info(fmt.Sprintf("WG keys updated (%s:%s) ", localIP.String(), pub))

	// notify service about new keys
	m.service.WireGuardSaveNewKeys(pub, priv, localIP.String())

	if isRotationStopped {
		// If there was no public key defined - start keys rotation
		m.StartKeysRotation()
	}

	return true, nil
}

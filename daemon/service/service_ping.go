

//go:build !fastping
// +build !fastping
package service

import (
	"fmt"
	"time"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/api/types"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/ping"
)

// PingServers ping vpn servers.
// In some cases the multiple (and simultaneous pings) are leading to OS crash on macOS and Windows.
// It happens when installed some third-party 'security' software.
// Therefore, we using ping algorithm which avoids simultaneous pings and doing it one-by-one
func (s *Service) PingServers(retryCount int, timeoutMs int) (map[string]int, error) {

	if s._vpn != nil {
		log.Info("Servers pinging skipped due to connected state")
		return nil, nil
	}

	if timeoutMs <= 0 {
		log.Debug("Servers pinging skipped: timeout argument value is 0")
		return nil, nil
	}

	// TODO: (to discuss) do we need to block pinging when VPNAccess==blocked
	/*
		if isBlocked, reasonDescription, err := s.IsConnectivityBlocked(); err == nil && isBlocked {
			log.Info("Servers pinging skipped: ", reasonDescription)
			return nil, nil
		}
	*/

	timeoutTime := time.Now().Add(time.Millisecond * time.Duration(timeoutMs))

	var geoLocation *types.GeoLookupResponse = nil
	if timeoutMs >= 3000 {
		l, err := s._api.GeoLookup(1500)
		if err != nil {
			log.Warning("(pinging) unable to obtain geolocation (fastest server detection could be not accurate):", err)
		}
		geoLocation = l
	} else {
		log.Warning("(pinging) not enough time to check geolocation (fastest server detection could be not accurate)")
	}

	// get servers IP
	// IPs will be sorted by distance from current location (nearest - first)
	hosts, err := s.getHostsToPing(geoLocation)
	if err != nil {
		log.Info("Servers ping failed: " + err.Error())
		return nil, err
	}

	result := make(map[string]int)

	funcPingIteration := func(onePingTimeoutMs int, timeout *time.Time) map[string]int {

		retMap := make(map[string]int)

		i := 0
		for _, h := range hosts {
			if s._vpn != nil {
				log.Info("Servers pinging stopped due to connected state")
				break
			}
			if timeout != nil && time.Now().Add(time.Millisecond*time.Duration(onePingTimeoutMs)).After(*timeout) {
				log.Info("Servers pinging stopped due max-timeout for this operation")
				break
			}

			if h == nil {
				continue
			}
			ipStr := h.String()
			if len(ipStr) <= 0 {
				continue
			}

			pinger, err := ping.NewPinger(ipStr)
			if err != nil {
				log.Error("Pinger creation error: " + err.Error())
				continue
			}

			pinger.SetPrivileged(true)
			pinger.Count = 1
			pinger.Timeout = time.Millisecond * time.Duration(onePingTimeoutMs)
			pinger.Run()
			stat := pinger.Statistics()
			i++

			if stat.AvgRtt > 0 {
				retMap[ipStr] = int(stat.AvgRtt / time.Millisecond)
			}

			if timeout == nil && len(retMap) > 0 && len(retMap)%10 == 0 {
				// periodically notify ping results when pinging in background
				s._evtReceiver.OnPingStatus(retMap)
			}
		}

		log.Info(fmt.Sprintf("Pinged %d of %d servers (%d successfully, timeout=%d)", i, len(hosts), len(retMap), onePingTimeoutMs))
		return retMap
	}

	// do not allow multiple ping request simultaneously
	if s._isServersPingInProgress {
		log.Info("Servers pinging skipped. Ping already in progress")
		return nil, nil
	}
	s._isServersPingInProgress = true

	// OS-specific preparations (e.g. we need to add servers IPs to firewall exceptions list)
	if err := s.implPingServersStarting(hosts); err != nil {
		log.Error("implPingServersStarting failed: " + err.Error())
	}

	// First ping iteration. Doing it fast. 300ms max for each server
	result = funcPingIteration(300, &timeoutTime)

	// The first ping result already received.
	// So, now there is no rush to do second ping iteration. Doing it in background.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic in background ping: ", r)
				if err, ok := r.(error); ok {
					log.ErrorTrace(err)
				}
			}

			if err := s.implPingServersStopped(hosts); err != nil {
				log.Error("implPingServersStopped failed: " + err.Error())
			}

			s._isServersPingInProgress = false
		}()

		ret := funcPingIteration(1000, nil)
		for k, v := range ret {
			if v <= 0 {
				continue
			}
			if oldVal, ok := result[k]; ok {
				if v < oldVal {
					result[k] = v
				}
			} else {
				result[k] = v
			}
		}
		s._evtReceiver.OnPingStatus(result)
	}()

	// Return first ping result
	// This result may not contain results for all servers
	return result, nil
}

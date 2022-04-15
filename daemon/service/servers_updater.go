//
//  Daemon for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
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

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ivpn/desktop-app/daemon/api"
	"github.com/ivpn/desktop-app/daemon/api/types"
	"github.com/ivpn/desktop-app/daemon/service/platform"
	"github.com/ivpn/desktop-app/daemon/service/platform/filerights"
)

type serversUpdater struct {
	servers           *types.ServersInfoResponse
	api               *api.API
	updatedNotifyChan chan struct{}
	token             string
}

// CreateServersUpdater - constructor for serversUpdater object
func CreateServersUpdater(apiObj *api.API) (IServersUpdater, error) {
	updater := &serversUpdater{api: apiObj}

	updater.updatedNotifyChan = make(chan struct{}, 1)

	return updater, nil
}

func (updater *serversUpdater) startService(token string) {
	updater.token = token
	updater.GetServers()

	// update servers list in background
	if err := updater.startUpdater(); err != nil {
		log.Error("Failed to start servers-list updater: ", err)
	}
}

// GetServers - get servers list.
// Use cached data (if exists), otherwise - download servers list.
func (s *serversUpdater) GetServers() (*types.ServersInfoResponse, error) {
	if s.servers != nil {
		return s.servers, nil
	}

	servers, _, _, err := readServersFromCache()

	if servers != nil && err == nil {
		s.servers = servers
		return servers, nil
	}

	return s.updateServers()
}

func (s *serversUpdater) startUpdater() error {
	go func(s *serversUpdater) {
		isFirstIteration := true
		for {
			updateDelay := time.Hour
			if _, err := s.updateServers(); err != nil {
				log.Error(err)
				if isFirstIteration {
					// The first try to update can be failed because of daemon is starting on OS boot
					// There could be not all connectivity initialized
					// Therefore - trying in 5min later
					updateDelay = time.Minute * 5
				}
			}
			isFirstIteration = false
			time.Sleep(updateDelay)
		}
	}(s)

	return nil
}

// UpdateServers - download servers list
func (s *serversUpdater) updateServers() (*types.ServersInfoResponse, error) {
	servers, _, err := s.api.ServersList(s.token)
	if err != nil {
		return servers, fmt.Errorf("failed to download servers list: %w", err)
	}
	log.Info(fmt.Sprintf("Updated servers info (%d OpenVPN)\n", len(servers.Servers)))

	s.servers = servers
	if err := writeServersToCache(servers); err != nil {
		log.Error("failed to save servers cache file: ", err)
	}

	select {
	case s.updatedNotifyChan <- struct{}{}:
		// notified
	default:
		// channel is full
	}

	return servers, nil
}

// UpdateNotifierChannel returns channel which is notifying when servers was updated
func (s *serversUpdater) UpdateNotifierChannel() chan struct{} {
	return s.updatedNotifyChan
}

func readServersFromCache() (svrs *types.ServersInfoResponse, apiIPsV4 []string, apiIPsV6 []string, e error) {

	serversFile := platform.ServersFile()

	_, err := os.Stat(serversFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil, fmt.Errorf("failed to read servers cache file: %w", err)
		}
		return nil, nil, nil, fmt.Errorf("failed to info about servers cache file: %w", err)
	}

	data, err := ioutil.ReadFile(serversFile)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read servers cache file: %w", err)
	}

	servers := new(types.ServersInfoResponse)
	if err := json.Unmarshal(data, servers); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal servers cache file: %w", err)
	}

	// check servers.json file has correct access rights (can we use it's data?)
	if err := filerights.CheckFileAccessRightsConfig(serversFile); err != nil {
		os.Remove(serversFile)
		// we can not use servers info from this file
		// but we can try to get IP addresses of alternate IP's
		// It is safe, because we are checking TLS server name for "api.ivpn.net" when accessing API (https)
		return nil, nil, nil, fmt.Errorf("skip reading servers cache file: %w", err)
	}

	return servers, nil, nil, nil
}

func writeServersToCache(servers *types.ServersInfoResponse) error {
	if servers == nil {
		return errors.New("nothing to save. Servers is null")
	}

	data, err := json.Marshal(servers)
	if err != nil {
		return fmt.Errorf("failed to marshal servers into a cache: %w", err)
	}

	if data == nil {
		return errors.New("failed to serialize servers")
	}

	return ioutil.WriteFile(platform.ServersFile(), data, filerights.DefaultFilePermissionsForConfig())
}

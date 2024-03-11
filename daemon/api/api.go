//
//  Daemon for IVPN Client Desktop
//  https://github.com/tahirmahm123/vpn-desktop-app-daemon
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

package api

import (
	"encoding/json"
	"fmt"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/config"
	"net"
	"sync"
	"time"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/api/types"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
	protocolTypes "github.com/tahirmahm123/vpn-desktop-app/daemon/protocol/types"
)

// API URLs
const (
	_defaultRequestTimeout     = time.Second * 10 // full request time (for each request)
	_defaultDialTimeout        = time.Second * 5  // time for t
	_apiPathPrefix             = "v3"
	_updateHost                = "repo.ivpn.net"
	_sessionNewPath            = _apiPathPrefix + "/auth"
	_serversPath               = "v2/servers-list?group=country,protocol"
	_sessionStatusPath         = _apiPathPrefix + "/details"
	_sessionDeletePath         = _apiPathPrefix + "/signout"
	_wgKeySetPath              = _apiPathPrefix + "/wg-keys"
	_geoLookupPath             = "/location"
	_forceDeviceLogoutByIdPath = _apiPathPrefix + "/logout/"
	_forceAllDevicesLogoutPath = _apiPathPrefix + "/logout-all"
)

// Alias - alias description of API request (can be requested by UI client)
type Alias struct {
	host string
	path string
	// If isArcIndependent!=true, the path will be updated: the "_<architecture>" will be added to filename
	// (see 'DoRequestByAlias()' for details)
	// Example:
	//		The "updateInfo_macOS" on arm64 platform will use file "/macos/update_arm64.json" (NOT A "/macos/update.json")
	isArcIndependent bool
}

// APIAliases - aliases of API requests (can be requested by UI client)
// NOTE: the aliases bellow are only for amd64 architecture!!!
// If isArcIndependent!=true: Filename construction for non-amd64 architectures: filename_<architecture>.<extensions>
// (see 'DoRequestByAlias()' for details)
// Example:
//
//	The "updateInfo_macOS" on arm64 platform will use file "/macos/update_arm64.json" (NOT A "/macos/update.json")
const (
	GeoLookupApiAlias string = "geo-lookup"
)

var APIAliases = map[string]Alias{
	GeoLookupApiAlias: {host: config.GetAPIHost(), path: _geoLookupPath},

	"updateInfo_Linux":   {host: _updateHost, path: "/stable/_update_info/update.json"},
	"updateSign_Linux":   {host: _updateHost, path: "/stable/_update_info/update.json.sign.sha256.base64"},
	"updateInfo_macOS":   {host: _updateHost, path: "/macos/update.json"},
	"updateSign_macOS":   {host: _updateHost, path: "/macos/update.json.sign.sha256.base64"},
	"updateInfo_Windows": {host: _updateHost, path: "/windows/update.json"},
	"updateSign_Windows": {host: _updateHost, path: "/windows/update.json.sign.sha256.base64"},

	"updateInfo_manual_Linux":   {host: _updateHost, path: "/stable/_update_info/update_manual.json"},
	"updateSign_manual_Linux":   {host: _updateHost, path: "/stable/_update_info/update_manual.json.sign.sha256.base64"},
	"updateInfo_manual_macOS":   {host: _updateHost, path: "/macos/update_manual.json"},
	"updateSign_manual_macOS":   {host: _updateHost, path: "/macos/update_manual.json.sign.sha256.base64"},
	"updateInfo_manual_Windows": {host: _updateHost, path: "/windows/update_manual.json"},
	"updateSign_manual_Windows": {host: _updateHost, path: "/windows/update_manual.json.sign.sha256.base64"},

	"updateInfo_beta_Linux":   {host: _updateHost, path: "/stable/_update_info/update_beta.json"},
	"updateSign_beta_Linux":   {host: _updateHost, path: "/stable/_update_info/update_beta.json.sign.sha256.base64"},
	"updateInfo_beta_macOS":   {host: _updateHost, path: "/macos/update_beta.json"},
	"updateSign_beta_macOS":   {host: _updateHost, path: "/macos/update_beta.json.sign.sha256.base64"},
	"updateInfo_beta_Windows": {host: _updateHost, path: "/windows/update_beta.json"},
	"updateSign_beta_Windows": {host: _updateHost, path: "/windows/update_beta.json.sign.sha256.base64"},
}

var log *logger.Logger

func init() {
	log = logger.NewLogger("api")
}

// IConnectivityInfo information about connectivity
type IConnectivityInfo interface {
	// IsConnectivityBlocked - returns nil if connectivity NOT blocked
	IsConnectivityBlocked() (err error)
}

type geolookup struct {
	mutex     sync.Mutex
	isRunning bool
	done      chan struct{}

	location types.GeoLookupResponse
	response []byte
	err      error
}

// API contains data about IVPN API servers
type API struct {
	mutex                 sync.Mutex
	alternateIPsV4        []net.IP
	lastGoodAlternateIPv4 net.IP
	alternateIPsV6        []net.IP
	lastGoodAlternateIPv6 net.IP
	connectivityChecker   IConnectivityInfo

	// last geolookups result
	geolookupV4 geolookup
	geolookupV6 geolookup
}

// CreateAPI creates new API object
func CreateAPI() (*API, error) {
	return &API{}, nil
}

func (a *API) SetConnectivityChecker(connectivityChecker IConnectivityInfo) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.connectivityChecker = connectivityChecker
}

// DownloadServersList - download servers list form API IVPN server
func (a *API) DownloadServersList() (*types.ServerListResponse, error) {
	//servers := new(types.ServerListResponse)
	servers, _, _, err := a.ServersList()
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// DoRequestByAlias do API request (by API endpoint alias). Returns raw data of response
func (a *API) DoRequestByAlias(apiAlias string, ipTypeRequired protocolTypes.RequiredIPProtocol) (responseData []byte, err error) {
	// For geolookup requests we have specific function
	if apiAlias == GeoLookupApiAlias {
		if ipTypeRequired != protocolTypes.IPv4 && ipTypeRequired != protocolTypes.IPv6 {
			return nil, fmt.Errorf("geolookup request failed: IP version not defined")
		}
		_, responseData, err = a.GeoLookup(0, ipTypeRequired)
		return responseData, err
	}

	//// get connection info by API alias
	//alias, ok := APIAliases[apiAlias]
	//if !ok {
	//	return nil, fmt.Errorf("unexpected request alias")
	//}
	//
	//if !alias.isArcIndependent {
	//	// If isArcIndependent!=true, the path will be updated: the "_<architecture>" will be added to filename
	//	// Example:
	//	//		The "updateInfo_macOS" on arm64 platform will use file "/macos/update_arm64.json" (NOT A "/macos/update.json"!)
	//	if runtime.GOARCH != "amd64" {
	//		extIdx := strings.Index(alias.path, ".")
	//		if extIdx > 0 {
	//			newPath := alias.path[:extIdx] + "_" + runtime.GOARCH + alias.path[extIdx:]
	//			alias.path = newPath
	//		}
	//	}
	//}
	//
	//return a.requestRaw(ipTypeRequired, alias.host, alias.path, "", "", nil, 0, 0)
	return nil, err
}

// SessionNew - try to register new session
//func (a *API) SessionNew(accountID string, wgPublicKey string, kemKeys types.KemPublicKeys, forceLogin bool, captchaID string, captcha string, confirmation2FA string) (
//	*types.SessionNewResponse,
//	*types.SessionNewErrorLimitResponse,
//	*types.APIErrorResponse,
//	string, // RAW response
//	error) {
//
//	var successResp types.SessionNewResponse
//	var errorLimitResp types.SessionNewErrorLimitResponse
//	var apiErr types.APIErrorResponse
//
//	rawResponse := ""
//
//	request := &types.SessionNewRequest{
//		AccountID:       accountID,
//		PublicKey:       wgPublicKey,
//		KemPublicKeys:   kemKeys,
//		ForceLogin:      forceLogin,
//		CaptchaID:       captchaID,
//		Captcha:         captcha,
//		Confirmation2FA: confirmation2FA}
//
//	data, err := a.requestRaw(protocolTypes.IPvAny, "", _sessionNewPath, "POST", "application/json", request, 0, 0)
//	if err != nil {
//		return nil, nil, nil, rawResponse, err
//	}
//
//	rawResponse = string(data)
//
//	// Check is it API error
//	if err := json.Unmarshal(data, &apiErr); err != nil {
//		return nil, nil, nil, rawResponse, fmt.Errorf("failed to deserialize API response: %w", err)
//	}
//
//	// success
//	if apiErr.Status == types.CodeSuccess {
//		if err := json.Unmarshal(data, &successResp); err != nil {
//			return nil, nil, nil, rawResponse, fmt.Errorf("failed to deserialize API response: %w", err)
//		}
//
//		return &successResp, nil, &apiErr, rawResponse, nil
//	}
//
//	// Session limit check
//	if apiErr.Status == types.CodeSessionsLimitReached {
//		if err := json.Unmarshal(data, &errorLimitResp); err != nil {
//			return nil, nil, nil, rawResponse, fmt.Errorf("failed to deserialize API response: %w", err)
//		}
//		return nil, &errorLimitResp, &apiErr, rawResponse, types.CreateAPIError(apiErr.Status, apiErr.Message)
//	}
//
//	return nil, nil, &apiErr, rawResponse, types.CreateAPIError(apiErr.Status, apiErr.Message)
//}

// SessionStatus - get session status
//func (a *API) SessionStatus(session string) (
//	*types.ServiceStatusAPIResp,
//	*types.APIErrorResponse,
//	error) {
//
//	var resp types.SessionStatusResponse
//	var apiErr types.APIErrorResponse
//
//	request := &types.SessionStatusRequest{Session: session}
//
//	data, err := a.requestRaw(protocolTypes.IPvAny, "", _sessionStatusPath, "POST", "application/json", request, 0, 0)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	// Check is it API error
//	if err := json.Unmarshal(data, &apiErr); err != nil {
//		return nil, nil, fmt.Errorf("failed to deserialize API response: %w", err)
//	}
//
//	// success
//	if apiErr.Status == types.CodeSuccess {
//		if err := json.Unmarshal(data, &resp); err != nil {
//			return nil, nil, fmt.Errorf("failed to deserialize API response: %w", err)
//		}
//		return &resp.ServiceStatus, &apiErr, nil
//	}
//
//	return nil, &apiErr, types.CreateAPIError(apiErr.Status, apiErr.Message)
//}

// SessionDelete - remove session
//func (a *API) SessionDelete(session string) error {
//	request := &types.SessionDeleteRequest{Session: session}
//	resp := &types.APIErrorResponse{}
//	if err := a.request("", _sessionDeletePath, "POST", "application/json", request, resp); err != nil {
//		return err
//	}
//	if resp.Status != types.CodeSuccess {
//		return types.CreateAPIError(resp.Status, resp.Message)
//	}
//	return nil
//}

// WireGuardKeySet - update WG key
func (a *API) WireGuardKeySet(session string, publicKey string) (
	successResp *types.WGKeysUpdateResponse,
	rawResponse string, // RAW response
	err error) {

	data, statusCode, err := a.requestRaw(_wgKeySetPath, "POST", types.WGKeyUpdateRequest{
		PublicKey: publicKey,
	}, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + session,
	})
	if err != nil {
		//fmt.Printf("Error from Server %s", err)
		return nil, rawResponse, err
	}

	rawResponse = string(data)

	// success
	if statusCode == 200 {
		if err := json.Unmarshal(data, &successResp); err != nil {
			return nil, rawResponse, fmt.Errorf("failed to deserialize API response Session New API Success: %w", err)
		}
		return successResp, rawResponse, nil
	}
	return nil, rawResponse, fmt.Errorf("request Failed with Status coode %d and Response: %s", statusCode, rawResponse)
}

// GeoLookup gets geolocation
func (a *API) GeoLookup(timeoutMs int, ipTypeRequired protocolTypes.RequiredIPProtocol) (location *types.GeoLookupResponse, rawData []byte, retErr error) {
	// There could be multiple Geolookup requests at the same time.
	// It doesn't make sense to make multiple requests to the API.
	// The internal function below reduces the number of similar API calls.
	// TODO : fix api location
	return nil, nil, nil
	//singletonFunc := func(ipType protocolTypes.RequiredIPProtocol) (*types.GeoLookupResponse, []byte, error) {
	//	// Each IP protocol has separate request
	//	var gl *geolookup
	//	if ipType == protocolTypes.IPv4 {
	//		gl = &a.geolookupV4
	//	} else if ipType == protocolTypes.IPv6 {
	//		gl = &a.geolookupV6
	//	} else {
	//		return nil, nil, fmt.Errorf("geolookup request failed: IP version not defined")
	//	}
	//	// Try to make API request (if not started yet). Only one API request allowed in the same time.
	//	func() {
	//		gl.mutex.Lock()
	//		defer gl.mutex.Unlock()
	//		// if API call is already running - do nosing, just wait for results
	//		if gl.isRunning {
	//			return
	//		}
	//		// mark: call is already running
	//		gl.isRunning = true
	//		gl.done = make(chan struct{})
	//		// do API call in routine
	//		go func() {
	//			defer func() {
	//				// API call finished
	//				gl.isRunning = false
	//				close(gl.done)
	//			}()
	//			gl.response, gl.err = a.requestRaw(ipType, "", _geoLookupPath, "GET", "", nil, timeoutMs, 0)
	//			if err := json.Unmarshal(gl.response, &gl.location); err != nil {
	//				gl.err = fmt.Errorf("failed to deserialize API response: %w", err)
	//			}
	//		}()
	//	}()
	//	// wait for API call result (for routine stop)
	//	<-gl.done
	//	return &gl.location, gl.response, gl.err
	//}
	//
	//// request Geolocation info
	//if ipTypeRequired != protocolTypes.IPvAny {
	//	location, rawData, retErr = singletonFunc(ipTypeRequired)
	//} else {
	//	location, rawData, retErr = singletonFunc(protocolTypes.IPv4)
	//	if retErr != nil {
	//		location, rawData, retErr = singletonFunc(protocolTypes.IPv6)
	//	}
	//}
	//
	//if retErr != nil {
	//	return nil, nil, retErr
	//}
	//return location, rawData, nil
}
func (a *API) ServersList() (
	successResp *types.ServerListResponse,
	statusCode int,
	rawResponse string, // RAW response
	err error) {

	data, statusCode, err := a.requestRaw(_serversPath, "GET", nil, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": config.GetAuthCredentials(),
	})
	if err != nil {
		//fmt.Printf("Error from Server %s", err)
		return nil, 0, rawResponse, err
	}
	rawResponse = string(data)
	if statusCode == 401 {
		return nil, statusCode, rawResponse, fmt.Errorf("Unauthenticated")
	}
	// success
	if statusCode == 200 {
		if err := json.Unmarshal(data, &successResp); err != nil {
			return nil, statusCode, rawResponse, fmt.Errorf("failed to deserialize API response Session New API Success: %w", err)
		}
		return successResp, statusCode, rawResponse, nil
	}
	return nil, statusCode, rawResponse, fmt.Errorf("request Failed with Status coode %d and Response: %s", statusCode, rawResponse)
}

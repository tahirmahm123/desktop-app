
package api

import (
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/api/types"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
	protocolTypes "github.com/tahirmahm123/vpn-desktop-app/daemon/protocol/types"
)

// API URLs
const (
	_defaultRequestTimeout = time.Second * 20 // full request time (for each request)
	_defaultDialTimeout    = time.Second * 5  // time for the dial to the API server (for each request)
	_apiHost               = "api.vpnpacket.com"
	// _apiHost           = "api.vpn.net"
	_updateHost        = "repo.vpn.net"
	_serversPath       = "v2/servers-list"
	_apiPathPrefix     = "v2"
	_sessionNewPath    = _apiPathPrefix + "/auth"
	_sessionStatusPath = _apiPathPrefix + "/details"
	_sessionDeletePath = _apiPathPrefix + "/signout"
	_wgKeySetPath      = _apiPathPrefix + "/session/wg/set"
	_geoLookupPath     = _apiPathPrefix + "/geo-lookup"
)

// Alias - alias description of API request (can be requested by UI client)
type Alias struct {
	host string
	path string
	// If isArchDependent==true, the path will be updated: the "_<architecture>" will be added to filename
	// (see 'DoRequestByAlias()' for details)
	// Example:
	//		The "updateInfo_macOS" on arm64 platform will use file "/macos/update_arm64.json" (NOT A "/macos/update.json")
	isArchDependent bool
}

// APIAliases - aliases of API requests (can be requested by UI client)
// NOTE: the aliases bellow are only for amd64 architecture!!!
// If isArchDependent==true: Filename construction for non-amd64 architectures: filename_<architecture>.<extensions>
// (see 'DoRequestByAlias()' for details)
// Example:
//		The "updateInfo_macOS" on arm64 platform will use file "/macos/update_arm64.json" (NOT A "/macos/update.json")

var APIAliases = map[string]Alias{
	"geo-lookup": {host: _apiHost, path: _geoLookupPath},

	"updateInfo_Linux":   {host: _updateHost, isArchDependent: true, path: "/stable/_update_info/update.json"},
	"updateSign_Linux":   {host: _updateHost, isArchDependent: true, path: "/stable/_update_info/update.json.sign.sha256.base64"},
	"updateInfo_macOS":   {host: _updateHost, isArchDependent: true, path: "/macos/update.json"},
	"updateSign_macOS":   {host: _updateHost, isArchDependent: true, path: "/macos/update.json.sign.sha256.base64"},
	"updateInfo_Windows": {host: _updateHost, isArchDependent: true, path: "/windows/update.json"},
	"updateSign_Windows": {host: _updateHost, isArchDependent: true, path: "/windows/update.json.sign.sha256.base64"},

	"updateInfo_manual_Linux":   {host: _updateHost, isArchDependent: true, path: "/stable/_update_info/update_manual.json"},
	"updateSign_manual_Linux":   {host: _updateHost, isArchDependent: true, path: "/stable/_update_info/update_manual.json.sign.sha256.base64"},
	"updateInfo_manual_macOS":   {host: _updateHost, isArchDependent: true, path: "/macos/update_manual.json"},
	"updateSign_manual_macOS":   {host: _updateHost, isArchDependent: true, path: "/macos/update_manual.json.sign.sha256.base64"},
	"updateInfo_manual_Windows": {host: _updateHost, isArchDependent: true, path: "/windows/update_manual.json"},
	"updateSign_manual_Windows": {host: _updateHost, isArchDependent: true, path: "/windows/update_manual.json.sign.sha256.base64"},
}

var log *logger.Logger

func init() {
	log = logger.NewLogger("api")
}

// IConnectivityInfo information about connectivity
type IConnectivityInfo interface {
	IsConnectivityBlocked() (isBlocked bool, reasonDescription string, err error)
}

// API contains data about VPNrvers
type API struct {
	mutex                 sync.Mutex
	alternateIPsV4        []net.IP
	lastGoodAlternateIPv4 net.IP
	alternateIPsV6        []net.IP
	lastGoodAlternateIPv6 net.IP
	connectivityChecker   IConnectivityInfo
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

// IsAlternateIPsInitialized - checks if the alternate IP initialized
func (a *API) IsAlternateIPsInitialized(IPv6 bool) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if IPv6 {
		return len(a.alternateIPsV6) > 0
	}
	return len(a.alternateIPsV4) > 0
}

func (a *API) GetLastGoodAlternateIP(IPv6 bool) net.IP {
	if IPv6 {
		return a.lastGoodAlternateIPv6
	}
	return a.lastGoodAlternateIPv4
}

func (a *API) SetLastGoodAlternateIP(IPv6 bool, ip net.IP) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if IPv6 {
		a.lastGoodAlternateIPv6 = ip
	}
	a.lastGoodAlternateIPv4 = ip
}

func (a *API) getAlternateIPs(IPv6 bool) []net.IP {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if IPv6 {
		return a.alternateIPsV6
	}
	return a.alternateIPsV4
}

// SetAlternateIPs save info about alternate servers IP addresses
func (a *API) SetAlternateIPs(IPv4List []string, IPv6List []string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.doSetAlternateIPs(false, IPv4List)
	a.doSetAlternateIPs(true, IPv6List)
	return nil
}

func (a *API) doSetAlternateIPs(IPv6 bool, IPs []string) error {
	if len(IPs) == 0 {
		log.Warning("Unable to set alternate API IP list. List is empty")
	}

	lastGoodIP := a.GetLastGoodAlternateIP(IPv6)

	ipList := make([]net.IP, 0, len(IPs))

	isLastIPExists := false
	for _, ipStr := range IPs {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}

		ipList = append(ipList, ip)

		if ip.Equal(lastGoodIP) {
			isLastIPExists = true
		}
	}

	if !isLastIPExists {
		if IPv6 {
			a.lastGoodAlternateIPv6 = nil
		} else {
			a.lastGoodAlternateIPv4 = nil
		}
	}

	// set new alternate IP list
	if IPv6 {
		a.alternateIPsV6 = ipList
	} else {
		a.alternateIPsV4 = ipList
	}

	return nil
}

// DoRequestByAlias do API request (by API endpoint alias). Returns raw data of response
func (a *API) DoRequestByAlias(apiAlias string, ipTypeRequired protocolTypes.RequiredIPProtocol) (responseData []byte, err error) {
	alias, ok := APIAliases[apiAlias]
	if !ok {
		return nil, fmt.Errorf("unexpected request alias")
	}

	if alias.isArchDependent {
		// If isArchDependent==true, the path will be updated: the "_<architecture>" will be added to filename
		// Example:
		//		The "updateInfo_macOS" on arm64 platform will use file "/macos/update_arm64.json" (NOT A "/macos/update.json"!)
		if runtime.GOARCH != "amd64" {
			extIdx := strings.Index(alias.path, ".")
			if extIdx > 0 {
				newPath := alias.path[:extIdx] + "_" + runtime.GOARCH + alias.path[extIdx:]
				alias.path = newPath
			}
		}
	}

	retData, retErr := a.requestRaw(alias.path, "", "", nil, 0, 0)

	return retData, retErr
}

// SessionNew - try to register new session
func (a *API) SessionNew(username string, password string) (
	*types.SessionNewResponse,
	*types.SessionNewErrorLimitResponse,
	*types.APIErrorResponse,
	string, // RAW response
	error) {

	var successResp types.SessionNewResponse
	var errorLimitResp types.SessionNewErrorLimitResponse
	var apiErr types.APIErrorResponse

	rawResponse := ""

	request := &types.SessionNewRequest{Username: username, Password: password}

	data, err := a.requestRaw(_sessionNewPath, "POST", "application/json", request, 0, 0)
	if err != nil {
		//fmt.Printf("Error from Server %s", err)
		return nil, nil, nil, rawResponse, err
	}

	rawResponse = string(data)
	//fmt.Printf("Data from Server %s", rawResponse)
	// Check is it API error
	if err := json.Unmarshal(data, &apiErr); err != nil {
		return nil, nil, nil, rawResponse, fmt.Errorf("failed to deserialize API response in Session New API Error: %w", err)
	}

	// success
	if apiErr.Status == types.CodeSuccess {
		if err := json.Unmarshal(data, &successResp); err != nil {
			return nil, nil, nil, rawResponse, fmt.Errorf("failed to deserialize API response Session New API Success: %w", err)
		}
		return &successResp, nil, &apiErr, rawResponse, nil
	}

	// Session limit check
	if apiErr.Status == types.CodeSessionsLimitReached {
		if err := json.Unmarshal(data, &errorLimitResp); err != nil {
			return nil, nil, nil, rawResponse, fmt.Errorf("failed to deserialize API response: %w", err)
		}
		return nil, &errorLimitResp, &apiErr, rawResponse, types.CreateAPIError(apiErr.Status, apiErr.Message)
	}

	return nil, nil, &apiErr, rawResponse, types.CreateAPIError(apiErr.Status, apiErr.Message)
}

// SessionStatus - get session status
func (a *API) SessionStatus(session string) (
	*types.ServiceStatusAPIResp,
	*types.APIErrorResponse,
	error) {

	var resp types.SessionStatusResponse
	var apiErr types.APIErrorResponse

	request := &types.SessionStatusRequest{Session: session}

	data, err := a.requestRaw(_sessionStatusPath, "POST", "application/json", request, 0, 0)
	if err != nil {
		return nil, nil, err
	}
	// Check is it API error
	if err := json.Unmarshal(data, &apiErr); err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize API response: %w", err)
	}

	// success
	if apiErr.Status == types.CodeSuccess {
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, nil, fmt.Errorf("failed to deserialize API response: %w", err)
		}
		return &resp.ServiceStatus, &apiErr, nil
	}

	return nil, &apiErr, types.CreateAPIError(apiErr.Status, apiErr.Message)
}

// SessionDelete - remove session
func (a *API) SessionDelete(session string) error {
	request := &types.SessionDeleteRequest{Session: session}
	resp, err := a.requestRaw(_sessionDeletePath, "POST", "application/json", request, 0, 0)
	if err != nil {
		return err
	}
	if string(resp) == "" {
		return nil
	}
	return err
}

// WireGuardKeySet - update WG key
func (a *API) WireGuardKeySet(session string, newPublicWgKey string, activePublicWgKey string) (localIP net.IP, err error) {
	request := &types.SessionWireGuardKeySetRequest{
		Session:            session,
		PublicKey:          newPublicWgKey,
		ConnectedPublicKey: activePublicWgKey}

	resp := &types.SessionsWireGuardResponse{}

	if err := a.request("", _wgKeySetPath, "POST", "application/json", request, resp); err != nil {
		return nil, err
	}

	if resp.Status != types.CodeSuccess {
		return nil, types.CreateAPIError(resp.Status, resp.Message)
	}

	localIP = net.ParseIP(resp.IPAddress)
	if localIP == nil {
		return nil, fmt.Errorf("failed to set WG key (failed to parse local IP in API response)")
	}

	return localIP, nil
}

// GeoLookup get geolocation
func (a *API) GeoLookup(timeoutMs int) (location *types.GeoLookupResponse, err error) {
	resp := &types.GeoLookupResponse{}

	if err := a.requestEx(_geoLookupPath, "GET", "", nil, resp, timeoutMs, 0); err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *API) ServersList(token string) (servers *types.ServersInfoResponse, rawResponse string, err error) {
	resp := &types.ServersInfoResponse{}
	req := &types.ServersListRequest{Session: token}
	response, err := a.requestRaw(_serversPath, "GET", "application/json", req, 0, 0)
	if err != nil {
		return nil, "", err
	}

	if err := json.Unmarshal(response, resp); err != nil {
		return nil, "", fmt.Errorf("failed to deserialize API response: %w", err)
	}

	return resp, string(response), nil
}

package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/protocol/types"
)

func getURL(host string, urlpath string) string {
	return "https://" + path.Join(host, urlpath)
}

func newRequest(urlPath string, method string, contentType string, body io.Reader) (*http.Request, error) {
	if len(method) == 0 {
		method = "GET"
	}

	req, err := http.NewRequest(method, urlPath, body)
	if err != nil {
		return nil, err
	}

	if len(contentType) > 0 {
		req.Header.Add("Content-type", contentType)
	}

	return req, nil
}

func (a *API) doRequest(urlPath string, method string, contentType string, request interface{}, timeoutMs int, timeoutDialMs int) (resp *http.Response, err error) {
	return a.doRequestAPIHost(urlPath, method, contentType, request, timeoutMs, timeoutDialMs)
}

func (a *API) doRequestAPIHost(urlPath string, method string, contentType string, request interface{}, timeoutMs int, timeoutDialMs int) (resp *http.Response, err error) {
	// isIPv6 := ipTypeRequired == types.IPv6

	// timeout time for full request
	timeout := _defaultRequestTimeout
	if timeoutMs > 0 {
		timeout = time.Millisecond * time.Duration(timeoutMs)
	}
	// timeout for the dial
	timeoutDial := _defaultDialTimeout
	if timeoutDialMs > 0 {
		timeoutDial = time.Millisecond * time.Duration(timeoutDialMs)
	}
	if timeoutDial > timeout {
		timeoutDial = 0
	}

	// When trying to access API server by alternate IPs (not by DNS name)
	// we need to configure TLS to use api.vpn.net hostname
	// (to avoid certificate errors)
	transCfg := &http.Transport{
		// NOTE: TLSClientConfig not in use in case of DialTLS defined
		TLSClientConfig: &tls.Config{
			ServerName: _apiHost,
		},

		// using certificate key pinning
		// DialTLS: makeDialer(APIVpnHashes, true, _apiHost, timeoutDial),
	}

	// configure http-client with preconfigured TLS transport
	client := &http.Client{Transport: transCfg, Timeout: timeout}

	data := []byte{}
	if request != nil {
		data, err = json.Marshal(request)
		if err != nil {
			return nil, err
		}
	}
	var tokenCheck types.CheckToken
	json.Unmarshal(data, &tokenCheck)
	bodyBuffer := bytes.NewBuffer(data)
	log.Info(fmt.Sprintf("Path: %s Data: %s", urlPath, string(data)))

	// try to access API server by host DNS
	var firstResp *http.Response
	var firstErr error
	// if isCanUseDNS {
	req, err := newRequest(getURL(_apiHost, urlPath), method, contentType, bodyBuffer)
	if tokenCheck.Session != "" {
		req.Header.Add("Authorization", "Bearer "+tokenCheck.Session)
	}
	if err != nil {
		return nil, err
	}

	firstResp, firstErr = client.Do(req)
	if firstErr == nil {
		return firstResp, firstErr
	}
	log.Warning("Failed to access " + _apiHost)

	return nil, fmt.Errorf("unable to access VPNrver: %w", firstErr)
}

func (a *API) requestRaw(urlPath string, method string, contentType string, requestObject interface{}, timeoutMs int, timeoutDialMs int) (responseData []byte, err error) {
	resp, err := a.doRequest(urlPath, method, contentType, requestObject, timeoutMs, timeoutDialMs)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	newStr := buf.String()
	if err != nil {
		return nil, fmt.Errorf("failed to get API HTTP response body: %w", err)
	}
	log.Info(fmt.Sprintf("API HTTP Path: %s Body: %s", urlPath, newStr))
	return buf.Bytes(), nil
}

func (a *API) request(host string, urlPath string, method string, contentType string, requestObject interface{}, responseObject interface{}) error {
	fmt.Println("Host " + host)
	return a.requestEx(urlPath, method, contentType, requestObject, responseObject, 0, 0)
}

func (a *API) requestEx(urlPath string, method string, contentType string, requestObject interface{}, responseObject interface{}, timeoutMs int, timeoutDialMs int) error {
	body, err := a.requestRaw(urlPath, method, contentType, requestObject, timeoutMs, timeoutDialMs)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, responseObject); err != nil {
		return fmt.Errorf("failed to deserialize API response: %w", err)
	}
	return nil
}

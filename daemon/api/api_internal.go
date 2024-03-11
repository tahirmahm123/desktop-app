package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/config"
	"io"
	"net/http"
	"path"
)

func getURL(urlPath string) string {
	return "https://" + path.Join(config.GetAPIHost(), urlPath)
}

func newRequest(urlPath string, method string, headers map[string]string, body io.Reader) (*http.Request, error) {
	if len(method) == 0 {
		method = "GET"
	}

	req, err := http.NewRequest(method, urlPath, body)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

func (a *API) doRequest(urlPath string, requestType string, request interface{}, headers map[string]string) (resp *http.Response, err error) {

	client := http.Client{}
	var data []byte
	if request != nil {
		data, err = json.Marshal(request)
		if err != nil {
			return nil, err
		}
	}

	bodyBuffer := bytes.NewBuffer(data)
	req, err := newRequest(urlPath, requestType, headers, bodyBuffer)
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (a *API) requestRaw(urlPath string, method string, requestObject interface{}, headers map[string]string) (responseData []byte, statusCode int, err error) {
	if headers == nil {
		headers = map[string]string{
			"Content-Type": "application/json",
		}
	}
	resp, err := a.doRequest(getURL(urlPath), method, requestObject, headers)
	if err != nil {
		return nil, 0, fmt.Errorf("API request failed: %w", err)
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	newStr := buf.String()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get API HTTP response body: %w", err)
	}
	log.Info(fmt.Sprintf("API HTTP Path: %s Status Code: %d Body: %s", urlPath, resp.StatusCode, newStr))
	return buf.Bytes(), resp.StatusCode, nil
}

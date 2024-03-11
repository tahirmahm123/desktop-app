package config

import (
	"encoding/base64"
	"fmt"
	"os"
)

const (
	appName     = "GoldenGuardVPN"
	app         = "Golden Guard VPN"
	apiURL      = "api.goldenguardvpn.com"
	apiUsername = "api-user@goldenguardvpn.com"
	apiPassword = "Fdma#eEe@)Td%9yrNu"
)

func GetAPIHost() string {
	return apiURL
}
func GetAuthCredentials() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(apiUsername+":"+apiPassword))
}
func GetAppName() string {
	return appName
}
func GetName() string {
	return app
}
func GetAppRoot() string {
	return fmt.Sprintf(appPath, GetName())
}
func OsName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
	}
	return hostname
}
func OsType() string {
	return osType
}

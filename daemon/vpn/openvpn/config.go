
package openvpn

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/netinfo"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/platform"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/platform/filerights"
)

// ConnectionParams represents OpenVPN connection parameters
type ConnectionParams struct {
	username          string
	password          string
	multihopExitSrvID string
	tcp               bool
	hostPort          int
	hostIP            net.IP
	proxyType         string
	proxyAddress      net.IP
	proxyPort         int
	proxyUsername     string
	proxyPassword     string
}
type Certificate struct {
	Certificate string `json:"certificate"`
}

// SetCredentials update WG credentials
func (c *ConnectionParams) SetCredentials(username, password string) {
	c.password = password
	c.username = username

	// MultiHop configuration is based just by adding "@exit_server_id" to the end of username
	// And forwarding this info on server
	if len(c.multihopExitSrvID) > 0 {
		c.username = fmt.Sprintf("%s@%s", username, c.multihopExitSrvID)
	}
}

// CreateConnectionParams creates OpenVPN connection parameters object
func CreateConnectionParams(
	multihopExitSrvID string,
	tcp bool,
	hostPort int,
	hostIP net.IP,
	proxyType string,
	proxyAddress net.IP,
	proxyPort int,
	proxyUsername string,
	proxyPassword string) ConnectionParams {

	return ConnectionParams{
		multihopExitSrvID: multihopExitSrvID,
		tcp:               tcp,
		hostPort:          hostPort,
		hostIP:            hostIP,
		proxyType:         proxyType,
		proxyAddress:      proxyAddress,
		proxyPort:         proxyPort,
		proxyUsername:     proxyUsername,
		proxyPassword:     proxyPassword}
}

// WriteConfigFile saves OpenVPN connection parameters into a config file
func (c *ConnectionParams) WriteConfigFile(
	localPort int,
	filePathToSave string,
	miAddr string,
	miPort int,
	logFile string,
	obfsproxyPort int,
	openVpnCertificate string,
	isCanUseV24Params bool) error {

	configText, err := c.generateConfiguration(openVpnCertificate, miAddr, miPort)
	if err != nil {
		return fmt.Errorf("failed to generate openvpn configuration : %w", err)
	}

	err = ioutil.WriteFile(filePathToSave, []byte(configText), 0600) // read\write only for privileged user
	if err != nil {
		return fmt.Errorf("failed to save OpenVPN configuration into a file: %w", err)
	}

	// only for Windows: Golang is not able to change file permissins in Windows style
	if err := filerights.WindowsChmod(filePathToSave, 0600); err != nil { // read\write only for privileged user
		return fmt.Errorf("failed to change OpenVPN configuration file permissions: %w", err)
	}

	log.Info("Configuring OpenVPN...\n",
		"=====================\n",
		configText,
		"\n=====================\n")

	return nil
}
func (c *ConnectionParams) generateConfiguration(openVpnCertificate string, miAddr string, miPort int) (cfg string, err error) {
	proto := "udp"
	if c.tcp {
		proto = "tcp"
	}

	cfgArr := make([]string, 0, 32)
	cfgArr = append(cfgArr, fmt.Sprintf("management %s %d", miAddr, miPort))
	cfgArr = append(cfgArr, "management-client")

	cfgArr = append(cfgArr, "management-hold")
	cfgArr = append(cfgArr, "auth-user-pass")
	cfgArr = append(cfgArr, "auth-nocache")

	cfgArr = append(cfgArr, "management-query-passwords")

	cfgArr = append(cfgArr, "management-signal")
	cfgArr = append(cfgArr, fmt.Sprintf("remote %s %d %s", c.hostIP, c.hostPort, proto))
	return strings.Replace(openVpnCertificate, "[REMOTES]", strings.Join(cfgArr, "\n"), 1), nil
}

func (c *ConnectionParams) generateConfigurationx(
	localPort int,
	miAddr string,
	miPort int,
	logFile string,
	obfsproxyPort int,
	openVpnCertificate string,
	isCanUseV24Params bool) (cfg []string, err error) {

	if obfsproxyPort > 0 {
		c.tcp = true
		c.hostPort = platform.ObfsproxyHostPort()
		c.proxyType = "socks"
		c.proxyAddress = net.IPv4(127, 0, 0, 1) // "127.0.0.1"
		c.proxyPort = obfsproxyPort
		c.proxyUsername = ""
		c.proxyPassword = ""
	}

	cfg = make([]string, 0, 32)

	cfg = append(cfg, "client")
	cfg = append(cfg, fmt.Sprintf("management %s %d", miAddr, miPort))
	cfg = append(cfg, "management-client")

	cfg = append(cfg, "management-hold")
	cfg = append(cfg, "auth-user-pass")
	cfg = append(cfg, "auth-nocache")

	cfg = append(cfg, "management-query-passwords")

	cfg = append(cfg, "management-signal")

	// Handshake Window --the TLS - based key exchange must finalize within n seconds of handshake initiation by any peer(default = 60 seconds).
	// If the handshake fails openvpn will attempt to reset our connection with our peer and try again.
	cfg = append(cfg, "hand-window 6")

	if isCanUseV24Params {
		cfg = append(cfg, "compress")
		cfg = append(cfg, "pull-filter ignore \"ping\"")
	} else {
		cfg = append(cfg, "comp-lzo no")
	}

	// To change default connection-check time:
	// 	pull-filter ignore "ping"
	//	keepalive 8 30
	cfg = append(cfg, "keepalive 8 30")

	// proxy
	if c.proxyType == "http" || c.proxyType == "socks" {

		// proxy authentication
		proxyAuthFile := ""
		if c.proxyUsername != "" && c.proxyPassword != "" {
			proxyAuthFile = "\"" + platform.OpenvpnProxyAuthFile() + "\""
			err := ioutil.WriteFile(platform.OpenvpnProxyAuthFile(), []byte(fmt.Sprintf("%s\n%s", c.proxyUsername, c.proxyPassword)), 0600)
			if err != nil {
				log.Error(err)
				return nil, fmt.Errorf("Failed to save file with proxy credentials: %w", err)
			}
		}

		// proxy config
		switch c.proxyType {
		case "http":
			cfg = append(cfg, "http-proxy-retry")
			cfg = append(cfg, fmt.Sprintf("http-proxy %s %d %s", c.proxyAddress.String(), c.proxyPort, proxyAuthFile))
			break
		case "socks":
			cfg = append(cfg, "socks-proxy-retry")
			cfg = append(cfg, fmt.Sprintf("socks-proxy %s %d %s", c.proxyAddress.String(), c.proxyPort, proxyAuthFile))
			break
		}
	}

	if len(logFile) > 0 && logger.IsEnabled() {
		cfg = append(cfg, fmt.Sprintf(`log "%s"`, logFile))
	}

	cfg = append(cfg, "dev tun")

	if c.tcp {
		cfg = append(cfg, "proto tcp")
	} else {
		cfg = append(cfg, "proto udp")
	}

	if c.hostIP.IsUnspecified() {
		return nil, errors.New("unable to connect. Host IP not defined")
	}
	if c.hostPort < 0 || c.hostPort > 65535 {
		return nil, errors.New("unable to connect. Invalid port")
	}

	cfg = append(cfg, fmt.Sprintf("remote %s %d", c.hostIP, c.hostPort))

	cfg = append(cfg, "resolv-retry infinite")
	if localPort > 0 {
		// NOTE:
		// Specifying the local port can lead to losing connectivity after OpenVPN RECONNECTING (observed on macOS)
		cfg = append(cfg, fmt.Sprintf("lport %d", localPort))
	} else {
		cfg = append(cfg, "nobind")
	}
	cfg = append(cfg, "persist-key")

	if _, err := os.Stat(platform.OpenvpnCaKeyFile()); os.IsNotExist(err) {
		return nil, errors.New("CA certificate not found")
	}
	cfg = append(cfg, fmt.Sprintf("ca \"%s\"", platform.OpenvpnCaKeyFile()))

	if _, err := os.Stat(platform.OpenvpnTaKeyFile()); os.IsNotExist(err) {
		return nil, errors.New("TLS auth key not found")
	}
	cfg = append(cfg, fmt.Sprintf("tls-auth \"%s\" 1", platform.OpenvpnTaKeyFile()))

	cfg = append(cfg, "cipher AES-256-CBC")
	cfg = append(cfg, "remote-cert-tls server")
	cfg = append(cfg, "verb 4")

	if upCmd := platform.OpenvpnUpScript(); upCmd != "" {
		cfg = append(cfg, "up \""+upCmd+"\"")
	}
	if downCmd := platform.OpenvpnDownScript(); downCmd != "" {
		cfg = append(cfg, "down \""+downCmd+"\"")
	}

	cfg = append(cfg, "script-security 2")

	if c.proxyAddress != nil && (c.proxyType == "http" || c.proxyType == "socks") {

		localGatewayAddress, err := netinfo.DefaultGatewayIP()
		if err != nil {
			return nil, fmt.Errorf("failed to get local gateway: %w", err)
		}

		if localGatewayAddress == nil {
			return nil, errors.New("internal error: LocalGatewayAddress not defined. Unable to generate OpenVPN configuration")
		}

		if c.proxyAddress.Equal(net.IPv4(127, 0, 0, 1)) {
			cfg = append(cfg, fmt.Sprintf("route %s 255.255.255.255 %s", c.hostIP.String(), localGatewayAddress.String()))
		} else {
			cfg = append(cfg, fmt.Sprintf("route %s 255.255.255.255 %s", c.proxyAddress, localGatewayAddress.String()))
		}
	}

	cfg, err = addUserDefinedParameters(cfg, openVpnCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to add user-defined parameters: %w", err)
	}

	return cfg, nil
}

// merge current parameters with user-defined parameters
func addUserDefinedParameters(currParams []string, userParams string) ([]string, error) {
	if len(userParams) <= 0 {
		return currParams, nil
	}

	// loop trough all openVpnCertificate defined by user
	// (looking if user-defined parameters overlap an existing parameters)
	tmpCfg := make([]string, 1)
	userLines := strings.Split(userParams, "\n")

	for _, cfgLine := range currParams {
		cfgParam := getParamFromConfigLine(cfgLine)
		cfgLineToSave := cfgLine

		for i, userLine := range userLines {
			userParam := getParamFromConfigLine(userLine)

			if len(userParam) > 0 && cfgParam == userParam {
				cfgLineToSave = userLine
				userLines[i] = ""
				break
			}
		}

		tmpCfg = append(tmpCfg, cfgLineToSave)
	}

	for _, userLine := range userLines {
		if len(userLine) > 0 {
			tmpCfg = append(tmpCfg, userLine)
		}
	}

	return tmpCfg, nil
}

func getParamFromConfigLine(line string) string {
	line = strings.TrimLeft(line, " \t")
	words := strings.Fields(line)

	if len(words) <= 0 || len(words[0]) <= 0 {
		return ""
	}
	// ignore comments
	if words[0][0] == '#' || words[0][0] == ';' {
		return ""
	}

	return strings.ToLower(words[0])
}

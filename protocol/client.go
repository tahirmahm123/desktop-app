package protocol

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	apitypes "github.com/ivpn/desktop-app-daemon/api/types"
	"github.com/ivpn/desktop-app-daemon/logger"
	"github.com/ivpn/desktop-app-daemon/protocol/types"
	"github.com/ivpn/desktop-app-daemon/vpn"
)

// Client for IVPN daemon
type Client struct {
	_port   int
	_secret uint64
	_conn   net.Conn

	_requestIdx int

	_defaultTimeout  time.Duration
	_receivers       map[*receiverChannel]struct{}
	_receiversLocker sync.Mutex

	_helloResponse types.HelloResp
}

// ResponseTimeout error
type ResponseTimeout struct {
}

func (e ResponseTimeout) Error() string {
	return "response timeout"
}

// CreateClient initialising new client for IVPN daemon
func CreateClient(port int, secret uint64) *Client {
	return &Client{
		_port:           port,
		_secret:         secret,
		_defaultTimeout: time.Second * 30,
		_receivers:      make(map[*receiverChannel]struct{})}
}

// Connect is connecting to daemon
func (c *Client) Connect() (err error) {
	if c._conn != nil {
		return fmt.Errorf("already connected")
	}

	logger.Info("Connecting...")

	c._conn, err = net.Dial("tcp", fmt.Sprintf(":%d", c._port))
	if err != nil {
		return fmt.Errorf("failed to connect to IVPN daemon (does IVPN daemon/service running?): %w", err)
	}

	logger.Info("Connected")

	// start receiver
	go c.receiverRoutine()

	if _, err := c.SendHello(); err != nil {
		return err
	}

	return nil
}

// SendHello - send initial message and get current status
func (c *Client) SendHello() (helloResponse types.HelloResp, err error) {
	if err := c.ensureConnected(); err != nil {
		return helloResponse, err
	}

	helloReq := types.Hello{Secret: c._secret, KeepDaemonAlone: true, GetStatus: true, Version: "1.0"}

	if err := c.sendRecvTimeOut(&helloReq, &c._helloResponse, time.Second*5); err != nil {
		if _, ok := errors.Unwrap(err).(ResponseTimeout); ok {
			return helloResponse, fmt.Errorf("Failed to send 'Hello' request (does another instance of IVPN Client running?): %w", err)
		}
		return helloResponse, fmt.Errorf("Failed to send 'Hello' request: %w", err)
	}
	return c._helloResponse, nil
}

// GetHelloResponse returns initialisation response from daemon
func (c *Client) GetHelloResponse() types.HelloResp {
	return c._helloResponse
}

// SessionNew creates new session
func (c *Client) SessionNew(accountID string, forceLogin bool) (apiStatus int, err error) {
	if err := c.ensureConnected(); err != nil {
		return 0, err
	}

	req := types.SessionNew{AccountID: accountID, ForceLogin: forceLogin}
	var resp types.SessionNewResp

	if err := c.sendRecv(&req, &resp); err != nil {
		return 0, err
	}

	if len(resp.Session.Session) <= 0 {
		return resp.APIStatus, fmt.Errorf("[%d] %s", resp.APIStatus, resp.APIErrorMessage)
	}

	return resp.APIStatus, nil
}

// SessionDelete remove session
func (c *Client) SessionDelete() error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	req := types.SessionDelete{}
	var resp types.EmptyResp

	if err := c.sendRecv(&req, &resp); err != nil {
		return err
	}

	return nil
}

// SessionStatus get session status
func (c *Client) SessionStatus() (ret types.SessionStatusResp, err error) {
	if err := c.ensureConnected(); err != nil {
		return ret, err
	}

	req := types.SessionStatus{}
	var resp types.SessionStatusResp

	if err := c.sendRecv(&req, &resp); err != nil {
		return ret, err
	}

	return resp, nil
}

// SetPreferences sends config parameter to daemon
// TODO: avoid using keys as a strings
func (c *Client) SetPreferences(key, value string) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	req := types.SetPreference{Key: key, Value: value}

	// TODO: daemon have to return confirmation
	if err := c.send(&req, 0); err != nil {
		return err
	}

	return nil
}

// FirewallSet change firewall state
func (c *Client) FirewallSet(isOn bool) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	// changing killswitch state
	req := types.KillSwitchSetEnabled{IsEnabled: isOn}
	var resp types.EmptyResp
	if err := c.sendRecv(&req, &resp); err != nil {
		return err
	}

	// requesting status
	state, err := c.FirewallStatus()
	if err != nil {
		return err
	}

	if state.IsEnabled != isOn {
		return fmt.Errorf("firewall state did not change [isEnabled=%v]", state.IsEnabled)
	}

	return nil
}

// FirewallAllowLan set configuration 'allow LAN'
func (c *Client) FirewallAllowLan(allow bool) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	// changing killswitch configuration
	req := types.KillSwitchSetAllowLAN{AllowLAN: allow, Synchronously: true}
	var resp types.EmptyResp
	if err := c.sendRecv(&req, &resp); err != nil {
		return err
	}

	return nil
}

// FirewallAllowLanMulticast set configuration 'allow LAN multicast'
func (c *Client) FirewallAllowLanMulticast(allow bool) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	// changing killswitch configuration
	req := types.KillSwitchSetAllowLANMulticast{AllowLANMulticast: allow, Synchronously: true}
	var resp types.EmptyResp
	if err := c.sendRecv(&req, &resp); err != nil {
		return err
	}

	return nil
}

// FirewallStatus get firewall state
func (c *Client) FirewallStatus() (state types.KillSwitchStatusResp, err error) {
	if err := c.ensureConnected(); err != nil {
		return state, err
	}

	// requesting status
	statReq := types.KillSwitchGetStatus{}
	if err := c.sendRecv(&statReq, &state); err != nil {
		return state, err
	}

	return state, nil
}

// GetServers gets servers list
func (c *Client) GetServers() (apitypes.ServersInfoResponse, error) {
	if err := c.ensureConnected(); err != nil {
		return apitypes.ServersInfoResponse{}, err
	}

	req := types.GetServers{}
	var resp types.ServerListResp

	if err := c.sendRecv(&req, &resp); err != nil {
		return resp.VpnServers, err
	}

	return resp.VpnServers, nil
}

// GetVPNState returns current VPN connection state
func (c *Client) GetVPNState() (vpn.State, types.ConnectedResp, error) {
	respConnected := types.ConnectedResp{}
	respDisconnected := types.DisconnectedResp{}
	respState := types.VpnStateResp{}

	if err := c.ensureConnected(); err != nil {
		return vpn.DISCONNECTED, respConnected, err
	}

	req := types.GetVPNState{}

	_, cmdBase, err := c.sendRecvAny(&req, &respConnected, &respDisconnected, &respState)
	if err != nil {
		return vpn.DISCONNECTED, respConnected, err
	}

	switch cmdBase.Command {
	case types.GetTypeName(respConnected):
		return vpn.CONNECTED, respConnected, nil

	case types.GetTypeName(respDisconnected):
		return vpn.DISCONNECTED, respConnected, nil

	case types.GetTypeName(respState):
		return respState.StateVal, respConnected, nil
	}

	return vpn.DISCONNECTED, respConnected, fmt.Errorf("failed to receive VPN state (not expected return type)")
}

// DisconnectVPN disconnect active VPN connection
func (c *Client) DisconnectVPN() error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	req := types.Disconnect{}
	respEmpty := types.EmptyResp{}
	respDisconnected := types.DisconnectedResp{}

	_, _, err := c.sendRecvAny(&req, &respEmpty, &respDisconnected)
	if err != nil {
		return err
	}

	return nil
}

// ConnectVPN - establish new VPN connection
func (c *Client) ConnectVPN(req types.Connect) (types.ConnectedResp, error) {
	respConnected := types.ConnectedResp{}
	respDisconnected := types.DisconnectedResp{}

	if err := c.ensureConnected(); err != nil {
		return respConnected, err
	}

	_, cmdBase, err := c.sendRecvAny(&req, &respConnected, &respDisconnected)
	if err != nil {
		return respConnected, err
	}

	switch cmdBase.Command {
	case types.GetTypeName(respConnected):
		return respConnected, nil

	case types.GetTypeName(respDisconnected):
		return respConnected, fmt.Errorf("%s", respDisconnected.ReasonDescription)
	}

	return respConnected, fmt.Errorf("connect request failed (not expected return type)")
}

// WGKeysGenerate regenerate WG keys
func (c *Client) WGKeysGenerate() error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	req := types.WireGuardGenerateNewKeys{}
	var resp types.EmptyResp
	if err := c.sendRecv(&req, &resp); err != nil {
		return err
	}

	return nil
}

// WGKeysRotationInterval changes WG keys rotation interval
func (c *Client) WGKeysRotationInterval(uinxTimeInterval int64) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	req := types.WireGuardSetKeysRotationInterval{Interval: uinxTimeInterval}
	var resp types.EmptyResp
	if err := c.sendRecv(&req, &resp); err != nil {
		return err
	}

	return nil
}

// PingServers changes WG keys rotation interval
func (c *Client) PingServers() (pingResults []types.PingResultType, err error) {
	if err := c.ensureConnected(); err != nil {
		return pingResults, err
	}

	req := types.PingServers{RetryCount: 4, TimeOutMs: 5000}
	var resp types.PingServersResp
	if err := c.sendRecv(&req, &resp); err != nil {
		return pingResults, err
	}

	return resp.PingResults, nil
}

// SetManualDNS - sets manual DNS for current VPN connection
func (c *Client) SetManualDNS(dns string) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	req := types.SetAlternateDns{DNS: dns}
	var resp types.SetAlternateDNSResp
	if err := c.sendRecv(&req, &resp); err != nil {
		return err
	}

	if resp.IsSuccess == false {
		return fmt.Errorf("DNS not changed")
	}

	return nil
}

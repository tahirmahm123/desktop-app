//
//  UI for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the UI for IVPN Client Desktop.
//
//  The UI for IVPN Client Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The UI for IVPN Client Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the UI for IVPN Client Desktop. If not, see <https://www.gnu.org/licenses/>.
//

import {
  enumValueName,
  isStrNullOrEmpty,
  getDistanceFromLatLonInKm,
} from "../helpers/helpers";
import {
  VpnTypeEnum,
  VpnStateEnum,
  PingQuality,
  PauseStateEnum,
  DnsEncryption,
} from "./types";

export default {
  namespaced: true,

  state: {
    connectionState: VpnStateEnum.DISCONNECTED,

    connectionInfo: null /*{
      VpnType: VpnTypeEnum.OpenVPN,
      ConnectedSince: new Date(),
      ClientIP: "",
      ServerIP: "",
      ExitServerID: "",
      ManualDNS: {
        DnsHost: "",      // string // DNS host IP address
	      Encryption: 0,    // DnsEncryption [	EncryptionNone = 0,	EncryptionDnsOverTls = 1,	EncryptionDnsOverHttps = 2]
	      DohTemplate: "",  // string // DoH/DoT template URI (for Encryption = DnsOverHttps or Encryption = DnsOverTls)
      },
      IsCanPause: null //(true/false)
    }*/,

    disconnectedInfo: {
      ReasonDescription: "",
    },

    pauseState: PauseStateEnum.Resumed,

    firewallState: {
      IsEnabled: null,
      IsPersistent: null,
      IsAllowLAN: null,
      IsAllowMulticast: null,
      IsAllowApiServers: null,
    },

    // The split-tunnelling configuration
    splitTunnelling: {
      IsEnabled: false,
      IsCanGetAppIconForBinary: false,
      // Split-Tunnelling (SplitTunnelStatus)
      // IsEnabled bool                     - Is ST enabled
      // IsCanGetAppIconForBinary bool      - This parameter informs availability of the functionality to get icon for particular binary
      //                                      (true - if commands GetAppIcon/AppIconResp  applicable for this platform)
      // SplitTunnelApps []string           - Information about applications added to ST configuration
      //                                      (applicable for Windows)
      // RunningApps []splittun.RunningApp  - Information about active applications running in Split-Tunnel environment
      //                                      (applicable for Linux)
      //                                      type RunningApp struct:
      //                                        Pid     int
      //                                        Ppid    int // The PID of the parent of this process.
      //                                        Cmdline string
      //                                        Exe         string  // The actual pathname of the executed command
      //                                        ExtIvpnRootPid int  // PID of the known parent process registered by AddPid() function
      //                                        ExtModifiedCmdLine string
    },

    dns: {
      DnsHost: "",
      Encryption: DnsEncryption.None,
      DohTemplate: "",
    },

    currentWiFiInfo: null, //{ SSID: "", IsInsecureNetwork: false },
    availableWiFiNetworks: null, // []{SSID: ""}

    // Servers hash object: serversHashed[gateway] = server
    serversHashed: {},
    servers: {
      data: [],
      certificate: "",
    },

    // true when servers pinging in progress
    isPingingServers: false,

    /*
    // SERVERS
    servers: {
      wireguard: [
        {
          gateway: "",
          country_code: "",
          country: "",
          city: "",
          latitude: 0,
	        longitude: 0,

          ping: ??? // property added after receiving ping info from daemon
          pingQuality: ??? // PingQuality (Good, Moderate, Bad) - property calculated after receiving ping info from daemon
          isIPv6: ??? // property calculated automatically after receiving servers data from the daemon (null - IPv6 not supported by this type of VPN)

          hosts: [
            {
              hostname: "",
              host: "",
              public_key: "",
              local_ip: "",
              ipv6:
              {
                local_ip: "",
                host: "",
                multihop_port: 0
              }
            }
          ]
        }
      ],
      openvpn: [
        {
          gateway: "",
          country_code: "",
          country: "",
          city: "",
          latitude: 0,
	        longitude: 0,
          ping: ??? // property added after receiving ping info from daemon
          pingQuality: ??? // PingQuality (Good, Moderate, Bad) - property calculated after receiving ping info from daemon

          hosts: [
            {
              hostname: "",
              host: "",
              multihop_port: 0
            }
          ]
        }
      ],
      config: {
        antitracker: {
          default: { ip: "", "multihop-ip": "" },
          hardcore: { ip: "", "multihop-ip": "" }
        },
        api: { ips: [""], ipv6s:[""] }
      }
    }*/
  },

  mutations: {
    connectionState(state, cs) {
      state.connectionState = cs;
      if (cs == VpnStateEnum.DISCONNECTED)
        state.pauseState = PauseStateEnum.Resumed;
    },
    connectionInfo(state, ci) {
      state.connectionInfo = ci;
      if (ci != null) {
        state.connectionState = VpnStateEnum.CONNECTED;
        state.disconnectedInfo = null;
      }
    },
    disconnected(state, disconnectionReason) {
      state.disconnectedInfo = { ReasonDescription: disconnectionReason };
      state.connectionState = VpnStateEnum.DISCONNECTED;
      state.pauseState = PauseStateEnum.Resumed;
      state.connectionInfo = null;
    },
    pauseState(state, val) {
      state.pauseState = val;
    },
    servers(state, serversObj) {
      updateServers(state, serversObj);
    },
    isPingingServers(state, val) {
      state.isPingingServers = val;
    },
    serversPingStatus(state, pingResultArray) {
      updateServersPings(state, pingResultArray);
    },
    firewallState(state, obj) {
      state.firewallState = obj;
    },
    // Split-Tunnelling
    splitTunnelling(state, val) {
      state.splitTunnelling = val;
    },
    dns(state, dns) {
      state.dns = dns;
    },

    currentWiFiInfo(state, currentWiFiInfo) {
      if (currentWiFiInfo != null && currentWiFiInfo.SSID == "")
        state.currentWiFiInfo = null;
      else state.currentWiFiInfo = currentWiFiInfo;
    },
    availableWiFiNetworks(state, availableWiFiNetworks) {
      state.availableWiFiNetworks = availableWiFiNetworks;
    },
  },

  getters: {
    isDisconnecting: (state) => {
      return state.connectionState === VpnStateEnum.DISCONNECTING;
    },
    isDisconnected: (state) => {
      return state.connectionState === VpnStateEnum.DISCONNECTED;
    },
    isConnecting: (state) => {
      switch (state.connectionState) {
        case VpnStateEnum.CONNECTING:
        case VpnStateEnum.WAIT:
        case VpnStateEnum.AUTH:
        case VpnStateEnum.GETCONFIG:
        case VpnStateEnum.ASSIGNIP:
        case VpnStateEnum.ADDROUTES:
        case VpnStateEnum.RECONNECTING:
        case VpnStateEnum.TCP_CONNECT:
          return true;
        default:
          return false;
      }
    },
    isConnected: (state) => {
      return state.connectionState === VpnStateEnum.CONNECTED;
    },
    vpnStateText: (state) => {
      return enumValueName(VpnStateEnum, state.connectionState);
    },
    activeServers(state, getters, rootState) {
      return getActiveServers(state, rootState);
    },
    antitrackerIp(state, getters, rootState) {
      let atConfig = state.servers.config.antitracker;
      if (atConfig == null) return null;
      let atIPs = rootState.settings.isAntitrackerHardcore
        ? atConfig.hardcore
        : atConfig.default;
      if (atIPs == null) return null;

      if (rootState.settings.vpnType === VpnTypeEnum.WireGuard) {
        // WireGuard (port-based MultiHop using the same DNS IP as SingleHop)
        return atIPs.ip;
      }
      // OpenVPN
      return rootState.settings.isMultiHop ? atIPs["multihop-ip"] : atIPs.ip;
    },
    isAntitrackerEnabled: (state) => {
      return isAntitrackerActive(state);
    },
    isAntitrackerHardcoreEnabled: (state) => {
      return isAntitrackerHardcoreActive(state);
    },
    fastestServer(state, getters, rootState) {
      let servers = getActiveServers(state, rootState);
      if (servers == null || servers.length <= 0) return null;

      let skipSvrs = rootState.settings.serversFastestExcludeList;
      let retSvr = null;

      // If there will not be any server with ping-info -
      // save the info about the first applicable server (which is not in skipSvrs)
      let fallbackSvr = null;

      for (let i = 0; i < servers.length; i++) {
        let curSvr = servers[i];
        if (!curSvr) continue;
        if (skipSvrs != null && skipSvrs.includes(curSvr.gateway)) continue;
        if (!fallbackSvr) fallbackSvr = curSvr;
        if (
          curSvr != null &&
          curSvr.ping &&
          curSvr.ping > 0 &&
          (retSvr == null || retSvr.ping > curSvr.ping)
        )
          retSvr = curSvr;
      }

      if (!retSvr) {
        // No fastest server detected (due to no ping info available)
        // Get nearest or first applicable server

        // get last known location
        const l = rootState.lastRealLocation;
        if (l) {
          try {
            // distance compare
            let compare = function (a, b) {
              var distA = getDistanceFromLatLonInKm(
                l.latitude,
                l.longitude,
                a.latitude,
                a.longitude
              );
              var distB = getDistanceFromLatLonInKm(
                l.latitude,
                l.longitude,
                b.latitude,
                b.longitude
              );
              if (distA === distB) return 0;
              if (distA < distB) return -1;
              return 1;
            };

            // sort servers by distance from last known real location
            let sortedSvrs = servers.slice().sort(compare);
            // get nearest server
            for (let i = 0; i < sortedSvrs.length; i++) {
              let curSvr = servers[i];
              if (skipSvrs != null && skipSvrs.includes(curSvr.gateway))
                continue;
              retSvr = curSvr;
              break;
            }
          } catch (e) {
            console.log(e);
          }
        }

        // If still not found: choose the first applicable server
        if (!retSvr) retSvr = fallbackSvr;
      }

      return retSvr;
    },
  },

  // can be called from renderer
  actions: {
    connectionInfo(context, ci) {
      // save current connection info
      context.commit("connectionInfo", ci);

      // Received 'connected' state
      // Connection can be triggered outside (not by current application instance)
      // So, we should just update received data in settings (vpnType, multihop, entry\exit servers)
      // (no consistency checks should be performed)
      const isMultiHop = isStrNullOrEmpty(ci.ExitServerID) ? false : true;
      context.commit("settings/vpnType", ci.VpnType, { root: true });
      context.dispatch("settings/isMultiHop", isMultiHop, { root: true });
      // it is important to read 'activeServers' only after vpnType was updated!
      const servers = context.getters.activeServers;
      const entrySvr = findServerByIp(servers, ci.ServerIP);
      context.commit("settings/serverEntry", entrySvr, { root: true });
      if (!isStrNullOrEmpty(ci.ExitServerID)) {
        const exitSvr = findServerByExitId(servers, ci.ExitServerID);
        context.commit("settings/serverExit", exitSvr, { root: true });
      }

      // save last DNS state
      context.commit("dns", ci.ManualDNS);
      updateDnsSettings(context);
    },
    pauseState(context, val) {
      context.commit("pauseState", val);

      if (val === PauseStateEnum.Resumed || val === PauseStateEnum.Resuming)
        context.dispatch("uiState/pauseConnectionTill", null, { root: true });
    },
    servers(context, value) {
      context.commit("servers", value);
      // notify 'settings' module about updated servers list
      // (it is required to update selected servers, if necessary)
      context.dispatch("settings/updateSelectedServers", null, { root: true });
    },
    // Split-Tunnelling
    splitTunnelling(state, val) {
      state.splitTunnelling = val;
    },
    dns(context, dns) {
      context.commit("dns", dns);
      // save current state to settings
      updateDnsSettings(context);
    },
  },
};

function updateDnsSettings(context) {
  // save current state to settings
  const isAntitracker = isAntitrackerActive(context.state);
  context.dispatch("settings/isAntitracker", isAntitracker, { root: true });

  if (isAntitracker === true) {
    const isAntitrackerHardcore = isAntitrackerHardcoreActive(context.state);
    context.dispatch("settings/isAntitrackerHardcore", isAntitrackerHardcore, {
      root: true,
    });
  }

  if (isAntitracker === false) {
    let currDnsState = context.state.dns;

    let isCustomDns = true;
    if (currDnsState == null || !currDnsState.DnsHost) isCustomDns = false;
    else
      context.dispatch("settings/dnsCustomCfg", currDnsState, { root: true });

    context.dispatch("settings/dnsIsCustom", isCustomDns, { root: true });
  }
}

function getActiveServers(state, rootState) {
  const vpnType = rootState.settings.vpnType;
  // const enableIPv6InTunnel = rootState.settings.enableIPv6InTunnel;
  // const showGatewaysWithoutIPv6 = rootState.settings.showGatewaysWithoutIPv6;
  console.log(vpnType);
  // if (vpnType === VpnTypeEnum.OpenVPN) {
  // IPv6 in not implemented for OpenVPN
  return state.servers.data;
  // }

  /*let wgServers = state.servers.wireguard;
  if (enableIPv6InTunnel == true && showGatewaysWithoutIPv6 != true) {
    // show only servers which support IPv6
    return wgServers.filter((s) => {
      return s.isIPv6;
    });
  }

  return wgServers;*/
}

function findServerByIp(servers, ip) {
  for (let i = 0; i < servers.length; i++) {
    const srv = servers[i];

    if (srv.hosts != null) {
      // wireguard/openvpn server
      for (let j = 0; j < srv.hosts.length; j++) {
        if (srv.hosts[j].host === ip) return srv;
      }
    }
  }
  return null;
}

function findServerByExitId(servers, id) {
  for (let i = 0; i < servers.length; i++) {
    const srv = servers[i];
    if (srv.gateway == null) continue;
    if (id === srv.gateway.split(".")[0]) return srv;
  }
}

function updateServersPings(state, pings) {
  // hash new ping result by host
  let hashedPings = {};
  for (let i = 0; i < pings.length; i++) {
    hashedPings[pings[i].Host] = pings[i].Ping;
  }

  function getPingQuality(pingMs) {
    if (pingMs < 100) return PingQuality.Good;
    if (pingMs < 300) return PingQuality.Moderate;
    return PingQuality.Bad;
  }

  let funcGetPing = function (s) {
    for (let i = 0; i < s.hosts.length; i++) {
      let pingValFoHost = hashedPings[s.hosts[i].ip];
      if (pingValFoHost != null) {
        s.ping = pingValFoHost;
        s.pingQuality = getPingQuality(s.ping);
        break;
      }
    }
  };

  // state.servers.wireguard.forEach((s) => {
  //   funcGetPing(s);
  // });

  state.servers.data.forEach((s) => {
    funcGetPing(s);
  });
}

function isServerSupportIPv6(server) {
  if (!server) return null;
  if (!server.hosts) return null;

  for (let h of server.hosts) {
    if (h && h.ipv6 && h.ipv6.local_ip) return true;
  }
  return false;
}

function updateServers(state, newServers) {
  if (newServers == null) return;
  console.log("new servers", newServers);
  // ensure all required properties are defined (even with empty values)
  let serversEmpty = {
    data: [],
    certificate: "",
  };
  newServers = Object.assign(serversEmpty, newServers);

  // prepare hash for new servers (hash by gateway id)
  function initNewServersAndCreateHash(hashObj, servers) {
    let retObj = hashObj;
    if (retObj == null) retObj = {};
    for (let i = 0; i < servers.length; i++) {
      let svr = servers[i];
      svr.ping = null; // initialize 'ping' field to support VUE reactivity for it
      svr.pingQuality = null;
      svr.isIPv6 = isServerSupportIPv6(svr);
      retObj[svr.flag] = svr; // hash
    }
    return retObj;
  }

  // let hash = initNewServersAndCreateHash(null, newServers.wireguard);
  state.serversHashed = initNewServersAndCreateHash(null, newServers.data);
  console.log(state.serversHashed);

  // copy ping value from old objects
  function copySvrsDataFromOld(oldServers, newServersHashed) {
    for (let i = 0; i < oldServers.length; i++) {
      let oldSrv = oldServers[i];
      let newSrv = newServersHashed[oldSrv.gateway];
      if (newSrv == null) {
        continue;
      }
      newSrv.ping = oldSrv.ping;
      newSrv.pingQuality = oldSrv.pingQuality;
    }
  }

  copySvrsDataFromOld(state.servers.data, state.serversHashed);

  // sort new servers (by country/city)
  // function compare(a, b) {
  //   let ret = a.country_code.localeCompare(b.country_code);
  //   if (ret != 0) return ret;
  //   return a.city.localeCompare(b.city);
  // }

  // newServers.data.sort(compare);

  // save servers
  state.servers = newServers;
}

function isAntitrackerActive(state) {
  let dnsIP = state.dns.DnsHost;
  if (isStrNullOrEmpty(dnsIP) || state.dns.Encryption != DnsEncryption.None)
    return false;

  let atConfig = state.servers.config.antitracker;
  switch (dnsIP) {
    case atConfig.default.ip:
    case atConfig.hardcore.ip:
    case atConfig.default["multihop-ip"]:
    case atConfig.hardcore["multihop-ip"]:
      return true;
    default:
  }
  return false;
}

function isAntitrackerHardcoreActive(state) {
  let dnsIP = state.dns.DnsHost;
  if (isStrNullOrEmpty(dnsIP) || state.dns.Encryption != DnsEncryption.None)
    return false;

  let atConfig = state.servers.config.antitracker;
  switch (dnsIP) {
    case atConfig.hardcore.ip:
    case atConfig.hardcore["multihop-ip"]:
      return true;
    default:
  }
  return false;
}

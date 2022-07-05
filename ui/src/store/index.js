//
//  UI for VPN Client Desktop
//  https://github.com/tahirmahm123/vpn-desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the UI for VPN Client Desktop.
//
//  The UI for VPN Client Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The UI for VPN Client Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the UI for VPN Client Desktop. If not, see <https://www.gnu.org/licenses/>.
//

import Vue from "vue";
import Vuex from "vuex";
import createSharedMutations from "./vuex-modules/shared-mutations.js";

import { isStrNullOrEmpty, IsRenderer } from "../helpers/helpers";

import account from "./module-account";
import vpnState from "./module-vpn-state";
import uiState from "./module-ui-state";
import settings from "./module-settings";
import { VpnTypeEnum } from "@/store/types";

Vue.use(Vuex);

let sharedMutationsOptions = {};
if (IsRenderer()) {
  // renderer
  sharedMutationsOptions.type = "renderer";
  sharedMutationsOptions.ipcRenderer = window.ipcSender.GetSafeIpcRenderer();
} else {
  // main
  sharedMutationsOptions.type = "main";
}

export default new Vuex.Store({
  plugins: [createSharedMutations(sharedMutationsOptions)],

  modules: { account, vpnState, uiState, settings },
  strict: true,
  state: {
    daemonConnectionState: null, // DaemonConnectionType from "./types";
    daemonVersion: "",
    daemonIsOldVersionError: false,
    daemonIsInstalling: false,

    configParams: {
      UserDefinedOvpnFile: "",
    },

    disabledFunctions: {
      WireGuardError: "",
      OpenVPNError: "",
      ObfsproxyError: "",
      SplitTunnelError: "",
    },

    dnsAbilities: {
      CanUseDnsOverTls: false,
      CanUseDnsOverHttps: false,
    },

    dnsPredefinedConfigurations: null, //[]  array elements: { DnsHost: "", Encryption: DnsEncryption.None, DohTemplate: "", }

    // true when we are requesting geo-lookup info on current moment
    isRequestingLocation: false,
    // Current location (be careful, in 'connected' state this object will contain info about 'VPN location')
    location: null, // {"ip_address":"","isp":"","organization":"","country":"","country_code":"","city":"","latitude": 0.0,"longitude":0.0,"isVpnServer":false}

    // true when we are requesting geo-lookup info on current moment
    isRequestingLocationIPv6: false,
    // Current location (be careful, in 'connected' state this object will contain info about 'VPN location')
    locationIPv6: null, // {"ip_address":"","isp":"","organization":"","country":"","country_code":"","city":"","latitude": 0.0,"longitude":0.0,"isVpnServer":false}

    // Contains current user (real) location OR (if connected) the last real user location
    // This parameter is using, for example, for calculating distance to a nearest server
    lastRealLocation: null,

    // Updates info example:
    /*
    // NOTE: if section 'generic' defined - the 'daemon' and 'uiClient' must be ignored
    {
      "generic": {
        "version": "3.2.45",
        "downloadLink": "https://repo.vpn.net/binaries/audit2021/VPN-3.2.40.dmg",
        "signature":    "https://repo.vpn.net/binaries/audit2021/VPN-3.2.40.dmg.sign.sha256.base64",
        "releaseNotes": [
          {
            "type": "new",
            "description": "New feature description"
          },
          {
            "type": "improved",
            "description": "Improvement description"
          },
          {
            "type": "fix",
            "description": "UI Bugfix description"
          }
        ]
      },
      "daemon": {
        "version": "2.12.7",
        "releaseNotes": [
          {
            "type": "new",
            "description": "New feature description"
          },
          {
            "type": "improved",
            "description": "Improvement description"
          }
        ]
      },
      "uiClient": {
        "version": "3.0.8",
        "releaseNotes": [
          {
            "type": "fix",
            "description": "UI Bugfix description"
          }
        ]
      }
    }*/
    latestVersionInfo: null,
  },

  getters: {
    getLastRealLocation: (state) => state.lastRealLocation,
    isCanUseIPv6InTunnel: (state) => {
      return (
        state.settings.vpnType !== VpnTypeEnum.OpenVPN // IPv6 is not implemented for OpenVPN yet
      );
    },
    isIPv4andIPv6LocationsEqual: (state) => {
      let l4 = state.location;
      let l6 = state.locationIPv6;

      if (!l4 || !l6) return true;

      if (
        l4.country != l6.country ||
        l4.city != l6.city ||
        (l4.isp != l6.isp && l4.isVpnServer != l6.isVpnServer)
      )
        return false;

      return true;
    },
    getIsInfoAvailableIPv4: (state) => {
      let l = state.location;
      if (l && l.ip_address) return true;
      return false;
    },
    getIsInfoAvailableIPv6: (state) => {
      let l = state.locationIPv6;
      if (l && l.ip_address) return true;
      return false;
    },
    getIsIPv6View: (state, getters) => {
      if (
        !getters.getIsInfoAvailableIPv4 &&
        !state.isRequestingLocation &&
        getters.getIsInfoAvailableIPv6
      )
        return true;

      return (
        state.uiState.isIPv6View &&
        (getters.getIsInfoAvailableIPv6 || state.isRequestingLocationIPv6)
      );
    },

    isWireGuardEnabled: (state) =>
      isStrNullOrEmpty(state.disabledFunctions.WireGuardError),
    isOpenVPNEnabled: (state) =>
      isStrNullOrEmpty(state.disabledFunctions.OpenVPNError),
    isObfsproxyEnabled: (state) =>
      isStrNullOrEmpty(state.disabledFunctions.ObfsproxyError),
    isSplitTunnelEnabled: (state) =>
      isStrNullOrEmpty(state.disabledFunctions.SplitTunnelError),

    getDnsAbilities: (state) => {
      return state.dnsAbilities;
    },
  },

  // can be called from main process
  mutations: {
    replaceState(state, val) {
      Object.assign(state, val);
    },
    daemonConnectionState(state, value) {
      state.daemonConnectionState = value;
    },
    daemonIsOldVersionError(state, value) {
      state.daemonIsOldVersionError = value;
    },
    daemonVersion(state, value) {
      state.daemonVersion = value;
    },
    daemonIsInstalling(state, value) {
      state.daemonIsInstalling = value;
    },
    latestVersionInfo(state, value) {
      state.latestVersionInfo = value;
    },
    disabledFunctions(state, disabledFuncs) {
      state.disabledFunctions = disabledFuncs;
    },
    dnsAbilities(state, dnsAbilities) {
      if (!dnsAbilities)
        dnsAbilities = { CanUseDnsOverTls: false, CanUseDnsOverHttps: false };

      state.dnsAbilities = dnsAbilities;
    },
    dnsPredefinedConfigurations(state, dnsPredefinedConfigurations) {
      state.dnsPredefinedConfigurations = dnsPredefinedConfigurations;
    },

    configParams(state, value) {
      state.configParams = value;
    },

    // LOCATION
    location(state, geoLocation) {
      state.location = geoLocation;

      if (!this.getters["vpnState/isConnected"]) {
        // save only real user location
        state.lastRealLocation = geoLocation;
      }
    },
    locationIPv6(state, geoLocation) {
      state.locationIPv6 = geoLocation;

      if (!this.getters["vpnState/isConnected"] && !state.location) {
        // save only real user location
        state.lastRealLocation = geoLocation;
      }
    },

    isRequestingLocation(state, value) {
      state.isRequestingLocation = value;
    },
    isRequestingLocationIPv6(state, value) {
      state.isRequestingLocationIPv6 = value;
    },
  },

  // can be called from renderer
  actions: {},
});

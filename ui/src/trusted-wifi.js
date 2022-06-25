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

import store from "@/store";
import daemonClient from "./daemon-client";

let lastProcessedRule = null; //{ SSID: null, isTrusted: null}

export function InitTrustedNetworks() {
  store.subscribe((mutation) => {
    try {
      if (
        mutation.type === "vpnState/currentWiFiInfo" ||
        mutation.type === "settings/wifi"
      ) {
        setTimeout(async () => processWifiChange(), 0);
      }
    } catch (e) {
      console.error("Error in InitTrustedNetworks handler", e);
    }
  });
}

// returns 'false' when not allowed to perform automatic VPN connection (e.g. on application startup):
//    - 'trustedNetworksControl' is enabled
//       AND the current WIFI is 'trusted'
//       AND 'disconnect VPN for trusted networks' action is enabled
export function IsCanAutoConnectForCurrentSSID() {
  const isLoggedIn = store.getters["account/isLoggedIn"];
  if (!isLoggedIn) {
    return false;
  }

  let wifi = store.state.settings.wifi;
  if (!wifi || !wifi.actions || wifi.trustedNetworksControl != true)
    return true;

  let currentWiFiInfo = store.state.vpnState.currentWiFiInfo;
  if (!currentWiFiInfo || !currentWiFiInfo.SSID) return true;

  const isTrusted = getTrustRuleForConfiguredNetwork(
    currentWiFiInfo.SSID,
    wifi.networks
  );

  if (isTrusted == true && wifi.actions.trustedDisconnectVpn == true)
    return false;
  return true;
}

async function processWifiChange() {
  // 1. trying to apply rules (if network is configured)
  // 2. if network not configured -> apply rules for default trust status
  // 3. if default trust status not defined -> for insecure network: connect VPN
  // Returns 'true' if a rule was applied for current WIFI network (including rule for 'connectVPNOnInsecureNetwork')

  const isLoggedIn = store.getters["account/isLoggedIn"];
  if (!isLoggedIn) {
    return false;
  }

  // trusted networks config
  let wifi = store.state.settings.wifi;
  if (!wifi) return false;

  // current network
  let currentWiFiInfo = store.state.vpnState.currentWiFiInfo;
  let isInsecureNetwork = null;
  let currSSID = null;
  if (currentWiFiInfo != null) {
    isInsecureNetwork = currentWiFiInfo.IsInsecureNetwork;
    currSSID = currentWiFiInfo.SSID;
  }

  if (!currSSID) {
    lastProcessedRule = null;
    return false;
  }

  // if trusted network control is enabled
  if (wifi.trustedNetworksControl == true && wifi.actions) {
    // get configuration for current network
    let trustRule = getTrustRuleForConfiguredNetwork(currSSID, wifi.networks);

    // if network not configured - get default trust operation for not configured networks
    if (trustRule == null) trustRule = wifi.defaultTrustStatusTrusted;

    if (trustRule != null) {
      // skip applying same rule if we did it already (for same network with the same actions)
      if (
        lastProcessedRule != null &&
        lastProcessedRule.SSID != null &&
        lastProcessedRule.SSID == currSSID &&
        lastProcessedRule.isTrusted == trustRule
      )
        return false;

      // apply rule
      await applyTrustRule(trustRule, wifi.actions);
      lastProcessedRule = {
        SSID: currSSID,
        isTrusted: trustRule,
      };
      return true;
    }
  }

  // check is it insecure network (if network still not processed by 'trusted networks' configuration)
  if (
    isInsecureNetwork == true &&
    wifi.connectVPNOnInsecureNetwork == true &&
    !(
      store.getters["vpnState/isConnected"] ||
      store.getters["vpnState/isConnecting"]
    )
  ) {
    // skip applying same rule if we did it already (for same network)
    if (
      lastProcessedRule != null &&
      lastProcessedRule.SSID != null &&
      lastProcessedRule.SSID == currSSID
    )
      return false;

    console.log(
      "Joined insecure network. Connecting (according to preferences) ..."
    );
    await daemonClient.Connect();
    lastProcessedRule = {
      SSID: currSSID,
      isTrusted: null,
    };
    return true;
  }

  lastProcessedRule = null;
  return false;
}

function getTrustRuleForConfiguredNetwork(currSSID, networks) {
  if (!currSSID || networks == null) return null;

  // check configuration for current network
  let networkConfigArr = networks.filter((wifi) => wifi.ssid == currSSID);
  if (networkConfigArr == null || networkConfigArr.length == 0) return null;
  let networkConfig = networkConfigArr[0];

  return networkConfig.isTrusted;
}

async function applyTrustRule(isTrusted, actions) {
  if (isTrusted == null || actions == null) return;
  if (isTrusted) {
    // trusted
    if (
      actions.trustedDisconnectVpn == true &&
      !store.getters["vpnState/isDisconnected"]
    ) {
      console.log(
        "Joined trusted network. Disconnecting (according to preferences) ..."
      );
      await daemonClient.Disconnect();
    }
    if (
      actions.trustedDisableFirewall == true &&
      store.state.vpnState.firewallState.IsEnabled != false
    ) {
      console.log(
        "Joined trusted network. Disabling firewall (according to preferences) ..."
      );
      await daemonClient.EnableFirewall(false);
    }
  } else {
    // untrusted
    if (
      actions.unTrustedEnableFirewall == true &&
      store.state.vpnState.firewallState.IsEnabled != true
    ) {
      console.log(
        "Joined untrusted network. Enabling firewall (according to preferences) ..."
      );
      await daemonClient.ResumeConnection();
      await daemonClient.EnableFirewall(true);
    }
    if (
      actions.unTrustedConnectVpn == true &&
      !(
        store.getters["vpnState/isConnected"] ||
        store.getters["vpnState/isConnecting"]
      )
    ) {
      console.log(
        "Joined untrusted network. Connecting (according to preferences) ..."
      );

      await daemonClient.Connect();
    }
  }
}

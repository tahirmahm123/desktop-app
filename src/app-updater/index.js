//
//  UI for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app-ui2
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

import store from "@/store";
import { Platform, PlatformEnum } from "@/platform/platform";
const NoUpdaterErrorMessage = "App updater not available for this platform";

function getUpdater() {
  try {
    // Can be loaded different updaters according to platform
    if (IsGenericUpdater()) return require("./updater_generic");
    return require("./updater_linux");
  } catch (e) {
    console.error("[ERROR] IsAbleToCheckUpdate :", e);
  }
  return null;
}

export function IsGenericUpdater() {
  if (Platform() === PlatformEnum.Linux) return false;
  return true;
}

export function IsAbleToCheckUpdate() {
  const updater = getUpdater();
  if (updater == null) return false;
  return true;
}

export function StartUpdateChecker(onHasUpdateCallback) {
  const updater = getUpdater();
  if (updater == null) {
    console.warn(NoUpdaterErrorMessage);
    return false;
  }

  try {
    const currDaemonVer = store.state.daemonVersion;
    const currUiVer = require("electron").app.getVersion();
    if (!currDaemonVer || !currUiVer) {
      console.warn(
        "Unable to start app update checker: current app versions undefined"
      );
      return false;
    }

    const doCheck = async function() {
      const updatesInfo = await CheckUpdates();

      try {
        if (
          updatesInfo &&
          onHasUpdateCallback &&
          updater.IsNewerVersion(updatesInfo, currDaemonVer, currUiVer)
        ) {
          onHasUpdateCallback(updatesInfo, currDaemonVer, currUiVer);
        }
      } catch (e) {
        console.error(e);
        return;
      }
    };
    // check for updates in 5 seconds after initialization
    setTimeout(doCheck, 1000 * 5);

    // start periodical update check
    setInterval(doCheck, 1000 * 60 * 60 * 12); // 12-hours interval
  } catch (e) {
    console.error(e);
    return false;
  }
  return true;
}

export async function CheckUpdates() {
  const updater = getUpdater();
  if (updater == null) {
    console.error("App updater not available for this platform");
    return null;
  }

  console.log("Checking for app updates...");
  try {
    let updatesInfo = await updater.CheckUpdates();

    if (!updatesInfo) return null;

    store.commit("latestVersionInfo", updatesInfo);
    return updatesInfo;
  } catch (e) {
    console.error(e);
  }
  return null;
}

export function Upgrade() {
  const updater = getUpdater();
  if (updater == null) {
    console.warn(NoUpdaterErrorMessage);
    return null;
  }
  if (updater.Upgrade) return updater.Upgrade(store.state.latestVersionInfo);
}

export function CancelDownload() {
  const updater = getUpdater();
  if (updater == null) {
    console.warn(NoUpdaterErrorMessage);
    return null;
  }
  if (updater.CancelDownload) return updater.CancelDownload();
}

export function Install() {
  const updater = getUpdater();
  if (updater == null) {
    console.warn(NoUpdaterErrorMessage);
    return null;
  }
  if (updater.Install) return updater.Install();
}

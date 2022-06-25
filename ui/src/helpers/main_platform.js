import { Platform, PlatformEnum } from "@/platform/platform";
const os = require("os");

async function winInstallFolder() {
  return await new Promise((resolve, reject) => {
    let Registry = require("winreg");
    let regKey = new Registry({
      hive: Registry.HKLM,
      key: "\\Software\\VPN Client",
    });

    regKey.get(Registry.DEFAULT_VALUE, function (err, item) {
      if (err) reject(`Error reading installation path (registry):${err}`);
      else resolve(item.value);
    });
  });
}

export async function GetPortInfoFilePath() {
  switch (Platform()) {
    case PlatformEnum.macOS:
      return "/Library/Application Support/VPN/port.txt";
    case PlatformEnum.Linux:
      return "/opt/vpn/mutable/port.txt";
    case PlatformEnum.Windows: {
      let dir = await winInstallFolder();
      return `${dir}\\etc\\port.txt`;
    }
    default:
      throw new Error(`Not supported platform: '${os.platform()}'`);
  }
}

export async function GetOpenSSLBinaryPath() {
  switch (Platform()) {
    case PlatformEnum.macOS:
      return "/usr/bin/openssl";
    case PlatformEnum.Linux:
      return "/usr/bin/openssl";
    case PlatformEnum.Windows: {
      if (os.arch() === "x64") {
        let dir = await winInstallFolder();
        return `${dir}\\OpenVPN\\x86_64\\openssl.exe`;
      } else throw new Error(`Not supported architecture: '${os.arch()}'`);
    }
    default:
      throw new Error(`Not supported platform: '${os.platform()}'`);
  }
}

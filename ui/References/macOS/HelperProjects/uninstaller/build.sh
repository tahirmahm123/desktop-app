#!/bin/bash

#save current dir
_BASE_DIR="$( pwd )"
_SCRIPT=`basename "$0"`
#enter the script folder
cd "$(dirname "$0")"
_SCRIPT_DIR="$( pwd )"

# check result of last executed command
function CheckLastResult
{
  if ! [ $? -eq 0 ]; then #check result of last command
    if [ -n "$1" ]; then
      echo $1
    else
      echo "FAILED"
    fi
    exit 1
  fi
}

# The Apple DevID certificate which will be used to sign VPN Agent (Daemon) binary
# The helper will check VPN Agent signature with this value
_SIGN_CERT="" # E.g. "WXXXXXXXXN". Specific value can be passed by command-line argument: -c <APPLE_DEVID_SERT>
while getopts ":c:" opt; do
  case $opt in
    c) _SIGN_CERT="$OPTARG"
    ;;
  esac
done

if [ -z "${_SIGN_CERT}" ]; then
  echo "Usage:"
  echo "    $0 -c <APPLE_DEVID_CERTIFICATE>"
  echo "    Example: $0 -c WXXXXXXXXN"
  exit 1
fi

if [ ! -f "../helper/net.vpn.client.Helper" ]; then
  echo " File not exists '../helper/net.vpn.client.Helper'. Please, compile helper project first."
  exit 1
fi

rm -fr bin
CheckLastResult

echo "[ ] *** Compiling VPN Installer / Uninstaller ***"

echo "[+] VPN Installer: updating certificate info in .plist ..."
echo "    Apple DevID certificate: '${_SIGN_CERT}'"
plutil -replace SMPrivilegedExecutables -xml \
        "<dict> \
      		<key>net.vpn.client.Helper</key> \
      		<string>identifier net.vpn.client.Helper and certificate leaf[subject.OU] = ${_SIGN_CERT}</string> \
      	</dict>" "VPN Installer-Info.plist" || CheckLastResult
plutil -replace SMPrivilegedExecutables -xml \
        "<dict> \
          <key>net.vpn.client.Helper</key> \
          <string>identifier net.vpn.client.Helper and certificate leaf[subject.OU] = ${_SIGN_CERT}</string> \
        </dict>" "VPN Uninstaller-Info.plist" || CheckLastResult

echo "[+] VPN Installer: make ..."
make
CheckLastResult

echo "[+] VPN Installer: VPN Installer.app ..."
mkdir -p "bin/VPN Installer.app/Contents/Library/LaunchServices" || CheckLastResult
mkdir -p "bin/VPN Installer.app/Contents/MacOS" || CheckLastResult
cp "../helper/net.vpn.client.Helper" "bin/VPN Installer.app/Contents/Library/LaunchServices" || CheckLastResult
cp "bin/VPN Installer" "bin/VPN Installer.app/Contents/MacOS" || CheckLastResult
cp "etc/install.sh" "bin/VPN Installer.app/Contents/MacOS" || CheckLastResult
cp "VPN Installer-Info.plist" "bin/VPN Installer.app/Contents/Info.plist" || CheckLastResult

echo "[+] VPN Installer: VPN Uninstaller.app ..."
mkdir -p "bin/VPN Uninstaller.app/Contents/MacOS" || CheckLastResult
cp "bin/VPN Uninstaller" "bin/VPN Uninstaller.app/Contents/MacOS" || CheckLastResult
cp "VPN Uninstaller-Info.plist" "bin/VPN Uninstaller.app/Contents/Info.plist" || CheckLastResult

echo "[ ] VPN Installer: Done"
echo "    ${_SCRIPT_DIR}/bin/VPN Installer.app"
echo "    ${_SCRIPT_DIR}/bin/VPN Uninstaller.app"

cd ${_BASE_DIR}

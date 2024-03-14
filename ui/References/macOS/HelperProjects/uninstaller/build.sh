#!/bin/bash
source ./app.sh

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

if [ ! -f "../helper/${HelperID}.Helper" ]; then
  echo " File not exists '../helper/${HelperID}.Helper'. Please, compile helper project first."
  exit 1
fi

rm -fr bin
CheckLastResult

echo "[ ] *** Compiling VPN Installer / Uninstaller ***"

echo "[+] $App Installer: updating certificate info in .plist ..."
echo "    Apple DevID certificate: '${_SIGN_CERT}'"
plutil -replace CFBundleExecutable -string "$App Installer" "VPN Installer-Info.plist" || CheckLastResult
plutil -replace CFBundleIdentifier -string "$HelperID.installer" "VPN Installer-Info.plist" || CheckLastResult
plutil -replace SMPrivilegedExecutables -xml \
        "<dict> \
      		<key>${HelperID}.Helper</key> \
      		<string>identifier \"${HelperID}.Helper\" and anchor apple generic and certificate 1[field.1.2.840.113635.100.6.2.6] /* exists */ and certificate leaf[field.1.2.840.113635.100.6.1.13] /* exists */ and certificate leaf[subject.OU] = \"${_SIGN_CERT}\"</string> \
      	</dict>" "VPN Installer-Info.plist" || CheckLastResult
plutil -replace CFBundleExecutable -string "$App Uninstaller" "VPN Uninstaller-Info.plist" || CheckLastResult
plutil -replace CFBundleIdentifier -string "$HelperID.uninstaller" "VPN Uninstaller-Info.plist" || CheckLastResult
plutil -replace SMPrivilegedExecutables -xml \
        "<dict> \
          <key>${HelperID}.Helper</key> \
          <string>identifier \"${HelperID}.Helper\" and anchor apple generic and certificate 1[field.1.2.840.113635.100.6.2.6] /* exists */ and certificate leaf[field.1.2.840.113635.100.6.1.13] /* exists */ and certificate leaf[subject.OU] = \"${_SIGN_CERT}\"</string> \
        </dict>" "VPN Uninstaller-Info.plist" || CheckLastResult

echo "[+] $App Installer: make ..."
rm -fr bin
mkdir bin

cc \
  -D APP_NAME="\"${AppName}\"" \
  -D APP_SLUG="\"${App}\"" \
  -D HELPER_LABEL="\"${HelperID}.Helper\"" \
  -m64 -framework Foundation \
  -mmacosx-version-min=10.6 \
  -D IS_INSTALLER=0 \
  -framework ServiceManagement \
  -framework Security \
  -Xlinker -sectcreate -Xlinker __TEXT -Xlinker __info_plist -Xlinker "VPN Uninstaller-Info.plist" \
  uninstaller.c \
  -o "bin/$App Uninstaller"

cc \
  -D HELPER_LABEL="\"${HelperID}.Helper\"" \
  -D APP_NAME="\"${AppName}\"" \
  -D APP_SLUG="\"${App}\"" \
  -m64 -framework Foundation \
  -mmacosx-version-min=10.6 \
  -D IS_INSTALLER=1 \
  -framework ServiceManagement \
  -framework Security \
  -Xlinker -sectcreate -Xlinker __TEXT -Xlinker __info_plist -Xlinker "VPN Installer-Info.plist" \
  uninstaller.c \
  -o "bin/$App Installer"
CheckLastResult

echo "[+] $App Installer: $App Installer.app ..."
mkdir -p "bin/$App Installer.app/Contents/Library/LaunchServices" || CheckLastResult
mkdir -p "bin/$App Installer.app/Contents/MacOS" || CheckLastResult
cp "../helper/$HelperID.Helper" "bin/$App Installer.app/Contents/Library/LaunchServices" || CheckLastResult
cp "bin/$App Installer" "bin/$App Installer.app/Contents/MacOS" || CheckLastResult
cp "etc/install.sh" "bin/$App Installer.app/Contents/MacOS" || CheckLastResult
cp "VPN Installer-Info.plist" "bin/$App Installer.app/Contents/Info.plist" || CheckLastResult

echo "[+] $App Installer: $App Uninstaller.app ..."
mkdir -p "bin/$App Uninstaller.app/Contents/MacOS" || CheckLastResult
cp "bin/$App Uninstaller" "bin/$App Uninstaller.app/Contents/MacOS" || CheckLastResult
cp "VPN Uninstaller-Info.plist" "bin/$App Uninstaller.app/Contents/Info.plist" || CheckLastResult

echo "[ ] VPN Installer: Done"
echo "    ${_SCRIPT_DIR}/bin/$App Installer.app"
echo "    ${_SCRIPT_DIR}/bin/$App Uninstaller.app"

cd ${_BASE_DIR}

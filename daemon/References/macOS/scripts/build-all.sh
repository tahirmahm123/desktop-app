#!/bin/sh

cd "$(dirname "$0")"

VERSION=""
DATE="$(date "+%Y-%m-%d")"
COMMIT="$(git rev-list -1 HEAD)"

while getopts ":v:" opt; do
  case $opt in
    v) VERSION="$OPTARG"
    ;;
#    \?) echo "Invalid option -$OPTARG" >&2
#   ;;
  esac
done

echo "############################################"
echo "### Building VPN Daemon"
echo "### OpenVPN and WireGuard will be also recompiled if they are not exists"

if [ "$#" -eq 0 ]
then
  echo "### Possible arguments:"
  echo "###   -norebuild    - do not rebuild openVPN and WireGuard binaries is they already compiled"
  echo "###   -debug        - compile VPN Daemon in debug mode"
  echo "###   -libvpn      - use XPC listener for notifying clients about daemon connection port (latest VPN UI not using XPC)"
  echo "###   -wifi         - enable wifi support (do not ask 'Enable WIFI support?' question before demon build start)"
fi
echo "############################################"

if [[ ! -f "../_deps/openvpn_inst/bin/openvpn" ]] || [[ ! -f "../_deps/wg_inst/wg" ]] || [[ ! -f "../_deps/wg_inst/wireguard-go" ]]
then
  echo "Please, check/modify required versions at the begining of scripts:"
  if [[ ! -f "../_deps/openvpn_inst/bin/openvpn" ]]
  then
    echo "    build-openvpn.sh"
  fi

  if [[ ! -f "../_deps/wg_inst/wg" ]] || [[ ! -f "../_deps/wg_inst/wireguard-go" ]]
  then
    echo "    build-wireguard.sh"
  fi

  read -p "Press enter to start ..."
fi

# Exit immediately if a command exits with a non-zero status.
set -e

function BuildOpenVPN
{
  echo "############################################"
  echo "### OpenVPN"
  echo "############################################"
  ./build-openvpn.sh
}

function BuildWireGuard
{
  echo "############################################"
  echo "### WireGuard"
  echo "############################################"
  ./build-wireguard.sh
}

function BuildObfs4proxy
{
  echo "############################################"
  echo "### obfs4proxy"
  echo "############################################"
  ./build-obfs4proxy.sh
}

if [ ! -z "$GITHUB_ACTIONS" ]; then
  echo "! GITHUB_ACTIONS detected ! It is just a build test."
  echo "! Skipped compilation of third-party dependencies: OpenVPN, WireGuard, obfs4proxy !"
else
  if [[ "$@" == *"-norebuild"* ]]
  then
      # check if we need to compile openvpn
      if [[ ! -f "../_deps/openvpn_inst/bin/openvpn" ]]
      then
        echo "OpenVPN not compiled"
        BuildOpenVPN
      else
        echo "OpenVPN already compiled. Skipping build."
      fi

      # check if we need to compile WireGuard
      if [[ ! -f "../_deps/wg_inst/wg" ]] || [[ ! -f "../_deps/wg_inst/wireguard-go" ]]
      then
        echo "WireGuard not compiled"
        BuildWireGuard
      else
        echo "WireGuard already compiled. Skipping build."
      fi

      # check if we need to compile obfs4proxy
      if [[ ! -f "../_deps/obfs4proxy_inst/obfs4proxy" ]]
      then
        echo "obfs4proxy not compiled"
        BuildObfs4proxy
      else
        echo "obfs4proxy already compiled. Skipping build."
      fi

  else
    # recompile openvpn, WireGuard, obfs4proxy
    BuildOpenVPN
    BuildWireGuard
    BuildObfs4proxy
  fi
fi
# updating servers.json
./update-servers.sh

echo "======================================================"
echo "=============== VPN Agent ==========================="
echo "======================================================"
echo "Version: $VERSION"
echo "Date   : $DATE"
echo "Commit : $COMMIT"

cd ../../../

BUILDTAGS_DEBUG=""
BUILDTAGS_NOWIFI=""
BUILDTAGS_USE_LIBVPN=""

if [[ "$@" == *"-debug"* ]]
then
  BUILDTAGS_DEBUG="debug"
fi

if [[ "$@" == *"-libvpn"* ]]
then
  BUILDTAGS_USE_LIBVPN="libvpn"
fi

if [[ "$@" != *"-wifi"* ]]
then
  echo ""
  echo "Enable WIFI support?"
  echo "(this will lead to some additional library dependencies for the final binary)"
  read -p "[y\n]? (N - default): " yn
  case $yn in
      [Yy]* )
          ;;
      [Nn]* )
        BUILDTAGS_NOWIFI="nowifi"
        ;;
      * )
        BUILDTAGS_NOWIFI="nowifi"
        ;;
  esac
fi

CGO_CFLAGS=-mmacosx-version-min=10.10 CGO_LDFLAGS=-mmacosx-version-min=10.10 go build -tags "${BUILDTAGS_NOWIFI} ${BUILDTAGS_USE_LIBVPN} ${BUILDTAGS_DEBUG}" -o "VPN Agent" -trimpath -ldflags "-v -X github.com/tahirmahm123/vpn-desktop-app/daemon/version._version=$VERSION -X github.com/tahirmahm123/vpn-desktop-app/daemon/version._commit=$COMMIT -X github.com/tahirmahm123/vpn-desktop-app/daemon/version._time=$DATE"

echo "Cpmpiled daemon binary: '$(pwd)/VPN Agent'"

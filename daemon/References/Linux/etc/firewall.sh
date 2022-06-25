#!/bin/bash

#
#  Daemon for VPN Client Desktop
#  https://github.com/tahirmahm123/vpn-desktop-app/daemon
#
#  Created by Stelnykovych Alexandr.
#  Copyright (c) 2020 Privatus Limited.
#
#  This file is part of the Daemon for VPN Client Desktop.
#
#  The Daemon for VPN Client Desktop is free software: you can redistribute it and/or
#  modify it under the terms of the GNU General Public License as published by the Free
#  Software Foundation, either version 3 of the License, or (at your option) any later version.
#
#  The Daemon for VPN Client Desktop is distributed in the hope that it will be useful,
#  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
#  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
#  details.
#
#  You should have received a copy of the GNU General Public License
#  along with the Daemon for VPN Client Desktop. If not, see <https://www.gnu.org/licenses/>.
#

# Useful commands
#   List all rules: 
#     sudo iptables -L -v
#     or
#     sudo iptables -S

IPv4BIN=iptables
IPv6BIN=ip6tables

LOCKWAITTIME=2

# main chains for VPN firewall
IN_VPN=VPN-IN
OUT_VPN=VPN-OUT
# VPN chains for VPN dependend rules (applicable when VPN enabled)
IN_VPN_IF=VPN-IN-VPN
OUT_VPN_IF=VPN-OUT-VPN
# chain for non-VPN depended exceptios (applicable all time when firewall enabled)
# can be used, for example, for 'allow LAN' functionality
IN_VPN_STAT_EXP=VPN-IN-STAT-EXP
OUT_VPN_STAT_EXP=VPN-OUT-STAT-EXP
# chain for non-VPN depended exceptios: only for ICMP protocol (ping)
IN_VPN_ICMP_EXP=VPN-IN-ICMP-EXP
OUT_VPN_ICMP_EXP=VPN-OUT-ICMP-EXP

# returns 0 if chain exists
function chain_exists()
{
    local bin=$1
    local chain_name=$2
    ${bin} -w ${LOCKWAITTIME} -n -L ${chain_name} >/dev/null 2>&1
}

function create_chain()
{
  local bin=$1
  local chain_name=$2
  chain_exists ${bin} ${chain_name} || ${bin} -w ${LOCKWAITTIME} -N ${chain_name}
}

# Checks if the VPN Firewall is enabled
# 0 - if enabled
# 1 - if not enabled
function get_firewall_enabled {
  chain_exists ${IPv4BIN} ${OUT_VPN}
}

# Load rules
function enable_firewall {
    get_firewall_enabled

    if (( $? == 0 )); then
      echo "Firewall is already enabled. Please disable it first" >&2
      return 0
    fi

    set -e

    if [ -f /proc/net/if_inet6 ]; then
      ### IPv6 ###

      # IPv6: define chains
      create_chain ${IPv6BIN} ${IN_VPN}
      create_chain ${IPv6BIN} ${OUT_VPN}
      create_chain ${IPv6BIN} ${IN_VPN_IF}
      create_chain ${IPv6BIN} ${OUT_VPN_IF}

      # IPv6: allow  local (lo) interface
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -o lo -j ACCEPT
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -i lo -j ACCEPT

      # allow link-local addresses
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -s FE80::/10 -j ACCEPT
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -d FE80::/10 -j ACCEPT

      # allow unique-local addresses
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -s FD00::/8 -j ACCEPT
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -d FD00::/8 -j ACCEPT

      # allow DHCP port (547out 546in)
      # ${IPv6BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -p udp --dport 547 -j ACCEPT
      # ${IPv6BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -p udp --dport 546 -j ACCEPT

      # IPv6: assign our chains to global (global -> VPN_CHAIN -> VPN_VPN_CHAIN)
      ${IPv6BIN} -w ${LOCKWAITTIME} -A OUTPUT -j ${OUT_VPN}
      ${IPv6BIN} -w ${LOCKWAITTIME} -A INPUT -j ${IN_VPN}
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -j ${OUT_VPN_IF}
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -j ${IN_VPN_IF}

      # IPv6: block everything by default
      ${IPv6BIN} -w ${LOCKWAITTIME} -P INPUT DROP
      ${IPv6BIN} -w ${LOCKWAITTIME} -P OUTPUT DROP

    else
      echo "IPv6 disabled: skipping IPv6 rules"
    fi

    ### IPv4 ###
 
    # define chains
    create_chain ${IPv4BIN} ${IN_VPN}
    create_chain ${IPv4BIN} ${OUT_VPN}

    create_chain ${IPv4BIN} ${IN_VPN_IF}
    create_chain ${IPv4BIN} ${OUT_VPN_IF}

    create_chain ${IPv4BIN} ${IN_VPN_STAT_EXP}
    create_chain ${IPv4BIN} ${OUT_VPN_STAT_EXP}

    create_chain ${IPv4BIN} ${IN_VPN_ICMP_EXP}
    create_chain ${IPv4BIN} ${OUT_VPN_ICMP_EXP}

    # allow  local (lo) interface
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -o lo -j ACCEPT
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -i lo -j ACCEPT

    # allow DHCP port (67out 68in)
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -p udp --dport 67 -j ACCEPT
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -p udp --dport 68 -j ACCEPT

    # enable all ICMP ping outgoing request (needed to be able to ping VPN servers)
    #${IPv4BIN} -A ${OUT_VPN} -p icmp --icmp-type 8 -d 0/0 -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT
    #${IPv4BIN} -A ${IN_VPN} -p icmp --icmp-type 0 -s 0/0 -m state --state ESTABLISHED,RELATED -j ACCEPT

    # assign our chains to global 
    # (global -> VPN_CHAIN -> VPN_VPN_CHAIN)
    # (global -> VPN_CHAIN -> IN_VPN_STAT_EXP)
    ${IPv4BIN} -w ${LOCKWAITTIME} -A OUTPUT -j ${OUT_VPN}
    ${IPv4BIN} -w ${LOCKWAITTIME} -A INPUT -j ${IN_VPN}
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -j ${OUT_VPN_IF}
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -j ${IN_VPN_IF}
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -j ${OUT_VPN_STAT_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -j ${IN_VPN_STAT_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN} -j ${OUT_VPN_ICMP_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_VPN} -j ${IN_VPN_ICMP_EXP}

    # block everything by default
    ${IPv4BIN} -w ${LOCKWAITTIME} -P INPUT DROP
    ${IPv4BIN} -w ${LOCKWAITTIME} -P OUTPUT DROP

    set +e

    echo "VPN Firewall enabled"
}

# Remove all rules
function disable_firewall {
    # Flush rules and delete custom chains

    ### allow everything by default ###
    ${IPv4BIN} -w ${LOCKWAITTIME} -P INPUT ACCEPT
    ${IPv4BIN} -w ${LOCKWAITTIME} -P OUTPUT ACCEPT
    ${IPv6BIN} -w ${LOCKWAITTIME} -P INPUT ACCEPT
    ${IPv6BIN} -w ${LOCKWAITTIME} -P OUTPUT ACCEPT

    ### IPv4 ###
    # '-D' Delete matching rule from chain
    ${IPv4BIN} -w ${LOCKWAITTIME} -D OUTPUT -j ${OUT_VPN}
    ${IPv4BIN} -w ${LOCKWAITTIME} -D INPUT -j ${IN_VPN}
    ${IPv4BIN} -w ${LOCKWAITTIME} -D ${OUT_VPN} -j ${OUT_VPN_IF}
    ${IPv4BIN} -w ${LOCKWAITTIME} -D ${IN_VPN} -j ${IN_VPN_IF}
    ${IPv4BIN} -w ${LOCKWAITTIME} -D ${OUT_VPN} -j ${OUT_VPN_STAT_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -D ${IN_VPN} -j ${IN_VPN_STAT_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -D ${OUT_VPN} -j ${OUT_VPN_ICMP_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -D ${IN_VPN} -j ${IN_VPN_ICMP_EXP}
    # '-F' Delete all rules in  chain or all chains
    ${IPv4BIN} -w ${LOCKWAITTIME} -F ${OUT_VPN_IF}
    ${IPv4BIN} -w ${LOCKWAITTIME} -F ${IN_VPN_IF}
    ${IPv4BIN} -w ${LOCKWAITTIME} -F ${OUT_VPN}
    ${IPv4BIN} -w ${LOCKWAITTIME} -F ${IN_VPN}
    ${IPv4BIN} -w ${LOCKWAITTIME} -F ${OUT_VPN_STAT_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -F ${IN_VPN_STAT_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -F ${OUT_VPN_ICMP_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -F ${IN_VPN_ICMP_EXP}
    # '-X' Delete a user-defined chain
    ${IPv4BIN} -w ${LOCKWAITTIME} -X ${OUT_VPN_IF}
    ${IPv4BIN} -w ${LOCKWAITTIME} -X ${IN_VPN_IF}
    ${IPv4BIN} -w ${LOCKWAITTIME} -X ${OUT_VPN}
    ${IPv4BIN} -w ${LOCKWAITTIME} -X ${IN_VPN}
    ${IPv4BIN} -w ${LOCKWAITTIME} -X ${OUT_VPN_STAT_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -X ${IN_VPN_STAT_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -X ${OUT_VPN_ICMP_EXP}
    ${IPv4BIN} -w ${LOCKWAITTIME} -X ${IN_VPN_ICMP_EXP}

    ### IPv6 ###
    ${IPv6BIN} -w ${LOCKWAITTIME} -D OUTPUT -j ${OUT_VPN}
    ${IPv6BIN} -w ${LOCKWAITTIME} -D INPUT -j ${IN_VPN}
    ${IPv6BIN} -w ${LOCKWAITTIME} -D ${OUT_VPN} -j ${OUT_VPN_IF}
    ${IPv6BIN} -w ${LOCKWAITTIME} -D ${IN_VPN} -j ${IN_VPN_IF}

    ${IPv6BIN} -w ${LOCKWAITTIME} -F ${OUT_VPN_IF}
    ${IPv6BIN} -w ${LOCKWAITTIME} -F ${IN_VPN_IF}
    ${IPv6BIN} -w ${LOCKWAITTIME} -F ${OUT_VPN}
    ${IPv6BIN} -w ${LOCKWAITTIME} -F ${IN_VPN}

    ${IPv6BIN} -w ${LOCKWAITTIME} -X ${OUT_VPN_IF}
    ${IPv6BIN} -w ${LOCKWAITTIME} -X ${IN_VPN_IF}
    ${IPv6BIN} -w ${LOCKWAITTIME} -X ${OUT_VPN}
    ${IPv6BIN} -w ${LOCKWAITTIME} -X ${IN_VPN}
    echo "VPN Firewall disabled"
}

function client_connected {
  IFACE=$1

  # allow all packets to VPN interface
  ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN_IF} -o ${IFACE} -j ACCEPT
  ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_VPN_IF} -i ${IFACE} -j ACCEPT

  if [ -f /proc/net/if_inet6 ]; then
      ### IPv6 ###

      # allow all packets to VPN interface
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${OUT_VPN_IF} -o ${IFACE} -j ACCEPT
      ${IPv6BIN} -w ${LOCKWAITTIME} -A ${IN_VPN_IF} -i ${IFACE} -j ACCEPT
    fi
}

function client_disconnected {
  ${IPv4BIN} -w ${LOCKWAITTIME} -F ${OUT_VPN_IF}
  ${IPv4BIN} -w ${LOCKWAITTIME} -F ${IN_VPN_IF}

  if [ -f /proc/net/if_inet6 ]; then
    ### IPv6 ###
    ${IPv6BIN} -w ${LOCKWAITTIME} -F ${OUT_VPN_IF}
    ${IPv6BIN} -w ${LOCKWAITTIME} -F ${IN_VPN_IF}
  fi
}

function add_exceptions {
  IN_CH=$1
  OUT_CH=$2
  shift 2
  EXP=$@

  create_chain ${IPv4BIN} ${IN_CH}
  create_chain ${IPv4BIN} ${OUT_CH}

  #add new rule
  # '-C' option is checking if the rule already exists (needed to avoid duplicates)
  ${IPv4BIN} -w ${LOCKWAITTIME} -C ${IN_CH} -s $@ -j ACCEPT || ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_CH} -s $@ -j ACCEPT
  ${IPv4BIN} -w ${LOCKWAITTIME} -C ${OUT_CH} -d $@ -j ACCEPT || ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_CH} -d $@ -j ACCEPT
}

function add_direction_exception {
  IN_CH=$1
  OUT_CH=$2

  #SRC_PORT=$3
  DST_ADDR=$4
  DST_PORT=$5
  PROTOCOL=$6

  create_chain ${IPv4BIN} ${IN_CH}
  create_chain ${IPv4BIN} ${OUT_CH}

  #add new rule
  # ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_CH}  -s ${DST_ADDR} -p ${PROTOCOL} --sport ${DST_PORT} --dport ${SRC_PORT} -j ACCEPT
  # ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_CH} -d ${DST_ADDR} -p ${PROTOCOL} --sport ${SRC_PORT} --dport ${DST_PORT} -j ACCEPT
  ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_CH}  -s ${DST_ADDR} -p ${PROTOCOL} --sport ${DST_PORT} -j ACCEPT
  ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_CH} -d ${DST_ADDR} -p ${PROTOCOL} --dport ${DST_PORT} -j ACCEPT
}

function remove_exceptions {
  IN_CH=$1
  OUT_CH=$2
  shift 2
  EXP=$@

  ${IPv4BIN} -w ${LOCKWAITTIME} -D ${IN_CH} -s $@ -j ACCEPT
  ${IPv4BIN} -w ${LOCKWAITTIME} -D ${OUT_CH} -d $@ -j ACCEPT
}

function remove_exceptions_icmp {
  IN_CH=$1
  OUT_CH=$2
  shift 2
  EXP=$@

  ${IPv4BIN} -w ${LOCKWAITTIME} -D ${IN_CH} -p icmp --icmp-type 0 -s $@ -m state --state ESTABLISHED,RELATED -j ACCEPT
  ${IPv4BIN} -w ${LOCKWAITTIME} -D ${OUT_CH} -p icmp --icmp-type 8 -d $@ -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT
}

function add_exceptions_icmp {
  IN_CH=$1
  OUT_CH=$2
  shift 2
  EXP=$@

  create_chain ${IPv4BIN} ${IN_CH}
  create_chain ${IPv4BIN} ${OUT_CH}

  # remove same rule if exists (just to avoid duplicates)
  ${IPv4BIN} -w ${LOCKWAITTIME} -D ${IN_CH} -p icmp --icmp-type 0 -s $@ -m state --state ESTABLISHED,RELATED -j ACCEPT
  ${IPv4BIN} -w ${LOCKWAITTIME} -D ${OUT_CH} -p icmp --icmp-type 8 -d $@ -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT

  #add new rule
  ${IPv4BIN} -w ${LOCKWAITTIME} -A ${IN_CH} -p icmp --icmp-type 0 -s $@ -m state --state ESTABLISHED,RELATED -j ACCEPT
  ${IPv4BIN} -w ${LOCKWAITTIME} -A ${OUT_CH} -p icmp --icmp-type 8 -d $@ -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT
}

function main {

    if [[ $1 = "-enable" ]] ; then

      enable_firewall

    elif [[ $1 = "-disable" ]] ; then

      disable_firewall

    elif [[ $1 = "-status" ]] ; then

      get_firewall_enabled

      if (( $? == 0 )); then
        echo "VPN Firewall is enabled"
        return 0
      else
        echo "VPN Firewall is disabled"
        return 1
      fi

    elif [[ $1 = "-add_exceptions" ]]; then

      shift
      add_exceptions ${IN_VPN_IF} ${OUT_VPN_IF} $@

    elif [[ $1 = "-remove_exceptions" ]]; then

      shift
      remove_exceptions ${IN_VPN_IF} ${OUT_VPN_IF} $@

    elif [[ $1 = "-add_exceptions_static" ]]; then
      
      shift
      add_exceptions ${IN_VPN_STAT_EXP} ${OUT_VPN_STAT_EXP} $@

    elif [[ $1 = "-remove_exceptions_static" ]]; then

      shift
      remove_exceptions ${IN_VPN_STAT_EXP} ${OUT_VPN_STAT_EXP} $@

    elif [[ $1 = "-add_exceptions_icmp" ]]; then

      shift
      add_exceptions_icmp ${IN_VPN_ICMP_EXP} ${OUT_VPN_ICMP_EXP} $@

    elif [[ $1 = "-remove_exceptions_icmp" ]]; then

      shift
      remove_exceptions_icmp ${IN_VPN_ICMP_EXP} ${OUT_VPN_ICMP_EXP} $@

    elif [[ $1 = "-connected" ]]; then
        IFACE=$2
        #SRC_ADDR=$3
        SRC_PORT=$4
        DST_ADDR=$5
        DST_PORT=$6
        PROTOCOL=$7

        # allow all communication trough vpn interface
        client_connected ${IFACE}

        # allow communication with host only srcPort <=> host.dstsPort
        add_direction_exception ${IN_VPN_IF} ${OUT_VPN_IF} ${SRC_PORT} ${DST_ADDR} ${DST_PORT} ${PROTOCOL}
    elif [[ $1 = "-disconnected" ]]; then
        shift
        client_disconnected
    else
        echo "Unknown command"
        return 2
    fi
}

main $@

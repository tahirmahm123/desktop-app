#!/bin/sh
sudo launchctl unload /Library/LaunchDaemons/net.vpn.client.Helper.plist
sudo rm /Library/LaunchDaemons/net.vpn.client.Helper.plist
sudo rm /Library/PrivilegedHelperTools/net.vpn.client.Helper

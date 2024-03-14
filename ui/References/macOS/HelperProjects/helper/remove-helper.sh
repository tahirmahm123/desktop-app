#!/bin/sh
source ../../app.sh
sudo launchctl unload /Library/LaunchDaemons/${HelperID}.Helper.plist
sudo rm /Library/LaunchDaemons/${HelperID}.Helper.plist
sudo rm /Library/PrivilegedHelperTools/${HelperID}.Helper

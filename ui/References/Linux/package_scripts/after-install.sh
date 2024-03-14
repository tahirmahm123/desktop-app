#!/bin/sh

echo "[*] After install (<%= version %> : <%= pkg %> : $1)"

DIR=/usr/share/applications
DESKTOP_FILE=$DIR/VPN.desktop
if [ -d "$DIR" ]; then
    echo "[ ] Installing .desktop file..."
    ln -fs /opt/vpn/ui/VPN.desktop $DESKTOP_FILE || echo "[!] Failed to create .desktop file: '$DESKTOP_FILE'"
else
    echo "[!] Unable to install .desktop file. Folder '$DIR' not exists"
fi

sudo chmod 4755 /opt/vpn/ui/bin/chrome-sandbox || echo "[!] Failed to 'chmod' for '/opt/vpn/ui/bin/chrome-sandbox'"
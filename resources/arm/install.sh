#!/usr/bin/env bash
if [ $(id -u) -ne 0 ]; then
	echo "Please run as root"
	exit 1
fi

INSTALL_DIR=/usr/local/go-metar-blink

echo "Installing go-metar-blink..."

mkdir $INSTALL_DIR || echo "Failed to make install dir"

# Copy this dist
cp -r ./ $INSTALL_DIR
cp go-metar-blink.service /etc/systemd/system/go-metar-blink.service

cp -n settings.json /boot/go-metar-blink.settings.example.json

# Install the service
systemctl enable go-metar-blink;
systemctl start go-metar-blink;

# Test if the service started successfully
systemctl is-active --quiet go-metar-blink;
if [[ $? -eq 0 ]]; then
	echo "Go Metar Blink service installed successfully!";
else
	echo "Service failed to start use `journalctl` to see systemctl logs.";
fi

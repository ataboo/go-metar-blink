#!/usr/bin/env bash
if [ $(id -u) -ne 0 ]; then
	echo "This script needs elevation to run"
	exit 1
fi

echo "Uninstalling Go Metar Blink..."

systemctl stop go-metar-blink
systemctl disable go-metar-blink

rm -r /usr/local/go-metar-blink || echo "Failed to delete app dir";
rm /etc/systemd/system/go-metar-blink.service || echo "Failed to delete service";
rm /boot/go-metar-blink.settings.example.json || "Failed to delete example settings"
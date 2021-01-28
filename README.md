# Go Metar Blink

Controls a map that shows aviation weather conditions through the colour and pattern of addressable LEDs.  The LEDs are controlled
by a Raspberry Pi with a level-shifting inverter.

## Build

The `dist` directory has a tar'd copy for arm6 that can be installed.  I may pack other platforms in the future but for now I'm only
aiming to use it on a Raspberry Pi Zero.

To build yourself, create a docker build container as shown in documentation/build-docker.md.  You should be able to build one of ws-2811's example files.

Run `resources/build.sh` from the project root to pack a tar file in `dist`.  
 
## Install
 
Copy and unpack the tar on the pi and run `install.sh` to setup the service.  You will likely need to elevate permissions as this script copies to `/usr/local`.

## Configuration

If `/boot/go-metar-blink.settings.json` is present, it will override the settings in `/usr/local/go-metar-blink/settings.json`.

## Uninstall

Run `uninstall.sh` to remove the service and delete the program from `/usr/local/go-metar-blink`.

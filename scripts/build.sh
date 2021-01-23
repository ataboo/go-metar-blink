# !/bin/bash

echo $PWD

APP="go-metar-blink"
OUTPUT_BIN=$APP-armv6

rm -rf ./build
mkdir -p build/$APP

cp settings.json ./build/$APP
cp -r resources/arm/* ./build/$APP

docker run --rm -v "$PWD":/usr/src/$APP --platform linux/arm/v6 -w /usr/src/$APP ws2811-builder:latest go build -o "./build/$APP/$OUTPUT_BIN" -v

tar -czvf $OUTPUT_BIN.tar -C build $APP

mv $OUTPUT_BIN.tar ./dist

echo "Built to dist/$OUTPUT_BIN.tar"

#!/bin/bash


name="icon"

echo "Processing icon ..."
set -e

if [ ! -f icon.svg ]; then
    echo "Error: icon.svg not found"
    exit 1
fi

if  ! command -v rsvg-convert &>/dev/null ; then
    echo "Error: rsvg-convert not found. please install librsvg and try again"
    exit 1
fi

mkdir -p icon.iconset

# run inkscape from the command line to generate the iconset formatted for icns
rsvg-convert --width=16 --height=16 icon.svg -o icon.iconset/icon_16x16.png
rsvg-convert --width=32 --height=32 icon.svg -o icon.iconset/icon_16x16@2x.png
rsvg-convert --width=32 --height=32 icon.svg -o icon.iconset/icon_32x32.png
rsvg-convert --width=64 --height=64 icon.svg -o icon.iconset/icon_32x32@2x.png
rsvg-convert --width=128 --height=128 icon.svg -o icon.iconset/icon_128x128.png
rsvg-convert --width=256 --height=256 icon.svg -o icon.iconset/icon_128x128@2x.png
rsvg-convert --width=256 --height=256 icon.svg -o icon.iconset/icon_256x256.png
rsvg-convert --width=512 --height=512 icon.svg -o icon.iconset/icon_256x256@2x.png
rsvg-convert --width=1024 --height=1024 icon.svg -o icon.iconset/icon_512x512.png
rsvg-convert --width=2048 --height=2048 icon.svg -o icon.iconset/icon_512x512@2x.png

# run osx iconutil app to convert the iconset to icns format
iconutil -c icns "icon.iconset"

#cleanup
rm -R "icon.iconset"


echo "Done."

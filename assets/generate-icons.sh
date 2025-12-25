#!/bin/bash

# This script generates simple emoji-based PNG icons for the tray
# Using macOS's ability to render emoji to images

# Check if ImageMagick is available
if ! command -v convert &> /dev/null; then
    echo "ImageMagick not found. Installing via brew..."
    brew install imagemagick
fi

# Create icons with emojis
# Work session - tomato emoji 🍅
convert -size 64x64 xc:transparent -font "Apple Color Emoji" -pointsize 48 \
    -gravity center -annotate +0+0 "🍅" icon-work.png

# Short break - coffee emoji ☕
convert -size 64x64 xc:transparent -font "Apple Color Emoji" -pointsize 48 \
    -gravity center -annotate +0+0 "☕" icon-short-break.png

# Long break - star emoji 🌟
convert -size 64x64 xc:transparent -font "Apple Color Emoji" -pointsize 48 \
    -gravity center -annotate +0+0 "🌟" icon-long-break.png

# Paused - pause emoji ⏸️
convert -size 64x64 xc:transparent -font "Apple Color Emoji" -pointsize 48 \
    -gravity center -annotate +0+0 "⏸️" icon-paused.png

echo "Icons generated successfully!"
ls -lh icon-*.png

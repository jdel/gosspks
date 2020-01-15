#!/bin/bash -ue
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [[ "$(uname -m)" == *"64"* ]]; then ARCH="amd64"; else ARCH="386"; fi

curl -o /usr/local/bin/gosspks -L "https://jdel.org/gosspks/releases/download/v0.3.0/gosspks-$PLATFORM-$ARCH"
chmod +x /usr/local/bin/gosspks
echo "Installed to /usr/local/bin/gosspks"
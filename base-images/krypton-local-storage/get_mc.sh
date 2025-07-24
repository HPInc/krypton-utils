#!/bin/sh
set -e
HTTP_OK=200
#URL="https://dl.min.io/client/mc/release/linux-amd64/mc"
URL=https://dl.min.io/client/mc/release/linux-amd64/archive/mc.RELEASE.2020-06-20T00-18-43Z
DEST=/usr/local/bin/mc

result=$(curl --write-out "%{http_code}" -sLo "$DEST" "$URL")
if [ "$result" -ne "$HTTP_OK" ]; then
	echo "fs client download failed. stopping."
	exit 1
fi
chmod +x "$DEST"

#!/bin/sh
set -e

SEMVER=0.0.1
FILES="bin/krypton-cli config/config.yaml install.sh tools/krypton-cli-completion.sh"

make

tar -czvf krypton-cli-$SEMVER.tar.gz $FILES

#!/bin/sh
BIN=bin/krypton-cli
KRYPTON_HOME="$HOME/.krypton-cli"
CONFIG="$KRYPTON_HOME/config.yaml"
CONFIG_AWS_ROOT="$KRYPTON_HOME/aws_root.cert"
AWS_ROOT_CERT_URL="https://www.amazontrust.com/repository/AmazonRootCA1.pem"
# completion
COMPLETION=tools/krypton-cli-completion.sh
COMPLETION_DEST=/usr/share/bash-completion/completions/krypton-cli

set -e

if ! $(cmp -s "$BIN" /usr/local/bin/krypton-cli ); then
  echo "updating krypton-cli"
  sudo cp "$BIN" /usr/local/bin
else
  echo "krypton-cli is already at the latest version"
fi

# ensure dir
mkdir -p "$KRYPTON_HOME"

if ! $(cmp -s "$CONFIG" config/config.yaml); then
  # backup existing config if any
  if [ -e "$CONFIG" ]; then
    echo "making a backup of current config"
    cp "$CONFIG" "$CONFIG.backup"
  fi
  # copy new config
  cp config/config.yaml "$CONFIG"
fi

# pull aws iot cert if it does not exist
if [ ! -e "$CONFIG_AWS_ROOT" ]; then
  echo "downloading aws root cert"
  curl -sLo "$CONFIG_AWS_ROOT" "$AWS_ROOT_CERT_URL"
fi

# install bash autocompletion script if there are changes
if ! $(cmp -s "$COMPLETION" "$COMPLETION_DEST"); then
  echo "installing bash auto-completion, please restart session for auto-complete changes"
  sudo cp "$COMPLETION" "$COMPLETION_DEST"
fi

echo "done!"

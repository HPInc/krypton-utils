#!/bin/sh
LINTER_IMAGE=ghcr.io/hpinc/krypton/golangci-lint
DIR=${1:-$(pwd)}

# flag overrides
# excluding S1034. See https://staticcheck.io/docs/checks#S1034
docker run --rm \
  -v "$DIR":/service \
  -w /service \
  "$LINTER_IMAGE" golangci-lint run \
  -v -eS1034 --tests=false --timeout=5m

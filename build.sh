#!/bin/bash
set -e
cd $(dirname $0)

PRONAME=$(basename $(pwd))
VERSION="$(git rev-parse --short HEAD)-$(date '+%Y%m%d')"

echo "set version: ${VERSION}"
go build -ldflags "-X main.buildVersion=${VERSION}" -o ./bin/$PRONAME *.go
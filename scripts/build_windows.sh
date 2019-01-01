#!/bin/bash

# get script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

source $DIR/build_common.sh

# build windows
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0
go build -o "$BINDIR/wrangle-windows-amd64.exe" -ldflags "$ldflags"
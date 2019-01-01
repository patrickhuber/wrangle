#!/bin/bash

# get script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

source $DIR/build_common.sh

export GOOS=darwin
export GOARCH=amd64
go build -o "$BINDIR/wrangle-darwin-amd64" -ldflags "$ldflags"
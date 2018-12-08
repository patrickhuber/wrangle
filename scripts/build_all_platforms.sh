#!/bin/bash

# get script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

# set the output directory
BINDIR=$DIR/../bin

## load the version
version=$(cat $DIR/../version.dat)

## create the ldflags
ldflags="-X main.version=$version"

# build MacOS
export GOOS=darwin
export GOARCH=amd64
go build -o "$BINDIR/wrangle-darwin-amd64" -ldflags "$ldflags"

# build linux
export GOOS=linux
export GOARCH=amd64
go build -o "$BINDIR/wrangle-linux-amd64" -ldflags "$ldflags"

# build windows
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0
go build -o "$BINDIR/wrangle-windows-amd64.exe" -ldflags "$ldflags"
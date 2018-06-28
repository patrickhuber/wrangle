#!/bin/bash

# build MacOS
GOOS=darwin
GOARCH=amd64
go build -o bin/cli-mgr-darwin-amd64 main.go

# build linux
GOOS=linux
GOARCH=amd64
go build -o bin/cli-mgr-linux-amd64 main.go

# build windows
GOOS=windows
GOARCH=amd64
go build -o bin/cli-mgr-windows-amd64.exe main.go
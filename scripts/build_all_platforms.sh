#!/bin/bash

# build MacOS
export GOOS=darwin
export GOARCH=amd64
go build -o bin/wrangle-darwin-amd64 main.go

# build linux
export GOOS=linux
export GOARCH=amd64
go build -o bin/wrangle-linux-amd64 main.go

# build windows
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0
go build -o bin/wrangle-windows-amd64.exe main.go
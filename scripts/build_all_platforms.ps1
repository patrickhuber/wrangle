$ErrorActionPreference = "Stop"

$datFilePath = Join-Path $PSScriptRoot "../version.dat"
$version = get-content $datFilePath

$bin = Join-Path $PSScriptRoot "../bin"

$ldflags = "-X main.version=$version"
$ldflags

# build MacOS
$ENV:GOOS="darwin"
$ENV:GOARCH="amd64"
go build -o "$bin/wrangle-darwin-amd64" -ldflags "$ldflags"

# build linux
$ENV:GOOS="linux"
$ENV:GOARCH="amd64"
go build -o "$bin/wrangle-linux-amd64" -ldflags "$ldflags"

# build windows
$ENV:GOOS="windows"
$ENV:GOARCH="amd64"
go build -o "$bin/wrangle-windows-amd64.exe" -ldflags "$ldflags"
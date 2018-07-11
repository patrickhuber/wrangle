$ErrorActionPreference = "Stop"

# build MacOS
$ENV:GOOS="darwin"
$ENV:GOARCH="amd64"
go build -o bin/wrangle-darwin-amd64 main.go

# build linux
$ENV:GOOS="linux"
$ENV:GOARCH="amd64"
go build -o bin/wrangle-linux-amd64 main.go

# build windows
$ENV:GOOS="windows"
$ENV:GOARCH="amd64"
go build -o bin/wrangle-windows-amd64.exe main.go
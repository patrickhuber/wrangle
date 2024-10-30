$ErrorActionPreference = "Stop"
$dir=$PSScriptRoot
$root=Resolve-Path(Join-Path $dir "..")
$dockerfile = Join-Path $root "docker" "Dockerfile.build.linux"
$dockerfile
docker build --build-arg LOCAL="true" -f "$dockerfile" . 
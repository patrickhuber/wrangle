$ErrorActionPreference = "Stop"
$dir=$PSScriptRoot
$root=Resolve-Path(Join-Path $dir "..")
$dockerfile=Join-Path $root "Dockerfile.test.linux"
docker build --build-arg LOCAL="true" -f $dockerfile . 
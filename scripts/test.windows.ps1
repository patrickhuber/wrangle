$ErrorActionPreference = "Stop"
$dir=$PSScriptRoot
$root=Resolve-Path(Join-Path $dir "..")
$dockerfile=Join-Path $root "Dockerfile.test.windows"
docker build --build-arg LOCAL="true" -f $dockerfile . 
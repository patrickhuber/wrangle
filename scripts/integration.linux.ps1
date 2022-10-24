$ErrorActionPreference = "Stop"
docker build --build-arg LOCAL="true" -f Dockerfile.integration.linux . 
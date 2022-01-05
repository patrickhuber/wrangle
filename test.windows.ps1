$ErrorActionPreference = "Stop"
docker build --build-arg LOCAL="true" -f Dockerfile.test.windows . 
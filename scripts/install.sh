#!/bin/bash

# exit when any command fails
set -e

# determine if linux or darwin
OS=$(uname -s)
case $OS in
Linux) export PLATFORM=linux;;
Darwin) export PLATFORM=darwin;;
*) echo "expected OSTYPE 'linux-gnu' or 'darwin'. found $OSTYPE. no installer is available for this OSTYPE"; exit 1 ;;
esac

echo "platform: '$PLATFORM'"

# determine the architecture
ARCH=$(uname -m)
case $ARCH in
x86_64) export ARCH=amd64;;
*) echo "expected 'uname -m' x86_64. found $ARCH. no installer is available for this architecture"; exit 1 ;;
esac

echo "architecture: '$ARCH'"

# variables
export VERSION=0.10.0
export ARCHIVE=wrangle-${PLATFORM}-${ARCH}.tgz
export URL=https://github.com/patrickhuber/wrangle/releases/download/${VERSION}/${ARCHIVE}
echo "downloading: '$ARCHIVE' from '$URL'"

# download the cli
wget $URL

# extract the executable
# remove the file
echo "extracting: ${ARCHIVE}"
tar -xfz ${ARCHIVE}

echo "cleanup: ${ARCHIVE}"
rm ${ARCHIVE}

# create the global configuration
echo "installing"
sudo wrangle bootstrap

# cleanup
echo "installing"
rm wrangle
#!/bin/sh

# determine if linux or darwin
OS=$(uname -s)
case $OS in
Linux) export PLATFORM=linux;;
Darwin) export PLATFORM=darwin;;
*) echo "expected OSTYPE 'linux-gnu' or 'darwin'. found $OSTYPE. no installer is available for this OSTYPE"; exit 1 ;;
esac

# determine the architecture
ARCH=$(uname -m)
case $ARCH in
x86_64) export ARCH=amd64;;
*) echo "expected 'uname -m' x86_64. found $ARCH. no installer is available for this architecture"; exit 1 ;;
esac

# variables
export VERSION=0.10.0
export ARCHIVE=wrangle-${PLATFORM}-${ARCHITECTURE}.tgz

# download the cli
wget https://github.com/patrickhuber/wrangle/releases/download/${VERSION}/${ARCHIVE}

# extract the executable
# remove the file
tar -xfz ${ARCHIVE}
rm ${ARCHIVE}

# create the global configuration
sudo wrangle bootstrap

# cleanup
rm wrangle
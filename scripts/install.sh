#!/bin/bash

# exit when any command fails
set -e

# check for jq
if ! command -v jq &> /dev/null
then
    echo "jq must be installed"
    exit 1
fi

# check for curl
if ! command -v curl &> /dev/null
then 
    echo "curl must be installed"
    exit 1
fi

# determine if linux or darwin
OS=$(uname -s)
case $OS in
Linux) export PLATFORM=linux;;
Darwin) export PLATFORM=darwin;;
*) echo "expected OS TYPE 'linux-gnu' or 'darwin'. found $OS. no installer is available for this OSTYPE"; exit 1 ;;
esac

echo "platform: '$PLATFORM'"

# determine the architecture
ARCH=$(uname -m)
case $ARCH in
x86_64) export ARCH=amd64;;
*) echo "expected 'uname -m' x86_64. found $ARCH. no installer is available for this architecture"; exit 1 ;;
esac

echo "architecture: '$ARCH'"

# get the latest version
echo "getting latest wrangle version from github"
json=$(curl 'https://api.github.com/repos/patrickhuber/wrangle/releases/latest')

export VERSION=$(echo $json | jq -r '.tag_name')

# variables
export ARCHIVE=wrangle-${VERSION}-${PLATFORM}-${ARCH}.tar.gz
export ARCHIVE_FOLDER=wrangle-${VERSION}-${PLATFORM}-${ARCH}
export URL=https://github.com/patrickhuber/wrangle/releases/download/${VERSION}/${ARCHIVE}
echo "downloading: '$ARCHIVE' from '$URL'"

# download the cli
wget $URL

# extract the executable
# remove the file
echo "extracting: ${ARCHIVE}"
mkdir -p ${ARCHIVE_FOLDER}
tar xvzf ${ARCHIVE} -C ${ARCHIVE_FOLDER} 

# create the global configuration and install packages
echo "installing"
sudo ${ARCHIVE_FOLDER}/wrangle bootstrap

# cleanup
echo "cleanup: ${ARCHIVE_FOLDER}"
rm -rf ${ARCHIVE_FOLDER}

echo "cleanup: ${ARCHIVE}"
rm -f ${ARCHIVE}
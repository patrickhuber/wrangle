#!/bin/bash

# get script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
BIN=$(realpath $DIR/../bin)
pushd $BIN
    tar czvf wrangle-darwin-amd64.tgz wrangle-darwin-amd64
    tar czvf wrangle-windows-amd64.tgz wrangle-windows-amd64.exe
    tar czvf wrangle-linux-amd64.tgz wrangle-linux-amd64
    zip -9 wrangle-windows-amd64.zip wrangle-windows-amd64.exe
popd
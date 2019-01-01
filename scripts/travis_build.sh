#!/bin/bash

# get script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    $DIR/build_darwin.sh
fi

if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
    $DIR/build_linux.sh
    $DIR/build_windows.sh
fi
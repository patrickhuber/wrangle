#!/bin/bash

# get script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

source $DIR/build_darwin.sh
source $DIR/build_linux.sh
source $DIR/build_windows.sh
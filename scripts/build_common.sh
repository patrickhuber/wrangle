#!/bin/bash

# get script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

# set the output directory
BINDIR=$DIR/../bin

## load the version
version=$(cat $DIR/../version.dat)

## create the ldflags
ldflags="-X main.version=$version"
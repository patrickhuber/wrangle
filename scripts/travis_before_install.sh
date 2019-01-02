#!/bin/bash
if [[ $TRAVIS_OS_NAME == 'osx' ]]; then
  brew update;
  brew install gnu-getopt;
  brew upgrade gnupg;
fi
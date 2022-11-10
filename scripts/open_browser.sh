#!/usr/bin/env bash

#if [ -z "$1" ];
#then
#    echo "$1: can't open address - missing address arg"
#    exit 1
#fi
#
#address_to_open="$1";
address_to_open=http://localhost:3001

case "$OSTYPE" in
  solaris*) open "$address_to_open" ;;
  darwin*)  open "$address_to_open" ;;
  linux*)   xdg-open "$address_to_open" ;;
  msys*)    start "$address_to_open" ;;
  cygwin*)  start "$address_to_open" ;;
  *)        echo "unknown: $OSTYPE"; exit 1 ;;
esac

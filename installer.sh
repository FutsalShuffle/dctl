#!/bin/bash

architecture=""
ostype=""
case $(uname -m) in
    x86_64) architecture="amd64" ;;
    arm)    dpkg --print-architecture | grep -q "arm64" && architecture="arm64" || architecture="arm" ;;
esac

case $OSTYPE in
    linux*)     ostype="linux" ;;
    darwin*)    ostype="darwin" ;;
#    win32)      ostype="windows" ;;
#    windows)    ostype="windows" ;;
#    msys)       ostype="windows" ;;
#    cygwin)     ostype="windows" ;;
esac

if [ "$architecture" == "" ];
  then
    echo "Unable to determine system architecture.";
    exit 1;
fi

if [ "$ostype" == "" ];
  then
    echo "Unable to determine system OS.";
    exit 1;
fi

url=https://github.com/FutsalShuffle/dctl/releases/download/v0.6/dctl_"$architecture"_"$ostype"

echo "$url"

wget -O dctl "$url"  \
&& chmod +x dctl \
&& mkdir -p ~/.dctl/ \
&& mv dctl ~/.dctl/ \
&& sudo ln -s ~/.dctl/dctl /usr/local/bin/dctl

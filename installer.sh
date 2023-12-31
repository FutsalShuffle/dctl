#!/bin/bash

architecture=""
ostype=""
case $(uname -m) in
    x86_64) architecture="amd64" ;;
    arm64)  architecture="arm64" ;;
esac

case $OSTYPE in
    linux*)     ostype="linux" ;;
    darwin*)    ostype="darwin" ;;
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

#newUrl=""
#newUrl=$(curl -s https://api.github.com/repos/FutsalShuffle/dctl/releases/latest | grep -o -E 'browser_download_url\": \"(.+)\/dctl_amd64_linux')
#echo $newUrl
url=https://github.com/FutsalShuffle/dctl/releases/download/v1.0/dctl_"$architecture"_"$ostype"

echo "$url"

wget -O dctl "$url"  \
&& chmod +x dctl \
&& mkdir -p ~/.dctl/ \
&& mv dctl ~/.dctl/ \
&& sudo ln -s ~/.dctl/dctl /usr/local/bin/dctl

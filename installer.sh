#!/bin/bash
wget https://github.com/FutsalShuffle/dctl/releases/download/v0.5/dctl \
&& chmod +x dctl \
&& mkdir ~/.dctl/ \
&& mv dctl ~/.dctl/ \
&& sudo ln -s ~/.dctl/dctl /usr/local/bin/dctl

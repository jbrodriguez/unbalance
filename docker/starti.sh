#!/bin/bash

VERSION="latest"

if [ -n "$1" ]
then
	VERSION=$1
fi

docker kill unbalance && docker rm unbalance

docker run -it \
-v /etc/localtime:/etc/localtime:ro \
-v /mnt/user/data:/config \
-v /mnt/user/data:/log \
-v /mnt:/mnt \
-v /root:/root \
-p 6237:6237 \
jbrodriguez/unbalance:$VERSION

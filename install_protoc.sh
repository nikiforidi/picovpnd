#!/bin/sh

apt update -y
apt install -y unzip
PROTOC_VERSION=$(curl -s "https://api.github.com/repos/protocolbuffers/protobuf/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')
wget -qO protoc.zip https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-$PROTOC_VERSION-linux-x86_64.zip
unzip -q protoc.zip bin/protoc -d /usr/local/
chmod a+x /usr/local/bin/protoc

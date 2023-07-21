#!/bin/bash

NAME=sankaku-viewer
OS=$1
ARCH=$2

if [[ "$ARCH" == "386" ]]; then
    DIR="$NAME-$OS-x86"
else
    DIR="$NAME-$OS-$ARCH"
fi

if [[ "$OS" == "windows" ]]; then
    NAME="$NAME.exe"
fi

# Build the project
echo "Building the project..."
mkdir -p $DIR
env GOOS=$OS GOARCH=$ARCH go build -v -o $DIR/$NAME -trimpath -ldflags "-s -w" cmd/sankaku/main.go
tar cf $DIR.tar $DIR --format=posix --owner=0 --group=0
gzip -9 $DIR.tar
rm -rf $DIR

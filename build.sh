#!/bin/bash

env GOAMD64=v3 go build -o sankaku.exe -v -trimpath -ldflags "-s -w" ./cmd/main.go
env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o sankaku -v -trimpath -ldflags "-s -w" ./cmd/main.go

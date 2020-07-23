#!/bin/bash

# Get packages
go get

# Build acker for linux, macOS and windows
GOOS=linux GOARCH=amd64 go build -a -trimpath -ldflags="-s -w" -o acker_linux_amd64
GOOS=darwin GOARCH=amd64 go build -a -trimpath -ldflags="-s -w" -o acker_macos_amd64
GOOS=windows GOARCH=amd64 go build -a -trimpath -ldflags="-s -w" -o acker_windows_amd64.exe

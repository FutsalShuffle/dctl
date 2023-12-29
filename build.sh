#!/bin/bash
env GOOS=darwin GOARCH=amd64 go build -o dctl_amd64_darwin .
env GOOS=darwin GOARCH=arm64 go build -o dctl_arm64_darwin .
env GOOS=linux GOARCH=amd64 go build -o dctl_amd64_linux .
env GOOS=windows GOARCH=amd64 go build -o dctl_amd64_win .

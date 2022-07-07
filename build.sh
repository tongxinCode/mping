#!/bin/bash
rm -rf build/
if [ ! -d build ]
then
    mkdir build
fi
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/mping.exe main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/mping main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o build/mping-arm main.go
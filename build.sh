#!/bin/bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.version=$1'" -o=bin/moxy-windows-amd64.exe
if [[ ! -e bin/moxy-windows-amd64.exe ]]; then
    exit 1
fi

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.version=$1'" -o=bin/moxy-darwin-amd64
if [[ ! -e bin/moxy-darwin-amd64 ]]; then
    exit 1
fi

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.version=$1'" -o=bin/moxy-linux-amd64
if [[ ! -e bin/moxy-linux-amd64 ]]; then
    exit 1
fi

ls -la bin
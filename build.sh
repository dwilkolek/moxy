#!/bin/bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X 'app.version=$1'" -o=bin/moxy-windows-amd64.exe ./app/cmd 
if [[ ! -e bin/moxy-windows-amd64.exe ]]; then
    exit 1
fi
cp bin/moxy-windows-amd64.exe bin/moxy-windows.exe

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -ldflags="-X 'app.version=$1'" -o=bin/moxy-darwin-amd64 ./app/cmd
if [[ ! -e bin/moxy-darwin-amd64 ]]; then
    exit 1
fi
cp bin/moxy-darwin-amd64 bin/moxy-mac

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'app.version=$1'" -o=bin/moxy-linux-amd64 ./app/cmd 
if [[ ! -e bin/moxy-linux-amd64 ]]; then
    exit 1
fi
cp bin/moxy-linux-amd64 bin/moxy-linux

ls -la bin
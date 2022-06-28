CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.version=$1'" -o=bin/moxy-windows.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.version=$1'" -o=bin/moxy-mac
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.version=$1'" -o=bin/moxy-linux

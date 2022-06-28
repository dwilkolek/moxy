export GOOS="linux"
go build -ldflags "-X main.version=$1" -o bin/moxy-linux github.com/dwilkolek/moxy
export GOOS="darwin"
go build -ldflags "-X main.version=$1" -o bin/moxy-mac github.com/dwilkolek/moxy 
export GOOS="windows"
go build -ldflags "-X main.version=$1" -o bin/moxy-windows.exe github.com/dwilkolek/moxy 

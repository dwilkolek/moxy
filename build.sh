export GOOS="linux"
go build -o bin/moxy-linux github.com/dwilkolek/moxy
export GOOS="darwin"
go build -o bin/moxy-mac github.com/dwilkolek/moxy 
export GOOS="windows"
go build -o bin/moxy-windows.exe github.com/dwilkolek/moxy 

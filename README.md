# MOXY
This is application to setup ssh tunnel and expose http servers under bastion.
It setup ssh tunnel with host with key authorizastion to destination specified in config and exposes http proxy to that server by adding Host header.

## Builiding
```
$Env:GOOS="linux"
go build -o bin/moxy-linux github.com/dwilkolek/moxy
$Env:GOOS="darwin"
go build -o bin/moxy-mac github.com/dwilkolek/moxy 
$Env:GOOS="windows"
go build -o bin/moxy-windows.exe github.com/dwilkolek/moxy 
```

## Config
Example `config.json`.
```
{
    "tunnel": {
        "userAndHost": "sshuser@localhost",
        "pathToPrivateKey": "/absolute/path/to/private/key/id_rsa",
        "destination": "localhost:8080"
    },
    "services": {
        "some": {
            "port": 9000,
            "headers": {
                "some": "value",
                "another":"header"
            } 
        },
        "other": {
            "port": 9001
        }
    }
}
```
This config will connect to `sshuser@localhost` using key `server/id_rsa` (use absolute path) and tunnel to `localhost:8080`.
It will expose 
- ssh unnel under random port (available std output)
- `localhost:9000` http server under available at `localhost:8080` called with header `Host: some.service`
- `localhost:9901` http server under available at `localhost:8080` called with header `Host: other.service`


## Running
### Windows
```
./moxy-windows.exe ./config.json
```
### Linux
```
chmod +x ./moxy-linux
./moxy-linux ./config.json
```
### Macos
```
chmod +x ./moxy-mac
./moxy-mac ./config.json
```

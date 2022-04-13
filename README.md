# MOXY
This is application to set up ssh tunnel and expose http servers under bastion.
It set up ssh tunnel with host with key authorization to destination specified in config and exposes http proxy with headers beatification required by external service.

## Building
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
                "host": "some.service"
                "some": "value"
            },
            "allowCors": true
        },
        "other": {
            "port": 9001
        }
    }
}
```
This config will connect to `sshuser@localhost` using key `/absolute/path/to/private/key/id_rsa` (use absolute path) and tunnel to `localhost:8080`.
It will expose 
- ssh tunnel under random port (available in std output)
- `localhost:9000` http server under available at tunneled `localhost:8080` called with headers `Host: some.service`, `"some: value"` and will handle preflight OPTION request to fulfill browsers cors check.
- `localhost:9901` http server under available at tunneled `localhost:8080`.


## Running
If you don't provide config as an arg then `config.json` will be used.

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

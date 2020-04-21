# Storage-Center

## Start Example

``` bash
go run server_main.go
go run client_main.go
```

You can use go run *_main.go -h to see the usage

## HTTP Usage

``` bash
curl localhost:9091/get?key=animal
curl localhost:9090/set -X POST -H "Content-Type:application/json" -d '{"name":"animal","value":"pig"}'
```

## Build

the output is placed in `build/`

``` bash
make client
make server
```

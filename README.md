# gostatic
super light http server

## How to Build
```
GOARCH=amd64 GOOS=linux go build  -ldflags "-linkmode external -extldflags -static -w"
```

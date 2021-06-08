* for MAC

```
go build -o parser main.go
```

* for WINDOWS

```
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o parser.exe main.go
```
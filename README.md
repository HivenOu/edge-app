编译树莓派运行的命令：CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$ldflags" -o app app.go
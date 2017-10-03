go get golang.org/x/sys/unix
set GOARCH=amd64
set GOOS=linux
go build 

set GOOS=windows
go build 
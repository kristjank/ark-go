go get golang.org/x/sys/unix
set GOARCH=amd64
set GOOS=linux
go build 

set GOOS=windows
go build 

REM linux
mkdir log
cd ..

REM windows
mkdir log
cd ..

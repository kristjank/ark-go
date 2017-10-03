
go get golang.org/x/sys/unix
set GOARCH=amd64
set GOOS=linux
go build 

set GOOS=windows
go build 

REM linux
if not exist "dist" mkdir dist
cd dist
mkdir linux
cd linux
move ..\..\arkgopool .
mkdir settings
copy ..\..\settings\banner.txt settings
copy ..\..\settings\sample.config.toml settings
cd ..

REM windows
if not exist "windows" mkdir windows
cd windows
move ..\..\arkgopool.exe .
mkdir settings
copy ..\..\settings\banner.txt settings
copy ..\..\settings\sample.config.toml settings
mkdir log


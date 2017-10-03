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
copy ..\..\arkgopool .
mkdir cfg
copy ..\..\cfg\banner.txt cfg
copy ..\..\cfg\sample.config.toml cfg
mkdir log
cd ..

REM windows
if not exist "windows" mkdir windows
cd windows
copy ..\..\arkgopool.exe .
mkdir cfg
copy ..\..\cfg\banner.txt cfg
copy ..\..\cfg\sample.config.toml cfg
mkdir log
cd ..


REM create archive
for /d %%X in (*) do "c:\Program Files\7-Zip\7z.exe" a -mx "%%X.zip" "%%X\*"

move linux.zip ARKGOStats-LinuxRelease.zip
move windows.zip ARKGOStats-WindowsRelease.zip


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
copy ..\..\arkgoserver .
mkdir cfg
copy ..\..\cfg\banner.txt cfg
copy ..\..\cfg\sample.config.toml cfg
mkdir log
cd ..

REM windows
if not exist "windows" mkdir windows
cd windows
copy ..\..\arkgoserver.exe .
mkdir cfg
copy ..\..\cfg\banner.txt cfg
copy ..\..\cfg\sample.config.toml cfg
mkdir log
cd ..

REM create archive
for /d %%X in (*) do "c:\Program Files\7-Zip\7z.exe" a -mx "%%X.zip" "%%X\*"

move linux.zip ARKGOServer-LinuxRelease.zip
move windows.zip ARKGOServer-WindowsRelease.zip


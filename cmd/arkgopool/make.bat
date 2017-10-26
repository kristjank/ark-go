REM go get golang.org/x/sys/unix
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
mkdir settings
copy ..\..\settings\banner.txt settings
copy ..\..\settings\sample.config.toml settings
mkdir log
cd ..

REM windows
if not exist "windows" mkdir windows
cd windows
copy ..\..\arkgopool.exe .
mkdir settings
copy ..\..\settings\banner.txt settings
copy ..\..\settings\sample.config.toml settings
mkdir log
cd ..

REM UsbKeyRelease
if not exist "usb" mkdir usb
cd usb
move ..\..\arkgopool.exe .
move ..\..\arkgopool .
mkdir settings
copy ..\..\settings\banner.txt settings
copy ..\..\settings\sample.config.toml settings
mkdir log
cd ..

REM create archive
for /d %%X in (*) do "c:\Program Files\7-Zip\7z.exe" a -mx "%%X.zip" "%%X\*"

move linux.zip ARKGO-LinuxRelease.zip
move windows.zip ARKGO-WindowsRelease.zip
move usb.zip ARKGO-UsbKeyRelease.zip


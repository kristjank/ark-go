set GOARCH=arm
set GOOS=linux
go build 

if not exist "dist" mkdir dist
cd dist

if not exist "linuxarm" mkdir linuxarm
cd linuxarm
move ..\..\arkgoserver .
mkdir cfg
copy ..\..\cfg\banner.txt cfg
copy ..\..\cfg\sample.config.toml cfg
mkdir log
cd ..

REM create archive
for /d %%X in (*) do "c:\Program Files\7-Zip\7z.exe" a -mx "%%X.zip" "%%X\*"

move linuxarm.zip ARKGOServer-LinuxRelease_ARM.zip


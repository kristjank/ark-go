go get golang.org/x/sys/unix
set GOARCH=amd64
set GOOS=linux
go build .\arkgopool.go .\database.go .\helpers.go .\payouts.go

set GOOS=windows
go build .\arkgopool.go .\database.go .\helpers.go .\payouts.go
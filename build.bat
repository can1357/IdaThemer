set GOOS=windows
set GOARCH=amd64
go build -o build/idathemer-windows-amd64.exe
set GOOS=linux
set GOARCH=amd64
go build -o build/idathemer-linux-amd64
set GOOS=darwin
set GOARCH=amd64
go build -o build/idathemer-darwin-amd64
set GOOS=linux
set GOARCH=arm64
go build -o build/idathemer-linux-arm64
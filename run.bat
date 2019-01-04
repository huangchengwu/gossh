%GOPATH%/bin/go-bindata.exe -o=./asset/asset.go -pkg=asset   www/ www/js



REM linux make 
REM set GOOS=linux 
REM go build run.go





REM windows make 
set GOOS=windows
go build run.go
run.exe

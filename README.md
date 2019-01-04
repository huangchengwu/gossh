使用GO语言开发的一个webssh

编辑静态文件
%GOPATH%/bin/go-bindata.exe -o=./asset/asset.go -pkg=asset   www/ www/js


编译linux
REM linux make 
REM set GOOS=linux 
REM go build run.go




编译windows
REM windows make 
set GOOS=windows
go build run.go
run.exe

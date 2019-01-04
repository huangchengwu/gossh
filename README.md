使用GO语言开发的一个webssh

编辑静态文件
%GOPATH%/bin/go-bindata.exe -o=./asset/asset.go -pkg=asset   www/ www/js


编译linux
REM linux make 
REM set GOOS=linux 
REM go build run.go




编译windows\n
\nREM windows make 
\nset GOOS=windows
\ngo build run.go
\nrun.exe

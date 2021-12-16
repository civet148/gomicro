@echo off

echo generating...

protoc -I=.  --proto_path=%GOPATH%\src  --micro_out=.  --go_out=.  echo.proto

echo generate over

pause

@echo off

REM LINUX COMPILATION


set GOOS=linux
set GOARCH=386

echo.
echo %GOOS%
echo %GOARCH%
echo.

go build -o smtpmaill smtpmail.go global.go authlogin.go encoding.go help.go config.go


set GOOS=linux
set GOARCH=amd64

echo.
echo %GOOS%
echo %GOARCH%
echo.

go build -o smtpmaill64 smtpmail.go global.go authlogin.go encoding.go help.go config.go

@echo off

REM WINDOWS COMPILATION


set GOOS=windows
set GOARCH=386

echo.
echo %GOOS%
echo %GOARCH%
echo.

go build -o smtpmail.exe smtpmail.go global.go authlogin.go encoding.go help.go config.go

set GOOS=windows
set GOARCH=amd64

echo.
echo %GOOS%
echo %GOARCH%
echo.

go build -o smtpmail64.exe smtpmail.go global.go authlogin.go encoding.go help.go config.go



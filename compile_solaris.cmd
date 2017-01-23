@echo off

REM SOLARIS COMPILATION


set GOOS=solaris
set GOARCH=amd64

echo.
echo %GOOS%
echo %GOARCH%
echo.

go build -o smtpmails64 smtpmail.go global.go authlogin.go encoding.go help.go config.go

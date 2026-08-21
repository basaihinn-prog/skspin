@echo off
setlocal

set "serviceDir=.\servicejili"
set "outputDir=.\bin"
set "GOOS=linux"
set "GOARCH=arm64"

if not exist "%outputDir%" (
    mkdir "%outputDir%"
)

sd "tada:" "jili:" %serviceDir%\jiliut\regRpc.go >NUL 2>&1
for /d %%G in (servicejili\*) do (
    echo %%G
    set "service=%%~nxG"
    CALL :build_service %service%
)
if /i "%~1"=="skip" exit /b 0
cd %outputDir%&&7z a %outputDir% *
goto :eof

:build_service
sd "tada_" "jili_" %serviceDir%\%service%\internal\const.go >NUL 2>&1
set "out=%service:tada_=jili_%"
go build -trimpath -ldflags "-s -w" -o %outputDir%\%out% %serviceDir%\%service%
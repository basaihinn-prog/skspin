@echo off
setlocal

set "serviceDir=.\service"
set "outputDir=bin"
set "GOOS=linux"
set "GOARCH=arm64"

 rmdir /s /q "%outputDir%"
if not exist "%outputDir%" (
    mkdir "%outputDir%"
)

for /d %%G in (service\*) do (
    echo %%G
    set "service=%%~nxG"
    CALL :build_service %service%
)

cd %outputDir%&&7z a %outputDir% *
goto :eof

:build_service
go build -trimpath -ldflags "-s -w" -o %outputDir%\%service% %serviceDir%\%service%
@echo off
setlocal

set "serviceDir=.\servicepp"
set "outputDir=.\bin"
set "GOOS=linux"
set "GOARCH=arm64"

if not exist "%outputDir%" (
    mkdir "%outputDir%"
)


for /d %%G in (servicepp\*) do (
    echo %%G
    set "service=%%~nxG"
    CALL :build_service %service%
)
if /i "%~1"=="skip" exit /b 0
cd %outputDir%&&7z a %outputDir% *
goto :eof

:build_service
go build -trimpath -ldflags "-s -w" -o %outputDir%\%out% %serviceDir%\%service%
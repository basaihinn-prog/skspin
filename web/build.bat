rd /s /q dist
rem call npm run build 
call npm run build:superHigh
set dst=dist
cd %dst%&&7z a %dst% *
pause 

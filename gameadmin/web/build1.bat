rd /s /q dist
call npm run build 
set dst=dist
cd %dst%&&7z a %dst% *
pause 

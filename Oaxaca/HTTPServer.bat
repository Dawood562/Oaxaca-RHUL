echo off
color 2
echo Server will be hosted on port 8000
echo To access website visit http://localhost:8000
python -m http.server
if %ERRORLEVEL% neq 0 goto ProcessError
goto End

:ProcessError
py -m http.server

:End
pause
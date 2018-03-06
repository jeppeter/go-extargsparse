
echo off

set BATDIR=%~dp0
set CURDIR=%CD%

cd %BATDIR%\releasechk && go build && .\releasechk.exe && del .\releasechk.exe && cd %CURDIR% || cd %CURDIR%
echo on
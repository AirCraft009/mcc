@echo off
setlocal

go build -o bin/mcc.exe -v ./Assembler-main

set BASEDIR=%~dp0..\
set MCC=%BASEDIR%bin\mcc.exe

echo %BASEDIR%
echo %MCC%

:: ---- stdlib ----
pushd "%BASEDIR%lib\stdlib\sources"

for %%F in (*.asm) do (
    echo Compiling stdlib %%F
    "%MCC%" -i "%%F" -o "%BASEDIR%lib\stdlib\obj\%%~nF.obj" -n=true -s=true
)

popd

:: ---- include ----
pushd "%BASEDIR%\lib\include\sources"

for %%F in (*.asm) do (
    echo Compiling include %%F
    "%MCC%" -i "%%F" -o "%BASEDIR%\lib\include\obj\%%~nF.obj" -n=true -s=true
)

popd

endlocal

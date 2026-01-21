@echo off
setlocal

set BASEDIR=%~dp0
set MCC=%BASEDIR%bin\mcc.exe

:: ---- stdlib ----
pushd "%BASEDIR%stdlib\sources"

for %%F in (*.asm) do (
    echo Compiling stdlib %%F
    "%MCC%" "%%F" --o "%BASEDIR%stdlib\obj\%%~nF.obj" --no_link=true
)

popd

:: ---- include ----
pushd "%BASEDIR%include\sources"

for %%F in (*.asm) do (
    echo Compiling include %%F
    "%MCC%" "%%F" --o "%BASEDIR%include\obj\%%~nF.obj" --no_link=true
)

popd

endlocal

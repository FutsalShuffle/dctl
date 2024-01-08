set DownloadUrl="https://github.com/FutsalShuffle/dctl/releases/download/v1.0/dctl_amd64_windows"
set saveTo="%userprofile%\.dctl\dctl.exe"
if not exist %userprofile%\.dctl mkdir %userprofile%\.dctl

powershell -Command "Invoke-WebRequest %DownloadUrl% -OutFile %saveTo%"
echo %userprofile%\.dctl\dctl.exe

@echo off
for /f "tokens=2*" %%a in ('reg query "HKLM\System\CurrentControlSet\Control\Session Manager\Environment" /v Path 2^>^&1^|find "REG_"') do @set oldPath=%%b
set newPath=%oldPath%;%userprofile%\.dctl
echo %newPath%
REG ADD "HKLM\System\CurrentControlSet\Control\Session Manager\Environment" /f /v Path /t REG_EXPAND_SZ /d "%newPath%"
setx dctli "dctli"

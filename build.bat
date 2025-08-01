@echo off
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
echo 开始构建...

REM 删除旧的构建文件
if exist build rd /s /q build

REM 构建程序
echo 正在编译...
go build -ldflags="-s -w" -o dmxsmart-auto-outbound.exe main.go
if %errorlevel% neq 0 (
    echo 构建失败！
    exit /b %errorlevel%
)

REM UPX压缩
echo 正在压缩可执行文件...
where upx >nul 2>nul
if %errorlevel% equ 0 (
    upx dmxsmart-auto-outbound.exe
) else (
    echo UPX未安装，使用Windows内置压缩...
    compact /c /exe:lzx dmxsmart-auto-outbound.exe >nul
    if %errorlevel% neq 0 (
        echo Windows压缩失败！
    )
)

powershell -Command "Compress-Archive -Path dmxsmart-auto-outbound.exe,config.yaml -DestinationPath dmxsmart-auto-outbound.zip -Force"
if %errorlevel% neq 0 (
    echo 创建ZIP文件失败！
    exit /b %errorlevel%
)

REM 清理临时文件
del dmxsmart-auto-outbound.exe

@REM 移动到 build 目录
mkdir build
move dmxsmart-auto-outbound.zip build

echo 构建完成！



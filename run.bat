@echo off
echo GoClaw AI Assistant
echo ===================
echo.

REM Check if config exists
if not exist "%USERPROFILE%\.goclaw\config.yaml" (
    echo Config file not found at %USERPROFILE%\.goclaw\config.yaml
    echo Please create it using config.example.yaml
    echo.
    pause
    exit /b 1
)

REM Start the server
echo Starting server on http://localhost:8080
echo Press Ctrl+C to stop
echo.
goclaw.exe serve

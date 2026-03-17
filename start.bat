@echo off
echo Starting GoClaw AI Assistant...

REM Check if config file exists
if not exist "%USERPROFILE%\.goclaw\config.yaml" (
    echo Config file not found. Creating default config...
    mkdir "%USERPROFILE%\.goclaw"
    copy config.example.yaml "%USERPROFILE%\.goclaw\config.yaml"
    echo Default config created at %USERPROFILE%\.goclaw\config.yaml
    echo Please edit it with your API keys and credentials.
)

REM Start the server
goclaw.exe serve

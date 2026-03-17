#!/bin/bash

# GoClaw startup script

echo "Starting GoClaw AI Assistant..."

# Check if config file exists
if [ ! -f ~/.goclaw/config.yaml ]; then
    echo "Config file not found. Creating default config..."
    mkdir -p ~/.goclaw
    cp config.example.yaml ~/.goclaw/config.yaml
    echo "Default config created at ~/.goclaw/config.yaml"
    echo "Please edit it with your API keys and credentials."
fi

# Start the server
./goclaw.exe serve

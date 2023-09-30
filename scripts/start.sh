#!/bin/bash
cd "$(dirname "$0")"
SCRIPT_DIR="$(pwd)"
echo "Starting container..."
echo "Clearing log file..."
> $SCRIPT_DIR/../log.txt
echo "Running image..."
docker run -d -p $(cat $SCRIPT_DIR/../.env | grep APP_PORT= | cut -d: -f2 | awk '/^/ { print $1":"$1 }') -v $SCRIPT_DIR/../log.txt:/usr/src/app/log.txt --name login-api ccu-login-api
echo "Container started"
#!/bin/bash
echo "Starting container..."
echo "Clearing log file..."
> log.txt
echo "Running image..."
docker run -d -p $(cat ../.env | grep APP_PORT= | cut -d: -f2 | awk '/^/ { print $1":"$1 }') -v ../log.txt:/usr/src/app/log.txt --name login-api ccu-login-api
echo "Container started"
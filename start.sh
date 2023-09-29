#!/bin/bash
> log.txt
docker run -d -p $(cat .env | grep APP_PORT= | cut -d: -f2 | awk '/^/ { print $1":"$1 }') -v $(pwd)/log.txt:/usr/src/app/log.txt --name login-api ccu-login-api
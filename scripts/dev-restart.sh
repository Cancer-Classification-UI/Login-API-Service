#!/bin/bash
# ONLY USE THIS WHEN IN DEVELOPMENT, DO NOT USE THIS IN PRODUCTION IT WILL
# BREAK THE IMAGE CAPTURE FOR THE REGISTRY
echo "Rebuilding image..."
docker build -t ccu-login-api .
./restart.sh
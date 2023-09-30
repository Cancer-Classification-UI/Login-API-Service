#!/bin/bash
echo "Killing and removing container..."
docker kill login-api
docker rm login-api
echo "Container stopped"

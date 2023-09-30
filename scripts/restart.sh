#!/bin/bash
cd "$(dirname "$0")"
echo "Restarting container..."
./stop.sh
./start.sh
#!/bin/bash
cd "$(dirname "$0")"
echo "Restarting..."
./stop.sh
./start.sh
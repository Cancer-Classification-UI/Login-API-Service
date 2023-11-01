#!/bin/bash
cd "$(dirname "$0")"
SCRIPT_DIR="$(pwd)"
echo "Starting containers..."

echo "Clearing log file..."
> $SCRIPT_DIR/../log.txt

echo "Creating compose from template..."

# Get URI from .env
MONGO_URI=$(grep -oP 'MONGODB_URI=\K.*' $SCRIPT_DIR/../.env)

# Extract the username, password, port, and database name from MONGO_URI
MONGO_USERNAME=$(echo $MONGO_URI | awk -F '://' '{print $2}' | awk -F ':' '{print $1}')
MONGO_PASSWORD=$(echo $MONGO_URI | awk -F '://' '{print $2}' | awk -F ':' '{print $2}' | awk -F '@' '{print $1}')

# Check if .env file exists, if it doesnt go for default values
if [ -f "$SCRIPT_DIR/../.env" ]; then
    sed -e "s/<MONGO_ROOT_USERNAME>/$MONGO_USERNAME/g" \
        -e "s/<MONGO_ROOT_PASSWORD>/$MONGO_PASSWORD/g" \
        -e "s/<API_PORT>/$(grep -oP 'APP_PORT=\K.*' $SCRIPT_DIR/../.env)/g" \
        -e "s/<MONGO_PORT>/$(grep -oP 'MONGODB_REDIRECT=\K.*' $SCRIPT_DIR/../.env)/g" \
        $SCRIPT_DIR/../docker-compose-template.yaml > $SCRIPT_DIR/../docker-compose.yaml
else 
    sed -e "s/<MONGO_ROOT_USERNAME>/ccu/g" \
        -e "s/<MONGO_ROOT_PASSWORD>/password/g" \
        -e "s/<API_PORT>/8084/g" \
        -e "s/<MONGO_PORT>/27084/g" \
        $SCRIPT_DIR/../docker-compose-template.yaml > $SCRIPT_DIR/../docker-compose.yaml
fi

cd $SCRIPT_DIR/../
docker-compose up -d

echo "Containers started"
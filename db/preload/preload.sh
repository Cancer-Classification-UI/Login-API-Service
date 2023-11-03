#!/bin/bash

MONGO_DB="$MONGO_INITDB_DATABASE"
MONGO_USER="$MONGO_INITDB_ROOT_USERNAME"
MONGO_PASSWORD="$MONGO_INITDB_ROOT_PASSWORD"
MONGO_AUTH_DB="admin"

# Directory containing JSON files
JSON_DIR="/preload"

# Wait for mongodb connetion
while true; do
  # Ping MongoDB
  if mongosh --eval 'db.runCommand({ ping: 1})' > /dev/null; then
    echo "Successful ping, MongoDB is ready!"
    break
  else
    echo "Waiting for MongoDB to start..."
    sleep 1
  fi
done

# Loop through JSON files in the directory
echo "Importing data..."
for json_file in "$JSON_DIR"/*.json; do
    if [ -f "$json_file" ]; then
        # Extract the filename (without the extension) as the collection name
        collection_name=$(basename "$json_file" .json)

        # Run mongoimport for the current JSON file
        echo mongoimport --host 127.0.0.1:27017 --db "$MONGO_DB" --collection "$collection_name" --file "$json_file" --username "$MONGO_USER" --password "$MONGO_PASSWORD" --authenticationDatabase "$MONGO_AUTH_DB"
        mongoimport --host 127.0.0.1:27017 --db "$MONGO_DB" --collection "$collection_name" --file "$json_file" --username "$MONGO_USER" --password "$MONGO_PASSWORD" --authenticationDatabase "$MONGO_AUTH_DB"

        echo "Imported $json_file into collection $collection_name"
    fi
done
# Add any other commands or logs you need
echo "Data import completed."


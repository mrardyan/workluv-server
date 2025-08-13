#!/bin/sh

# Wait for database and Redis to be ready
echo "Waiting for services to be ready..."
sleep 10

# Start the application
echo "Starting Go server..."
exec ./main

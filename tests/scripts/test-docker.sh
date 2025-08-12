#!/bin/bash

echo "🧪 Testing Docker Compose Setup"
echo "================================"

# Build the images
echo "📦 Building Docker images..."
docker-compose build

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed"
    exit 1
fi

echo "✅ Docker build successful"

# Start services
echo "🚀 Starting services..."
docker-compose up -d postgres redis

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Check service health
echo "🔍 Checking service health..."
docker-compose ps

# Test database connection
echo "🗄️  Testing database connection..."
docker-compose exec postgres pg_isready -U postgres

# Test Redis connection
echo "🔴 Testing Redis connection..."
docker-compose exec redis redis-cli ping

# Stop services
echo "🛑 Stopping services..."
docker-compose down

echo "✅ Docker setup test completed successfully!"

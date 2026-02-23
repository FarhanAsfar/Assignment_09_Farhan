#!/bin/bash

echo "=================================================================="
echo "    Airbnb Market Scraping System - Quick Start"
echo "=================================================================="
echo ""

# Step 1: Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ Failed to install Go dependencies"
    exit 1
fi
echo "✅ Go dependencies installed"
echo ""

# Step 2: Start PostgreSQL
echo "🐘 Starting PostgreSQL database..."
docker-compose up -d
if [ $? -ne 0 ]; then
    echo "❌ Failed to start PostgreSQL"
    exit 1
fi
echo "✅ PostgreSQL started"
echo ""

# Step 3: Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
sleep 3

# Step 4: Check database connection
echo "🔍 Verifying database connection..."
docker exec airbnb-postgres psql -U postgres -d airbnb_scraper -c "SELECT 1;" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "⚠️  Database not ready yet, waiting 5 more seconds..."
    sleep 5
fi
echo "✅ Database is ready"
echo ""

# Step 5: Start scraping
echo "=================================================================="
echo "🚀 Starting Airbnb scraper..."
echo "=================================================================="
echo ""

go run main.go

# Done
echo ""
echo "=================================================================="
echo "✅ Scraping complete!"
echo "=================================================================="
echo ""
echo "To view analytics anytime, run:"
echo "  go run main.go --show-stats"
echo ""
#!/bin/bash

# VPS Deployment Script for Contact Management API
# This script applies all necessary fixes for VPS deployment

echo "🚀 Starting VPS deployment fixes..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "Please run this script from the project root directory (where go.mod is located)"
    exit 1
fi

print_status "1. Pulling latest changes from Git..."
git pull origin main

if [ $? -ne 0 ]; then
    print_error "Failed to pull from Git. Please check your Git configuration."
    exit 1
fi

print_status "2. Stopping any running application..."
pkill -f "./bin/server" || true
pkill -f "go run ./cmd/server" || true

print_status "3. Resetting database..."
mysql -u yudo -p'P@ssw0rd' -e "USE getcontact; DROP TABLE IF EXISTS contacts; DROP TABLE IF EXISTS users; DROP TABLE IF EXISTS schema_migrations;" 2>/dev/null || print_warning "Database reset may have failed - continuing..."

print_status "4. Running database migrations..."
make migrate-up

if [ $? -ne 0 ]; then
    print_error "Migration failed! Check the error above."
    exit 1
fi

print_status "5. Building application..."
go build -o bin/server ./cmd/server

if [ $? -ne 0 ]; then
    print_error "Build failed! Check the error above."
    exit 1
fi

print_status "6. Starting application..."
./bin/server &

# Wait a moment for the app to start
sleep 3

print_status "7. Testing application health..."
curl -s -X GET "http://localhost:8080/api/v1/health" > /dev/null

if [ $? -eq 0 ]; then
    print_status "✅ Application is running successfully!"
    print_status "🌐 API available at: http://13.229.87.19:8080"
    print_status "📊 Health check: http://13.229.87.19:8080/api/v1/health"
else
    print_error "❌ Application health check failed!"
    print_status "Check logs with: tail -f logs/app.log"
fi

print_status "8. Checking running processes..."
ps aux | grep -E "(server|go)" | grep -v grep

print_status "🎉 VPS deployment completed!"
echo ""
echo "Useful commands:"
echo "- Check logs: tail -f logs/app.log"
echo "- Restart app: pkill -f './bin/server' && ./bin/server &"
echo "- Check database: mysql -u yudo -p'P@ssw0rd' -e 'USE getcontact; SHOW TABLES;'"

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
print_status "Checking for existing processes..."
ps aux | grep -E "(server|go)" | grep -v grep || print_status "No existing processes found"

# Kill any existing server processes more aggressively
pkill -9 -f "./bin/server" || true
pkill -9 -f "go run ./cmd/server" || true

# Wait a moment to ensure processes are killed
sleep 2

# Double check no processes are running on port 8080
if lsof -ti:8080 > /dev/null 2>&1; then
    print_warning "Port 8080 still in use, forcing kill..."
    lsof -ti:8080 | xargs kill -9 || true
    sleep 1
fi

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
print_status "Launching ./bin/server..."
./bin/server > logs/app.log 2>&1 &

# Get the PID of the background process
SERVER_PID=$!
print_status "Server started with PID: $SERVER_PID"

# Wait a moment for the app to start
sleep 5

# Check if the process is still running
if kill -0 $SERVER_PID 2>/dev/null; then
    print_status "Server process is running"
else
    print_error "Server process failed to start!"
    print_status "Checking logs..."
    tail -20 logs/app.log
    exit 1
fi

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

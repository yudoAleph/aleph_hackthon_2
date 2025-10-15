#!/bin/bash

# VPS Troubleshooting Script for Contact Management API
# This script diagnoses and fixes common VPS deployment issues

echo "🔍 Starting VPS troubleshooting..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

print_info() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "Please run this script from the project root directory (where go.mod is located)"
    exit 1
fi

print_status "1. Checking system information..."
echo "Current directory: $(pwd)"
echo "User: $(whoami)"
echo "Date: $(date)"

print_status "2. Checking if port 8080 is in use..."
if lsof -i :8080 > /dev/null 2>&1; then
    print_info "Port 8080 is in use:"
    lsof -i :8080
else
    print_warning "Port 8080 is not in use!"
fi

print_status "3. Checking running processes..."
ps aux | grep -E "(server|go)" | grep -v grep || print_warning "No Go server processes found"

print_status "4. Checking if binary exists..."
if [ -f "bin/server" ]; then
    print_info "Binary exists: $(ls -la bin/server)"
else
    print_error "Binary does not exist! Need to build first."
    print_status "Building application..."
    go build -o bin/server ./cmd/server
    if [ $? -ne 0 ]; then
        print_error "Build failed!"
        exit 1
    fi
fi

print_status "5. Checking environment configuration..."
if [ -f ".env" ]; then
    print_info ".env file found in root"
    grep -E "^PORT|^DB_HOST|^DB_USER" .env || print_warning "Some config variables missing"
elif [ -f "configs/.env" ]; then
    print_info ".env file found in configs/"
    grep -E "^PORT|^DB_HOST|^DB_USER" configs/.env || print_warning "Some config variables missing"
else
    print_error "No .env file found!"
fi

print_status "6. Testing database connection..."
mysql -u yudo -p'P@ssw0rd' -e "USE getcontact; SHOW TABLES;" 2>/dev/null
if [ $? -eq 0 ]; then
    print_info "Database connection successful"
else
    print_error "Database connection failed!"
fi

print_status "7. Checking recent logs..."
if [ -f "logs/app.log" ]; then
    print_info "Last 10 lines of app.log:"
    tail -10 logs/app.log
else
    print_warning "No app.log file found"
fi

print_status "8. Attempting to start application..."
print_info "Stopping any existing processes..."
pkill -9 -f "./bin/server" || true
sleep 2

print_info "Starting application..."
./bin/server > logs/app.log 2>&1 &
SERVER_PID=$!
print_info "Server started with PID: $SERVER_PID"

# Wait for startup
sleep 3

# Check if process is still running
if kill -0 $SERVER_PID 2>/dev/null; then
    print_status "✅ Server process is running!"
else
    print_error "❌ Server process failed to start!"
    print_info "Checking logs for errors:"
    tail -20 logs/app.log
    exit 1
fi

print_status "9. Testing health endpoint..."
if curl -s -f http://localhost:8080/api/v1/health > /dev/null; then
    print_status "✅ Health check passed!"
    curl -s http://localhost:8080/api/v1/health
else
    print_error "❌ Health check failed!"
fi

print_status "10. Testing external access..."
EXTERNAL_IP="13.229.87.19"
if curl -s -f --max-time 10 http://$EXTERNAL_IP:8080/api/v1/health > /dev/null; then
    print_status "✅ External access works!"
else
    print_warning "❌ External access failed!"
    print_info "This might be due to firewall settings or binding issues"
    print_info "Check if application is binding to 0.0.0.0 instead of localhost"
fi

print_status "11. Network diagnostics..."
print_info "Checking listening ports:"
netstat -tlnp | grep :8080 || print_warning "Port 8080 not listening"

print_info "Checking firewall status:"
if command -v ufw >/dev/null 2>&1; then
    ufw status | grep 8080 || print_warning "Port 8080 might be blocked by UFW"
else
    print_info "UFW not found, checking iptables..."
    iptables -L | grep 8080 || print_info "No specific iptables rules for port 8080"
fi

print_status "🎯 Troubleshooting completed!"
echo ""
echo "📋 Summary of findings:"
echo "- Server PID: $SERVER_PID"
echo "- Port 8080 status: $(lsof -i :8080 > /dev/null 2>&1 && echo 'IN USE' || echo 'FREE')"
echo "- Health check: $(curl -s -f http://localhost:8080/api/v1/health > /dev/null 2>&1 && echo 'PASS' || echo 'FAIL')"
echo "- External access: $(curl -s -f --max-time 5 http://$EXTERNAL_IP:8080/api/v1/health > /dev/null 2>&1 && echo 'WORKING' || echo 'BLOCKED')"
echo ""
echo "🔧 Useful commands:"
echo "- Check logs: tail -f logs/app.log"
echo "- Restart app: pkill -f './bin/server' && ./bin/server &"
echo "- Test API: curl http://localhost:8080/api/v1/health"
echo "- Check port: netstat -tlnp | grep 8080"

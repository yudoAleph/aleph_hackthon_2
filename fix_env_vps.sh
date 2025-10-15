#!/bin/bash

# Fix .env configuration on VPS
# This script updates the database credentials and ensures proper configuration

echo "🔧 Fixing .env configuration on VPS..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Please run this script from the project root directory"
    exit 1
fi

# Create or update the .env file in configs/ with correct database credentials
cat > configs/.env << 'EOF'
# Server Configuration
PORT=8080
ENVIRONMENT=production
ALLOWED_ORIGINS=*

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=yudo
DB_PASSWORD=P@ssw0rd
DB_NAME=getcontact
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=HackthonII-2025
EOF

echo "✅ .env file updated with correct database credentials"
echo "📋 Configuration:"
echo "- DB_USER: yudo"
echo "- DB_PASSWORD: P@ssw0rd"
echo "- DB_NAME: getcontact"
echo "- PORT: 8080"

# Verify the file was created correctly
if [ -f "configs/.env" ]; then
    echo "✅ configs/.env file exists"
    echo "📄 File contents:"
    cat configs/.env
else
    echo "❌ Failed to create configs/.env file"
    exit 1
fi

echo ""
echo "🎯 Next steps:"
echo "1. Run database migration: make migrate-up"
echo "2. Build and start application: ./deploy_vps.sh"
echo "3. Test API: curl http://localhost:8080/health"

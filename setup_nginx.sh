#!/bin/bash

# Nginx Reverse Proxy Setup for aleph_hackthon_2
# This sets up nginx to proxy port 80 to application port 9001

echo "🌐 Setting up Nginx reverse proxy..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    print_error "Please run this script as root (sudo)"
    exit 1
fi

print_status "1. Installing nginx (if not already installed)..."
if ! command -v nginx >/dev/null 2>&1; then
    apt-get update
    apt-get install -y nginx
else
    print_status "Nginx is already installed"
fi

print_status "2. Creating nginx configuration..."
cat > /etc/nginx/sites-available/aleph_hackthon_2 << 'NGINX_CONF'
server {
    listen 80;
    server_name 13.229.87.19;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;

    # Proxy to application
    location / {
        proxy_pass http://127.0.0.1:9001\;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeout settings
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://127.0.0.1:9001/health\;
        access_log off;
    }
}
NGINX_CONF

print_status "3. Enabling the site..."
ln -sf /etc/nginx/sites-available/aleph_hackthon_2 /etc/nginx/sites-enabled/

print_status "4. Removing default nginx site (optional)..."
rm -f /etc/nginx/sites-enabled/default

print_status "5. Testing nginx configuration..."
nginx -t
if [ $? -ne 0 ]; then
    print_error "Nginx configuration test failed!"
    exit 1
fi

print_status "6. Reloading nginx..."
systemctl reload nginx

print_status "7. Checking nginx status..."
systemctl status nginx --no-pager -l

print_status "8. Testing the proxy..."
curl -I http://localhost/
curl -I http://13.229.87.19/

print_status "🎉 Nginx reverse proxy setup completed!"
echo ""
echo "📋 Configuration Summary:"
echo "- Public URL: http://13.229.87.19"
echo "- Internal App: http://127.0.0.1:9001"
echo "- Nginx Config: /etc/nginx/sites-available/aleph_hackthon_2"
echo ""
echo "🔧 Useful commands:"
echo "- Check nginx: systemctl status nginx"
echo "- Reload nginx: systemctl reload nginx"
echo "- View logs: tail -f /var/log/nginx/access.log"
echo "- Test config: nginx -t"

#!/bin/bash

# Port Forwarding Setup for VPS
# Forward public port 80 to application port 9001

echo "🔧 Setting up port forwarding from port 80 to port 9001..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

print_status "1. Checking current iptables rules..."
iptables -t nat -L PREROUTING --line-numbers

print_status "2. Adding port forwarding rule..."
# Forward port 80 to port 9001
iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 9001

print_status "3. Checking updated rules..."
iptables -t nat -L PREROUTING --line-numbers

print_status "4. Making iptables rules persistent..."

# For Ubuntu/Debian systems
if command -v netfilter-persistent >/dev/null 2>&1; then
    print_status "Using netfilter-persistent..."
    netfilter-persistent save
    netfilter-persistent reload
elif command -v iptables-persistent >/dev/null 2>&1; then
    print_status "Using iptables-persistent..."
    iptables-save > /etc/iptables/rules.v4
else
    print_warning "No persistent iptables solution found."
    print_warning "Rules will be lost after reboot."
    print_warning "Consider installing: apt-get install iptables-persistent"
fi

print_status "5. Checking firewall status..."
if command -v ufw >/dev/null 2>&1; then
    print_info "UFW status:"
    ufw status | grep 80 || print_warning "Port 80 might not be allowed in UFW"
    print_info "To allow port 80: sudo ufw allow 80"
fi

print_status "6. Testing port forwarding..."
# Test if port 80 is accessible
timeout 5 bash -c "</dev/tcp/localhost/80" && print_status "Port 80 is accessible" || print_warning "Port 80 connection test failed"

print_status "🎉 Port forwarding setup completed!"
echo ""
echo "📋 Summary:"
echo "- Public IP: 13.229.87.19"
echo "- Public Port: 80"
echo "- Internal Port: 9001"
echo "- Application: aleph_hackthon_2"
echo ""
echo "🌐 Access your application at: http://13.229.87.19"
echo ""
echo "🔧 Useful commands:"
echo "- Check rules: sudo iptables -t nat -L PREROUTING"
echo "- Remove rule: sudo iptables -t nat -D PREROUTING 1 (replace 1 with rule number)"
echo "- Restart iptables: sudo netfilter-persistent reload"

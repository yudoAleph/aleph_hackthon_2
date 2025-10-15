#!/bin/bash

# Contact Management API Test Script
# This script tests all API endpoints with example data

BASE_URL="http://localhost:8080"

echo "Testing Contact Management API"
echo "================================"

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Warning: jq is not installed. JSON responses will not be formatted."
    echo "Install jq with: brew install jq (macOS) or apt-get install jq (Linux)"
    JQ_FORMAT=""
else
    JQ_FORMAT="| jq ."
fi

# 1. Health Check
echo "1. Health Check:"
eval "curl -s -X GET \"$BASE_URL/health\" -H \"Content-Type: application/json\" $JQ_FORMAT"

# 2. Register User
echo -e "\n2. Register User:"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/register" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "email": "test@example.com",
    "phone": "+1234567890",
    "password": "password123"
  }')

if [ -n "$JQ_FORMAT" ]; then
    echo "$REGISTER_RESPONSE" | jq .
else
    echo "$REGISTER_RESPONSE"
fi

# 3. Login
echo -e "\n3. Login:"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

if [ -n "$JQ_FORMAT" ]; then
    echo "$LOGIN_RESPONSE" | jq .
else
    echo "$LOGIN_RESPONSE"
fi

# Extract token
if [ -n "$JQ_FORMAT" ]; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
else
    # Simple token extraction without jq
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
fi

echo -e "\nToken: $TOKEN"

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ] && [ "$TOKEN" != "" ]; then
    echo -e "\n4. Get Profile:"
    eval "curl -s -X GET \"$BASE_URL/api/v1/profile\" \
      -H \"Content-Type: application/json\" \
      -H \"Authorization: Bearer $TOKEN\" $JQ_FORMAT"

    echo -e "\n5. Create Contact:"
    CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/contacts" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "full_name": "Test Contact",
        "phone": "+1234567891",
        "email": "contact@example.com"
      }')

    if [ -n "$JQ_FORMAT" ]; then
        echo "$CREATE_RESPONSE" | jq .
    else
        echo "$CREATE_RESPONSE"
    fi

    echo -e "\n6. List Contacts:"
    eval "curl -s -X GET \"$BASE_URL/api/v1/contacts\" \
      -H \"Content-Type: application/json\" \
      -H \"Authorization: Bearer $TOKEN\" $JQ_FORMAT"

    echo -e "\n7. Update Profile:"
    eval "curl -s -X PUT \"$BASE_URL/api/v1/profile\" \
      -H \"Content-Type: application/json\" \
      -H \"Authorization: Bearer $TOKEN\" \
      -d '{
        \"full_name\": \"Updated Test User\",
        \"phone\": \"+1234567899\"
      }' $JQ_FORMAT"

else
    echo -e "\n❌ Failed to get authentication token. Please check if the server is running and credentials are correct."
fi

echo -e "\nAPI testing completed!"
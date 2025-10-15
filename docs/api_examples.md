# Contact Management API - cURL Examples

Base URL: <http://localhost:8080>

## Health Check

### Check API Health

```bash
curl -X GET "http://localhost:8080/health" \
  -H "Content-Type: application/json"
```

**Response:**

```json
{
  "status": "healthy"
}
```

## Authentication

### 1. User Registration

```bash
curl -X POST "http://localhost:8080/api/v1/register" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "password": "password123"
  }'
```

**Notes:**
- Email field is required and must be a valid email format
- Invalid email formats will return a 400 error with message "Invalid email format"

**Response:**

```json
{
  "status": 1,
  "status_code": 201,
  "message": "Registration success",
  "data": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "avatar_url": null,
    "token": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
}
```

### 2. User Login

```bash
curl -X POST "http://localhost:8080/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

**Notes:**
- Email field is required and must be a valid email format
- Invalid email formats will return a 400 error with message "Invalid email format"

**Response:**

```json
{
  "status": 1,
  "status_code": 200,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## Protected Routes (Require Authorization Header)

**Note:** Replace `YOUR_JWT_TOKEN` with the actual token from login response

### 3. Get User Profile

```bash
curl -X GET "http://localhost:8080/api/v1/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**

```json
{
  "status": 1,
  "status_code": 200,
  "message": "Profile loaded successfully",
  "data": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "avatar_url": null
  }
}
```

### 4. Update User Profile

```bash
curl -X PUT "http://localhost:8080/api/v1/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "full_name": "John Smith",
    "phone": "+1234567891"
  }'
```

**Response:**

```json
{
  "status": 1,
  "status_code": 200,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "full_name": "John Smith",
    "email": "john.doe@example.com",
    "phone": "+1234567891",
    "avatar_url": null
  }
}
```

## Contact Management

### 5. List Contacts (with pagination and search)

```bash
# Get all contacts
curl -X GET "http://localhost:8080/api/v1/contacts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# With search query
curl -X GET "http://localhost:8080/api/v1/contacts?q=john" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# With pagination
curl -X GET "http://localhost:8080/api/v1/contacts?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**

```json
{
  "status": 1,
  "status_code": 200,
  "message": "Contacts loaded successfully",
  "data": {
    "count": 2,
    "page": 1,
    "limit": 10,
    "contacts": [
      {
        "id": 1,
        "full_name": "Jane Smith",
        "phone": "+1234567892",
        "email": "jane@example.com",
        "favorite": true
      },
      {
        "id": 2,
        "full_name": "Bob Johnson",
        "phone": "+1234567893",
        "email": "bob@example.com",
        "favorite": false
      }
    ]
  }
}
```

### 6. Create Contact

```bash
curl -X POST "http://localhost:8080/api/v1/contacts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "full_name": "Alice Wilson",
    "phone": "+1234567894",
    "email": "alice@example.com"
  }'
```

**Notes:**
- Email field is optional but must be a valid email format if provided
- Invalid email formats will return a 400 error with message "Invalid email format"

**Response:**

```json
{
  "status": 1,
  "status_code": 201,
  "message": "Contact created successfully",
  "data": {
    "id": 3,
    "full_name": "Alice Wilson",
    "phone": "+1234567894",
    "email": "alice@example.com",
    "favorite": false
  }
}
```

### 7. Get Contact Details

```bash
curl -X GET "http://localhost:8080/api/v1/contacts/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**

```json
{
  "status": 1,
  "status_code": 200,
  "message": "Contact detail loaded",
  "data": {
    "id": 1,
    "full_name": "Jane Smith",
    "phone": "+1234567892",
    "email": "jane@example.com",
    "favorite": true
  }
}
```

### 8. Update Contact

```bash
curl -X PUT "http://localhost:8080/api/v1/contacts/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "full_name": "Jane Smith Updated",
    "phone": "+1234567895",
    "email": "jane.updated@example.com",
    "favorite": true
  }'
```

**Notes:**
- Email field is optional but must be a valid email format if provided
- Invalid email formats will return a 400 error with message "Invalid email format"

**Response:**

```json
{
  "status": 1,
  "status_code": 200,
  "message": "Contact updated successfully",
  "data": {
    "id": 1,
    "full_name": "Jane Smith Updated",
    "phone": "+1234567895",
    "email": "jane.updated@example.com",
    "favorite": true
  }
}
```

### 9. Delete Contact

```bash
curl -X DELETE "http://localhost:8080/api/v1/contacts/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**

```json
{
  "status": 1,
  "status_code": 200,
  "message": "Contact deleted successfully",
  "data": {}
}
```

## Error Responses

### Authentication Error (401)

```json
{
  "status": 0,
  "status_code": 401,
  "message": "Unauthorized",
  "data": {}
}
```

### Validation Error (400)

```json
{
  "status": 0,
  "status_code": 400,
  "message": "Invalid request format",
  "data": {
    "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
  }
}
```

### Email Validation Error (400)

```json
{
  "status": 0,
  "status_code": 400,
  "message": "Invalid email format",
  "data": {
    "error": "Email must be a valid email address"
  }
}
```

### Not Found Error (404)

```json
{
  "status": 0,
  "status_code": 404,
  "message": "Contact not found",
  "data": {}
}
```

## Quick Test Script

Create a file `test_api.sh` with the following content to test all endpoints:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "Testing Contact Management API"
echo "================================"

# 1. Health Check
echo "1. Health Check:"
curl -s -X GET "$BASE_URL/health" | jq .

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
echo $REGISTER_RESPONSE | jq .

# 3. Login
echo -e "\n3. Login:"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')
echo $LOGIN_RESPONSE | jq .

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')
echo -e "\nToken: $TOKEN"

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
  # 4. Get Profile
  echo -e "\n4. Get Profile:"
  curl -s -X GET "$BASE_URL/api/v1/profile" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" | jq .

  # 5. Create Contact
  echo -e "\n5. Create Contact:"
  curl -s -X POST "$BASE_URL/api/v1/contacts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "full_name": "Test Contact",
      "phone": "+1234567891",
      "email": "contact@example.com"
    }' | jq .

  # 6. List Contacts
  echo -e "\n6. List Contacts:"
  curl -s -X GET "$BASE_URL/api/v1/contacts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" | jq .
fi

echo -e "\nAPI testing completed!"
```

Make it executable and run:

```bash
chmod +x test_api.sh
./test_api.sh
```

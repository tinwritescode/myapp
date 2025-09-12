# Authentication Implementation

## Overview

The application now includes JWT-based authentication middleware that protects API endpoints and provides user context for URL shortener operations.

## Features

- ✅ JWT token validation
- ✅ User context extraction
- ✅ Protected and public route separation
- ✅ Configurable JWT secret
- ✅ Comprehensive error handling

## Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=myapp

# Server Configuration
SERVER_PORT=8080

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

## API Endpoints

### Public Endpoints (No Authentication Required)
- `GET /api/v1/ping` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /{short_code}` - URL redirection

### Protected Endpoints (Authentication Required)
- `GET /api/v1/users` - Get users
- `POST /api/v1/urls` - Create short URL
- `GET /api/v1/urls` - List URLs
- `GET /api/v1/urls/{id}` - Get URL details
- `PUT /api/v1/urls/{id}` - Update URL
- `DELETE /api/v1/urls/{id}` - Delete URL
- `GET /api/v1/urls/{id}/stats` - Get URL statistics

## Authentication Flow

1. **Register/Login**: Get JWT token from `/api/v1/auth/register` or `/api/v1/auth/login`
2. **Include Token**: Add `Authorization: Bearer <token>` header to protected requests
3. **Access Protected Resources**: Use the token to access user-specific endpoints

## Usage Examples

### Login Request
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Protected Request
```http
POST /api/v1/urls
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "original_url": "https://example.com/very/long/url"
}
```

## Middleware Functions

### AuthMiddleware()
- Validates JWT tokens
- Sets user context in Gin context
- Returns 401 for invalid/missing tokens

### OptionalAuthMiddleware()
- Validates JWT tokens if present
- Continues without authentication if no token
- Useful for endpoints that work with or without auth

### Helper Functions
- `GetUserID(c *gin.Context) (uint, bool)` - Extract user ID
- `GetUserEmail(c *gin.Context) (string, bool)` - Extract user email
- `GetUsername(c *gin.Context) (string, bool)` - Extract username
- `RequireAuth(c *gin.Context) bool` - Check if user is authenticated

## Security Features

- JWT token validation with expiration
- Secure secret key management via environment variables
- Proper error handling for authentication failures
- User context isolation (users can only access their own URLs)

## Testing

Use the provided HTTP request files in `requests/` directory:
- `requests/urls/complete-flow.http` - Full authentication flow
- `requests/urls/create-url.http` - URL creation with auth
- `requests/urls/get-urls.http` - URL listing with auth

## Error Codes

- `UNAUTHORIZED` - Missing or invalid authentication
- `INVALID_TOKEN` - Malformed or expired JWT token
- `TOKEN_EXPIRED` - JWT token has expired

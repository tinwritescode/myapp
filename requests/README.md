# HTTP Request Files

This directory contains HTTP request files for testing the REST API endpoints.

## Structure

```
requests/
├── environments/          # Environment-specific configurations
│   ├── local.http        # Local development
│   ├── staging.http      # Staging environment
│   └── production.http   # Production environment
├── auth/                 # Authentication endpoints
│   ├── login.http        # User login
│   └── register.http     # User registration
├── users/                # User management endpoints
│   ├── get-users.http    # Get users list
│   └── user-profile.http # User profile operations
├── shared/               # Shared utilities
│   ├── variables.http    # Global variables
│   └── health-check.http # Health check endpoints
└── README.md            # This file
```

## Usage

### VS Code REST Client Extension

1. Install the "REST Client" extension in VS Code
2. Open any `.http` file
3. Click "Send Request" above each request
4. Switch environments using the environment selector

### Environment Switching

- **Local**: `http://localhost:8080`
- **Staging**: `https://staging-api.example.com`
- **Production**: `https://api.example.com`

### Authentication

Most endpoints require authentication. The login request automatically stores the JWT token in the `authToken` variable for subsequent requests.

### Variables

- `{{baseUrl}}` - API base URL
- `{{apiVersion}}` - API version (v1)
- `{{contentType}}` - Content type (application/json)
- `{{authToken}}` - JWT token (set after login)
- `{{userId}}` - Current user ID (set after login)

## Available Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### Users
- `GET /api/v1/users` - Get users list (with pagination, search, filters)
- `GET /api/v1/users/{id}` - Get user by ID (placeholder)
- `PUT /api/v1/users/{id}` - Update user (placeholder)
- `DELETE /api/v1/users/{id}` - Delete user (placeholder)

### Health
- `GET /api/v1/ping` - Health check
- `GET /swagger/index.html` - API documentation

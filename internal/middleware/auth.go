package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tinwritescode/myapp/internal/dto/common"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWT secret key - will be set from config
var jwtSecret []byte

// SetJWTSecret sets the JWT secret from configuration
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, common.NewErrorResponseWithCode(common.UNAUTHORIZED, "Authorization header is required"))
			c.Abort()
			return
		}

		// Check if token starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, common.NewErrorResponseWithCode(common.UNAUTHORIZED, "Invalid authorization header format"))
			c.Abort()
			return
		}

		// Parse and validate token
		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.NewErrorResponseWithCode(common.INVALID_TOKEN, "Invalid or expired token"))
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
// Useful for endpoints that work with or without authentication
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Check if token starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Invalid format, continue without authentication
			c.Next()
			return
		}

		// Parse and validate token
		claims, err := validateToken(tokenString)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)

		c.Next()
	}
}

// validateToken parses and validates a JWT token
func validateToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, common.NewAppError(common.INVALID_TOKEN, "Unexpected signing method", nil)
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, common.NewAppError(common.INVALID_TOKEN, "Failed to parse token", err)
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, common.NewAppError(common.INVALID_TOKEN, "Invalid token claims", nil)
}

// RequireAuth is a helper function to check if user is authenticated
func RequireAuth(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	if uid, ok := userID.(uint); ok {
		return uid, true
	}

	return 0, false
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	if emailStr, ok := email.(string); ok {
		return emailStr, true
	}

	return "", false
}

// GetUsername extracts username from context
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("user_username")
	if !exists {
		return "", false
	}

	if usernameStr, ok := username.(string); ok {
		return usernameStr, true
	}

	return "", false
}

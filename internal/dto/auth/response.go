package auth

import "time"

// UserInfo represents user information in responses
type UserInfo struct {
	ID       uint   `json:"id" example:"1"`
	Email    string `json:"email" example:"user@example.com"`
	Username string `json:"username" example:"johndoe"`
	FullName string `json:"full_name" example:"John Doe"`
	IsActive bool   `json:"is_active" example:"true"`
}

// RegisterResponse represents the response for user registration
type RegisterResponse struct {
	User UserInfo `json:"user"`
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	Token        string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refresh_token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt    time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`
	User         UserInfo  `json:"user"`
}

// RefreshTokenResponse represents the response for token refresh
type RefreshTokenResponse struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`
}

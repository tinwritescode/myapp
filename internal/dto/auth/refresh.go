package auth

import "time"

// RefreshTokenRequest represents the request for refreshing tokens
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"a1b2c3d4e5f6..."`
}

// RefreshTokenResponse represents the response for refreshing tokens
type RefreshTokenResponse struct {
	Token        string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refresh_token" example:"a1b2c3d4e5f6..."`
	ExpiresAt    time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`
	User         UserInfo  `json:"user"`
}

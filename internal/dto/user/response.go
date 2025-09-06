package user

import (
	"time"

	"github.com/tinwritescode/myapp/internal/dto/common"
)

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Email     string    `json:"email" example:"user@example.com"`
	Username  string    `json:"username" example:"johndoe"`
	FullName  string    `json:"full_name" example:"John Doe"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}

// GetUsersResponse represents the response for getting users
type GetUsersResponse struct {
	common.PaginatedResponse
	Data []UserResponse `json:"data"`
}

// GetUserResponse represents the response for getting a single user
type GetUserResponse struct {
	common.BaseResponse
	Data UserResponse `json:"data"`
}

// UpdateUserResponse represents the response for updating a user
type UpdateUserResponse struct {
	common.BaseResponse
	Data UserResponse `json:"data"`
}

package url

import "time"

// CreateURLRequest represents the request to create a new URL
type CreateURLRequest struct {
	OriginalURL string     `json:"original_url" binding:"required,url" example:"https://example.com/very/long/url"`
	ShortCode   *string    `json:"short_code,omitempty" binding:"omitempty,alphanum,min=3,max=8" example:"abc123"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

// GetURLsRequest represents the request to get URLs with pagination and filtering
type GetURLsRequest struct {
	Page     int     `form:"page" binding:"min=1" example:"1"`
	Limit    int     `form:"limit" binding:"min=1,max=100" example:"10"`
	Search   *string `form:"search" example:"example"`
	IsActive *bool   `form:"is_active" example:"true"`
	SortBy   string  `form:"sort_by" binding:"omitempty,oneof=created_at updated_at click_count" example:"created_at"`
	SortDir  string  `form:"sort_dir" binding:"omitempty,oneof=asc desc" example:"desc"`
}

// UpdateURLRequest represents the request to update a URL
type UpdateURLRequest struct {
	OriginalURL *string    `json:"original_url,omitempty" binding:"omitempty,url" example:"https://example.com/updated/url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	IsActive    *bool      `json:"is_active,omitempty" example:"true"`
}

// RedirectRequest represents the request for URL redirection
type RedirectRequest struct {
	ShortCode string `uri:"short_code" binding:"required,alphanum" example:"abc123"`
}

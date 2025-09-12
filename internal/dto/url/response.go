package url

import (
	"time"

	"github.com/tinwritescode/myapp/internal/dto/common"
)

// URLResponse represents a URL in API responses
type URLResponse struct {
	ID          uint       `json:"id" example:"1"`
	OriginalURL string     `json:"original_url" example:"https://example.com/very/long/url"`
	ShortCode   string     `json:"short_code" example:"abc123"`
	UserID      *uint      `json:"user_id,omitempty" example:"1"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	ClickCount  int64      `json:"click_count" example:"42"`
	IsActive    bool       `json:"is_active" example:"true"`
	CreatedAt   time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// CreateURLResponse represents the response when creating a URL
type CreateURLResponse struct {
	common.BaseResponse
	Data URLResponse `json:"data"`
}

// GetURLsResponse represents the response when getting URLs with pagination
type GetURLsResponse struct {
	common.PaginatedResponse
	Data []URLResponse `json:"data"`
}

// GetURLResponse represents the response when getting a single URL
type GetURLResponse struct {
	common.BaseResponse
	Data URLResponse `json:"data"`
}

// UpdateURLResponse represents the response when updating a URL
type UpdateURLResponse struct {
	common.BaseResponse
	Data URLResponse `json:"data"`
}

// DeleteURLResponse represents the response when deleting a URL
type DeleteURLResponse struct {
	common.BaseResponse
}

// URLStatsResponse represents the response for URL statistics
type URLStatsResponse struct {
	common.BaseResponse
	Data URLStats `json:"data"`
}

// URLStats represents URL statistics
type URLStats struct {
	URLResponse
	RecentClicks []ClickEvent `json:"recent_clicks,omitempty"`
}

// ClickEvent represents a click event
type ClickEvent struct {
	ID        uint      `json:"id" example:"1"`
	URLID     uint      `json:"url_id" example:"1"`
	IPAddress string    `json:"ip_address" example:"192.168.1.1"`
	UserAgent string    `json:"user_agent" example:"Mozilla/5.0..."`
	Referer   *string   `json:"referer,omitempty" example:"https://google.com"`
	ClickedAt time.Time `json:"clicked_at" example:"2024-01-01T00:00:00Z"`
}

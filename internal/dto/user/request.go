package user

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" binding:"omitempty,email" example:"user@example.com"`
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=20" example:"johndoe"`
	FullName *string `json:"full_name,omitempty" example:"John Doe"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

// GetUsersRequest represents query parameters for getting users
type GetUsersRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1" example:"1"`
	Limit    int    `form:"limit,default=10" binding:"min=1,max=100" example:"10"`
	Search   string `form:"search" example:"john"`
	IsActive *bool  `form:"is_active" example:"true"`
	SortBy   string `form:"sort_by,default=created_at" example:"created_at"`
	SortDir  string `form:"sort_dir,default=desc" binding:"oneof=asc desc" example:"desc"`
}

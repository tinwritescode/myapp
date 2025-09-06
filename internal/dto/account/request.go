package account

// CreateAccountRequest represents the request body for creating an account
type CreateAccountRequest struct {
	UserID      uint   `json:"user_id" binding:"required" example:"1"`
	AccountType string `json:"account_type" binding:"required,oneof=checking savings credit" example:"checking"`
	Balance     int64  `json:"balance,omitempty" example:"1000"`
	Currency    string `json:"currency,omitempty" binding:"omitempty,len=3" example:"USD"`
	IsActive    *bool  `json:"is_active,omitempty" example:"true"`
}

// UpdateAccountRequest represents the request body for updating an account
type UpdateAccountRequest struct {
	AccountType *string `json:"account_type,omitempty" binding:"omitempty,oneof=checking savings credit" example:"checking"`
	Balance     *int64  `json:"balance,omitempty" example:"1000"`
	Currency    *string `json:"currency,omitempty" binding:"omitempty,len=3" example:"USD"`
	IsActive    *bool   `json:"is_active,omitempty" example:"true"`
}

// TransferRequest represents the request body for transferring money
type TransferRequest struct {
	FromAccountID uint   `json:"from_account_id" binding:"required" example:"1"`
	ToAccountID   uint   `json:"to_account_id" binding:"required" example:"2"`
	Amount        int64  `json:"amount" binding:"required,gt=0" example:"100"`
	Description   string `json:"description,omitempty" example:"Transfer to savings"`
}

// GetAccountsRequest represents query parameters for getting accounts
type GetAccountsRequest struct {
	UserID      *uint  `form:"user_id" example:"1"`
	AccountType string `form:"account_type" example:"checking"`
	IsActive    *bool  `form:"is_active" example:"true"`
	Page        int    `form:"page,default=1" binding:"min=1" example:"1"`
	Limit       int    `form:"limit,default=10" binding:"min=1,max=100" example:"10"`
	SortBy      string `form:"sort_by,default=created_at" example:"created_at"`
	SortDir     string `form:"sort_dir,default=desc" binding:"oneof=asc desc" example:"desc"`
}

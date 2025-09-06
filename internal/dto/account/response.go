package account

import (
	"time"

	"github.com/tinwritescode/myapp/internal/dto/common"
)

// AccountResponse represents an account in API responses
type AccountResponse struct {
	ID          uint      `json:"id" example:"1"`
	UserID      uint      `json:"user_id" example:"1"`
	AccountType string    `json:"account_type" example:"checking"`
	Balance     int64     `json:"balance" example:"1000"`
	Currency    string    `json:"currency" example:"USD"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}

// GetAccountsResponse represents the response for getting accounts
type GetAccountsResponse struct {
	common.PaginatedResponse
	Data []AccountResponse `json:"data"`
}

// CreateAccountResponse represents the response for creating an account
type CreateAccountResponse struct {
	common.BaseResponse
	Data AccountResponse `json:"data"`
}

// GetAccountResponse represents the response for getting a single account
type GetAccountResponse struct {
	common.BaseResponse
	Data AccountResponse `json:"data"`
}

// UpdateAccountResponse represents the response for updating an account
type UpdateAccountResponse struct {
	common.BaseResponse
	Data AccountResponse `json:"data"`
}

// TransferResponse represents the response for money transfer
type TransferResponse struct {
	common.BaseResponse
	Data struct {
		TransactionID string          `json:"transaction_id" example:"txn_123456"`
		FromAccount   AccountResponse `json:"from_account"`
		ToAccount     AccountResponse `json:"to_account"`
		Amount        int64           `json:"amount" example:"100"`
		Description   string          `json:"description" example:"Transfer to savings"`
		ProcessedAt   time.Time       `json:"processed_at" example:"2024-01-01T12:00:00Z"`
	} `json:"data"`
}

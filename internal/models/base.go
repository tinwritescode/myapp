package models

import (
	"time"

	"github.com/tinwritescode/myapp/internal/dto/user"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	BaseModel
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	FullName string `json:"full_name"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// ToResponse converts User model to UserResponse DTO
func (u *User) ToResponse() user.UserResponse {
	return user.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type Account struct {
	BaseModel
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	AccountType string `gorm:"not null" json:"account_type"`
	Balance     int64  `gorm:"default:0" json:"balance"`
	Currency    string `gorm:"default:'USD'" json:"currency"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tinwritescode/myapp/internal/database"
	"github.com/tinwritescode/myapp/internal/dto/common"
	"github.com/tinwritescode/myapp/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(email, username, password, fullName string) (string, string, time.Time, *models.User, error)
	Login(email, password string) (string, string, time.Time, *models.User, error)
	RefreshToken(refreshToken string) (string, string, time.Time, *models.User, error)
	RevokeRefreshToken(refreshToken string) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	db *gorm.DB
}

// JWT secret key - will be set from config
var jwtSecret []byte

// SetJWTSecret sets the JWT secret from configuration
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	userServiceInstance UserService
)

func NewUserService() UserService {
	return &userService{
		db: database.GetDB(),
	}
}

func GetUserService() UserService {
	if userServiceInstance == nil {
		userServiceInstance = NewUserService()
	}
	return userServiceInstance
}

func (s *userService) Register(email, username, password, fullName string) (string, string, time.Time, *models.User, error) {
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", email, username).First(&existingUser).Error; err == nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.EMAIL_ALREADY_USED, "user with this email or username already exists", nil)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to process password", err)
	}

	// Create user
	user := models.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
		FullName: fullName,
		IsActive: true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to create user", err)
	}

	// Generate JWT token
	token, err := s.generateJWT(&user)
	if err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to generate token", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to generate refresh token", err)
	}

	// Set expiration time (7 days for refresh token)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Store refresh token in database
	refreshTokenRecord := models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expirationTime,
		IsActive:  true,
	}

	if err := s.db.Create(&refreshTokenRecord).Error; err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to store refresh token", err)
	}

	return token, refreshToken, expirationTime, &user, nil
}

func (s *userService) Login(email, password string) (string, string, time.Time, *models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INVALID_CREDENTIALS, "invalid credentials", err)
	}

	if !user.IsActive {
		return "", "", time.Time{}, nil, common.NewAppError(common.UNAUTHORIZED, "account is deactivated", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INVALID_CREDENTIALS, "invalid credentials", err)
	}

	token, err := s.generateJWT(&user)
	if err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to generate token", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to generate refresh token", err)
	}

	// Set expiration time (7 days for refresh token)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Store refresh token in database
	refreshTokenRecord := models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expirationTime,
		IsActive:  true,
	}

	if err := s.db.Create(&refreshTokenRecord).Error; err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to store refresh token", err)
	}

	return token, refreshToken, expirationTime, &user, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, common.NewAppError(common.USER_NOT_FOUND, "user not found", err)
	}
	return &user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, common.NewAppError(common.USER_NOT_FOUND, "user not found", err)
	}
	return &user, nil
}

func (s *userService) generateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// RefreshToken validates a refresh token and returns new access and refresh tokens
func (s *userService) RefreshToken(refreshToken string) (string, string, time.Time, *models.User, error) {
	// Find the refresh token in database
	var tokenRecord models.RefreshToken
	if err := s.db.Where("token = ? AND is_active = ?", refreshToken, true).First(&tokenRecord).Error; err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INVALID_TOKEN, "invalid refresh token", err)
	}

	// Check if token is expired
	if tokenRecord.IsExpired() {
		// Mark token as inactive
		s.db.Model(&tokenRecord).Update("is_active", false)
		return "", "", time.Time{}, nil, common.NewAppError(common.TOKEN_EXPIRED, "refresh token expired", nil)
	}

	// Get the user
	var user models.User
	if err := s.db.First(&user, tokenRecord.UserID).Error; err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.USER_NOT_FOUND, "user not found", err)
	}

	// Check if user is active
	if !user.IsActive {
		return "", "", time.Time{}, nil, common.NewAppError(common.UNAUTHORIZED, "account is deactivated", nil)
	}

	// Generate new access token
	newAccessToken, err := s.generateJWT(&user)
	if err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to generate access token", err)
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to generate refresh token", err)
	}

	// Set expiration time (7 days for refresh token)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Invalidate old refresh token
	if err := s.db.Model(&tokenRecord).Update("is_active", false).Error; err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to invalidate old token", err)
	}

	// Store new refresh token in database
	newTokenRecord := models.RefreshToken{
		Token:     newRefreshToken,
		UserID:    user.ID,
		ExpiresAt: expirationTime,
		IsActive:  true,
	}

	if err := s.db.Create(&newTokenRecord).Error; err != nil {
		return "", "", time.Time{}, nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to store new refresh token", err)
	}

	return newAccessToken, newRefreshToken, expirationTime, &user, nil
}

// RevokeRefreshToken invalidates a refresh token
func (s *userService) RevokeRefreshToken(refreshToken string) error {
	var tokenRecord models.RefreshToken
	if err := s.db.Where("token = ?", refreshToken).First(&tokenRecord).Error; err != nil {
		return common.NewAppError(common.INVALID_TOKEN, "refresh token not found", err)
	}

	if err := s.db.Model(&tokenRecord).Update("is_active", false).Error; err != nil {
		return common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to revoke token", err)
	}

	return nil
}

// generateRefreshToken generates a cryptographically secure random refresh token
func (s *userService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

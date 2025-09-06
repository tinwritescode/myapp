package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tinwritescode/myapp/internal/database"
	"github.com/tinwritescode/myapp/internal/dto/common"
	"github.com/tinwritescode/myapp/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(email, username, password, fullName string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsers(page, limit int, search *string, isActive *bool, sortBy, sortDir string) ([]models.User, int64, error)
}

type userService struct {
	db *gorm.DB
}

var jwtSecret = []byte("your-secret-key")

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

func (s *userService) Register(email, username, password, fullName string) (*models.User, error) {
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", email, username).First(&existingUser).Error; err == nil {
		return nil, common.NewAppError(common.EMAIL_ALREADY_USED, "user with this email or username already exists", nil)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to process password", err)
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
		return nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to create user", err)
	}

	return &user, nil
}

func (s *userService) Login(email, password string) (string, *models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", nil, common.NewAppError(common.INVALID_CREDENTIALS, "invalid credentials", err)
	}

	if !user.IsActive {
		return "", nil, common.NewAppError(common.UNAUTHORIZED, "account is deactivated", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, common.NewAppError(common.INVALID_CREDENTIALS, "invalid credentials", err)
	}

	token, err := s.generateJWT(&user)
	if err != nil {
		return "", nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "failed to generate token", err)
	}

	return token, &user, nil
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

func (s *userService) GetUsers(page, limit int, search *string, isActive *bool, sortBy, sortDir string) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	query := s.db

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if search != nil && *search != "" {
		query = query.Where("email LIKE ? OR username LIKE ?", "%"+*search+"%", "%"+*search+"%")
	}

	// Count total records
	if err := query.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortDir)).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

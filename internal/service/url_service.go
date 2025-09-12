package service

import (
	"fmt"
	"time"

	"github.com/tinwritescode/myapp/internal/database"
	"github.com/tinwritescode/myapp/internal/dto/common"
	"github.com/tinwritescode/myapp/internal/models"
	"github.com/tinwritescode/myapp/pkg/utils"
	"gorm.io/gorm"
)

type URLService interface {
	CreateURL(originalURL string, shortCode *string, userID *uint, expiresAt *time.Time) (*models.URL, error)
	GetURLByShortCode(shortCode string) (*models.URL, error)
	GetURLByID(id uint, userID *uint) (*models.URL, error)
	GetURLs(userID *uint, page, limit int, search *string, isActive *bool, sortBy, sortDir string) ([]models.URL, int64, error)
	UpdateURL(id uint, userID *uint, originalURL *string, expiresAt *time.Time, isActive *bool) (*models.URL, error)
	DeleteURL(id uint, userID *uint) error
	IncrementClickCount(shortCode string) error
	GetURLStats(id uint, userID *uint) (*models.URL, error)
}

type urlService struct {
	db *gorm.DB
}

var (
	urlServiceInstance URLService
)

func NewURLService() URLService {
	return &urlService{
		db: database.GetDB(),
	}
}

func GetURLService() URLService {
	if urlServiceInstance == nil {
		urlServiceInstance = NewURLService()
	}
	return urlServiceInstance
}

func (s *urlService) CreateURL(originalURL string, shortCode *string, userID *uint, expiresAt *time.Time) (*models.URL, error) {
	// Validate original URL
	if err := utils.ValidateURL(originalURL); err != nil {
		return nil, common.NewAppError(common.VALIDATION_ERROR, fmt.Sprintf("Invalid URL: %s", err.Error()), err)
	}

	// Normalize URL
	normalizedURL := utils.NormalizeURL(originalURL)

	// Generate short code if not provided
	var finalShortCode string
	if shortCode != nil {
		if err := utils.ValidateShortCode(*shortCode); err != nil {
			return nil, common.NewAppError(common.VALIDATION_ERROR, fmt.Sprintf("Invalid short code: %s", err.Error()), err)
		}
		finalShortCode = *shortCode
	} else {
		generatedCode, err := utils.GenerateShortCode(utils.ShortCodeLength)
		if err != nil {
			return nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to generate short code", err)
		}
		finalShortCode = generatedCode
	}

	// Check if short code already exists
	var existingURL models.URL
	if err := s.db.Where("short_code = ?", finalShortCode).First(&existingURL).Error; err == nil {
		if shortCode != nil {
			return nil, common.NewAppError(common.SHORT_CODE_ALREADY_EXISTS, "Short code already exists", nil)
		}
		// If auto-generated, try again with a different code
		generatedCode, err := utils.GenerateShortCode(utils.ShortCodeLength)
		if err != nil {
			return nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to generate unique short code", err)
		}
		finalShortCode = generatedCode
	}

	// Check if URL already exists for the same user
	if userID != nil {
		var existingUserURL models.URL
		if err := s.db.Where("original_url = ? AND user_id = ?", normalizedURL, *userID).First(&existingUserURL).Error; err == nil {
			return &existingUserURL, nil // Return existing URL
		}
	}

	// Create URL
	url := models.URL{
		OriginalURL: normalizedURL,
		ShortCode:   finalShortCode,
		UserID:      userID,
		ExpiresAt:   expiresAt,
		ClickCount:  0,
		IsActive:    true,
	}

	if err := s.db.Create(&url).Error; err != nil {
		return nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to create URL", err)
	}

	return &url, nil
}

func (s *urlService) GetURLByShortCode(shortCode string) (*models.URL, error) {
	var url models.URL
	if err := s.db.Where("short_code = ? AND is_active = ?", shortCode, true).First(&url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewAppError(common.URL_NOT_FOUND, "URL not found", err)
		}
		return nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to get URL", err)
	}

	// Check if URL has expired
	if utils.IsURLExpired(url.ExpiresAt) {
		return nil, common.NewAppError(common.URL_EXPIRED, "URL has expired", nil)
	}

	return &url, nil
}

func (s *urlService) GetURLByID(id uint, userID *uint) (*models.URL, error) {
	var url models.URL
	query := s.db.Where("id = ?", id)

	// If userID is provided, ensure user owns the URL
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.First(&url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewAppError(common.URL_NOT_FOUND, "URL not found", err)
		}
		return nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to get URL", err)
	}

	return &url, nil
}

func (s *urlService) GetURLs(userID *uint, page, limit int, search *string, isActive *bool, sortBy, sortDir string) ([]models.URL, int64, error) {
	var urls []models.URL
	var total int64
	query := s.db

	// Filter by user if provided
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Filter by active status
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Search functionality
	if search != nil && *search != "" {
		query = query.Where("original_url LIKE ? OR short_code LIKE ?", "%"+*search+"%", "%"+*search+"%")
	}

	// Count total records
	if err := query.Model(&models.URL{}).Count(&total).Error; err != nil {
		return nil, 0, common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to count URLs", err)
	}

	// Get paginated results
	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortDir)).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&urls).Error; err != nil {
		return nil, 0, common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to get URLs", err)
	}

	return urls, total, nil
}

func (s *urlService) UpdateURL(id uint, userID *uint, originalURL *string, expiresAt *time.Time, isActive *bool) (*models.URL, error) {
	// Get existing URL
	url, err := s.GetURLByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if originalURL != nil {
		if err := utils.ValidateURL(*originalURL); err != nil {
			return nil, common.NewAppError(common.VALIDATION_ERROR, fmt.Sprintf("Invalid URL: %s", err.Error()), err)
		}
		url.OriginalURL = utils.NormalizeURL(*originalURL)
	}

	if expiresAt != nil {
		url.ExpiresAt = expiresAt
	}

	if isActive != nil {
		url.IsActive = *isActive
	}

	// Save changes
	if err := s.db.Save(url).Error; err != nil {
		return nil, common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to update URL", err)
	}

	return url, nil
}

func (s *urlService) DeleteURL(id uint, userID *uint) error {
	// Check if URL exists and user has permission
	_, err := s.GetURLByID(id, userID)
	if err != nil {
		return err
	}

	// Soft delete
	if err := s.db.Delete(&models.URL{}, id).Error; err != nil {
		return common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to delete URL", err)
	}

	return nil
}

func (s *urlService) IncrementClickCount(shortCode string) error {
	if err := s.db.Model(&models.URL{}).Where("short_code = ?", shortCode).Update("click_count", gorm.Expr("click_count + 1")).Error; err != nil {
		return common.NewAppError(common.INTERNAL_SERVER_ERROR, "Failed to increment click count", err)
	}
	return nil
}

func (s *urlService) GetURLStats(id uint, userID *uint) (*models.URL, error) {
	// Get URL with stats
	url, err := s.GetURLByID(id, userID)
	if err != nil {
		return nil, err
	}

	return url, nil
}

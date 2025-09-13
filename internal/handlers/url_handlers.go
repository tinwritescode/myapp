package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tinwritescode/myapp/internal/dto/common"
	"github.com/tinwritescode/myapp/internal/dto/url"
	"github.com/tinwritescode/myapp/internal/middleware"
	"github.com/tinwritescode/myapp/internal/service"
	"github.com/tinwritescode/myapp/pkg/utils"
)

func getURLService() service.URLService {
	return service.GetURLService()
}

// @Summary Create URL
// @Description Create a new short URL
// @Tags urls
// @Accept json
// @Produce json
// @Param request body url.CreateURLRequest true "URL creation details"
// @Success 201 {object} url.CreateURLResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 409 {object} common.ErrorResponse
// @Router /urls [post]
func CreateURL(c *gin.Context) {
	var req url.CreateURLRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	// Get user ID from context if authenticated
	var userID *uint
	if uid, exists := middleware.GetUserID(c); exists {
		userID = &uid
	}

	urlService := getURLService()
	createdURL, err := urlService.CreateURL(req.OriginalURL, req.ShortCode, userID, req.ExpiresAt)
	if err != nil {
		handleURLError(c, err)
		return
	}

	response := url.CreateURLResponse{
		BaseResponse: common.BaseResponse{
			Success: true,
			Message: "URL created successfully",
		},
		Data: createdURL.ToResponse(),
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Create Public URL
// @Description Create a new short URL without authentication
// @Tags urls
// @Accept json
// @Produce json
// @Param request body url.CreateURLRequest true "URL creation details"
// @Success 201 {object} url.CreateURLResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 409 {object} common.ErrorResponse
// @Router /urls/public [post]
func CreatePublicURL(c *gin.Context) {
	var req url.CreateURLRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	// Get user ID from context if authenticated
	var userID *uint
	if uid, exists := middleware.GetUserID(c); exists {
		userID = &uid
	}

	// Create URL with user ID if authenticated, otherwise public
	urlService := getURLService()
	createdURL, err := urlService.CreateURL(req.OriginalURL, req.ShortCode, userID, req.ExpiresAt)
	if err != nil {
		handleURLError(c, err)
		return
	}

	response := url.CreateURLResponse{
		BaseResponse: common.BaseResponse{
			Success: true,
			Message: "URL created successfully",
		},
		Data: createdURL.ToResponse(),
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Get URLs
// @Description Get URLs with pagination and filtering
// @Tags urls
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param is_active query bool false "Filter by active status"
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_dir query string false "Sort direction" Enums(asc, desc) default(desc)
// @Success 200 {object} url.GetURLsResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /urls [get]
func GetURLs(c *gin.Context) {
	var req url.GetURLsRequest
	if !middleware.BindQuery(c, &req) {
		return
	}

	// Get user ID from context if authenticated
	var userID *uint
	if uid, exists := middleware.GetUserID(c); exists {
		userID = &uid
	}

	urlService := getURLService()
	urls, total, err := urlService.GetURLs(userID, req.Page, req.Limit, req.Search, req.IsActive, req.SortBy, req.SortDir)
	if err != nil {
		handleURLError(c, err)
		return
	}

	// Map to URL response
	urlResponses := make([]url.URLResponse, len(urls))
	for i, u := range urls {
		urlResponses[i] = u.ToResponse()
	}

	totalPages := utils.CalculateTotalPages(int(total), req.Limit)
	response := url.GetURLsResponse{
		PaginatedResponse: common.PaginatedResponse{
			BaseResponse: common.BaseResponse{
				Success: true,
				Message: "URLs retrieved successfully",
			},
			Pagination: common.Pagination{
				Page:       req.Page,
				Limit:      req.Limit,
				Total:      total,
				TotalPages: totalPages,
			},
		},
		Data: urlResponses,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get URL by ID
// @Description Get a specific URL by ID
// @Tags urls
// @Accept json
// @Produce json
// @Param id path int true "URL ID"
// @Success 200 {object} url.GetURLResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /urls/{id} [get]
func GetURLByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid URL ID"))
		return
	}

	// Get user ID from context if authenticated
	var userID *uint
	if uid, exists := middleware.GetUserID(c); exists {
		userID = &uid
	}

	urlService := getURLService()
	urlData, err := urlService.GetURLByID(uint(id), userID)
	if err != nil {
		handleURLError(c, err)
		return
	}

	response := url.GetURLResponse{
		BaseResponse: common.BaseResponse{
			Success: true,
			Message: "URL retrieved successfully",
		},
		Data: urlData.ToResponse(),
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update URL
// @Description Update a URL
// @Tags urls
// @Accept json
// @Produce json
// @Param id path int true "URL ID"
// @Param request body url.UpdateURLRequest true "URL update details"
// @Success 200 {object} url.UpdateURLResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /urls/{id} [put]
func UpdateURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid URL ID"))
		return
	}

	var req url.UpdateURLRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	// Get user ID from context if authenticated
	var userID *uint
	if uid, exists := middleware.GetUserID(c); exists {
		userID = &uid
	}

	urlService := getURLService()
	updatedURL, err := urlService.UpdateURL(uint(id), userID, req.OriginalURL, req.ExpiresAt, req.IsActive)
	if err != nil {
		handleURLError(c, err)
		return
	}

	response := url.UpdateURLResponse{
		BaseResponse: common.BaseResponse{
			Success: true,
			Message: "URL updated successfully",
		},
		Data: updatedURL.ToResponse(),
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Delete URL
// @Description Delete a URL
// @Tags urls
// @Accept json
// @Produce json
// @Param id path int true "URL ID"
// @Success 200 {object} url.DeleteURLResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /urls/{id} [delete]
func DeleteURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid URL ID"))
		return
	}

	// Get user ID from context if authenticated
	var userID *uint
	if uid, exists := middleware.GetUserID(c); exists {
		userID = &uid
	}

	urlService := getURLService()
	err = urlService.DeleteURL(uint(id), userID)
	if err != nil {
		handleURLError(c, err)
		return
	}

	response := url.DeleteURLResponse{
		BaseResponse: common.BaseResponse{
			Success: true,
			Message: "URL deleted successfully",
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Redirect to original URL
// @Description Redirect to the original URL using short code
// @Tags urls
// @Param short_code path string true "Short code"
// @Success 302 {string} string "Redirect to original URL"
// @Failure 404 {object} common.ErrorResponse
// @Router /{short_code} [get]
func RedirectURL(c *gin.Context) {
	shortCode := c.Param("short_code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Short code is required"))
		return
	}

	urlService := getURLService()
	urlData, err := urlService.GetURLByShortCode(shortCode)
	if err != nil {
		handleURLError(c, err)
		return
	}

	// Increment click count
	if err := urlService.IncrementClickCount(shortCode); err != nil {
		// Log error but don't fail the redirect
		// You might want to use a proper logger here
	}

	// Redirect to original URL
	c.Redirect(http.StatusFound, urlData.OriginalURL)
}

// @Summary Get URL statistics
// @Description Get statistics for a specific URL
// @Tags urls
// @Accept json
// @Produce json
// @Param id path int true "URL ID"
// @Success 200 {object} url.URLStatsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /urls/{id}/stats [get]
func GetURLStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid URL ID"))
		return
	}

	// Get user ID from context if authenticated
	var userID *uint
	if uid, exists := middleware.GetUserID(c); exists {
		userID = &uid
	}

	urlService := getURLService()
	urlData, err := urlService.GetURLStats(uint(id), userID)
	if err != nil {
		handleURLError(c, err)
		return
	}

	response := url.URLStatsResponse{
		BaseResponse: common.BaseResponse{
			Success: true,
			Message: "URL statistics retrieved successfully",
		},
		Data: url.URLStats{
			URLResponse: urlData.ToResponse(),
			// RecentClicks would be populated from a separate table in a real implementation
		},
	}

	c.JSON(http.StatusOK, response)
}

// handleURLError handles URL-specific errors
func handleURLError(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	if appErr, ok := err.(*common.AppError); ok {
		switch appErr.Code {
		case common.VALIDATION_ERROR:
			statusCode = http.StatusBadRequest
		case common.SHORT_CODE_ALREADY_EXISTS:
			statusCode = http.StatusConflict
		case common.URL_NOT_FOUND:
			statusCode = http.StatusNotFound
		case common.URL_EXPIRED:
			statusCode = http.StatusGone
		case common.UNAUTHORIZED:
			statusCode = http.StatusUnauthorized
		case common.FORBIDDEN:
			statusCode = http.StatusForbidden
		case common.INTERNAL_SERVER_ERROR:
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, common.NewErrorResponseWithCode(appErr.Code, appErr.Message))
	} else {
		c.JSON(statusCode, common.NewErrorResponse(err.Error()))
	}
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinwritescode/myapp/internal/dto/auth"
	"github.com/tinwritescode/myapp/internal/dto/common"
	"github.com/tinwritescode/myapp/internal/middleware"
	"github.com/tinwritescode/myapp/internal/service"
)

func getUserService() service.UserService {
	return service.GetUserService()
}

// @Summary Register user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Registration details"
// @Success 201 {object} auth.RegisterResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 409 {object} common.ErrorResponse
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req auth.RegisterRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	userService := getUserService()
	token, refreshToken, expiresAt, user, err := userService.Register(req.Email, req.Username, req.Password, req.FullName)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if appErr, ok := err.(*common.AppError); ok {
			switch appErr.Code {
			case common.EMAIL_ALREADY_USED:
				statusCode = http.StatusConflict
			case common.INTERNAL_SERVER_ERROR:
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, common.NewErrorResponseWithCode(appErr.Code, appErr.Message))
		} else {
			c.JSON(statusCode, common.NewErrorResponse(err.Error()))
		}
		return
	}

	// Return user data with token (without password)
	response := auth.RegisterResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: auth.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			FullName: user.FullName,
			IsActive: user.IsActive,
		},
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login credentials"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req auth.LoginRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	// Use service to login user
	userService := getUserService()
	token, refreshToken, expiresAt, user, err := userService.Login(req.Email, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if appErr, ok := err.(*common.AppError); ok {
			switch appErr.Code {
			case common.INVALID_CREDENTIALS, common.UNAUTHORIZED:
				statusCode = http.StatusUnauthorized
			case common.INTERNAL_SERVER_ERROR:
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, common.NewErrorResponseWithCode(appErr.Code, appErr.Message))
		} else {
			c.JSON(statusCode, common.NewErrorResponse(err.Error()))
		}
		return
	}

	// Return response
	response := auth.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: auth.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			FullName: user.FullName,
			IsActive: user.IsActive,
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} auth.RefreshTokenResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	userService := getUserService()
	token, refreshToken, expiresAt, user, err := userService.RefreshToken(req.RefreshToken)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if appErr, ok := err.(*common.AppError); ok {
			switch appErr.Code {
			case common.INVALID_TOKEN, common.TOKEN_EXPIRED:
				statusCode = http.StatusUnauthorized
			case common.UNAUTHORIZED:
				statusCode = http.StatusUnauthorized
			case common.USER_NOT_FOUND:
				statusCode = http.StatusUnauthorized
			case common.INTERNAL_SERVER_ERROR:
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, common.NewErrorResponseWithCode(appErr.Code, appErr.Message))
		} else {
			c.JSON(statusCode, common.NewErrorResponse(err.Error()))
		}
		return
	}

	// Return response
	response := auth.RefreshTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: auth.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			FullName: user.FullName,
			IsActive: user.IsActive,
		},
	}

	c.JSON(http.StatusOK, response)
}

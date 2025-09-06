package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinwritescode/myapp/internal/dto/auth"
	"github.com/tinwritescode/myapp/internal/dto/common"
	"github.com/tinwritescode/myapp/internal/dto/user"
	"github.com/tinwritescode/myapp/internal/middleware"
	"github.com/tinwritescode/myapp/internal/service"
	"github.com/tinwritescode/myapp/pkg/utils"
)

func getUserService() service.UserService {
	return service.GetUserService()
}

// @Summary Get users
// @Description Get users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param is_active query bool false "Filter by active status"
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_dir query string false "Sort direction" Enums(asc, desc) default(desc)
// @Success 200 {object} user.GetUsersResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var req user.GetUsersRequest
	if !middleware.BindQuery(c, &req) {
		return
	}

	userService := getUserService()
	users, total, err := userService.GetUsers(req.Page, req.Limit, &req.Search, req.IsActive, req.SortBy, req.SortDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

	// Map to user response
	userResponses := make([]user.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = u.ToResponse()
	}

	totalPages := utils.CalculateTotalPages(int(total), req.Limit)
	response := user.GetUsersResponse{
		PaginatedResponse: common.PaginatedResponse{
			BaseResponse: common.BaseResponse{
				Success: true,
				Message: "Users retrieved successfully",
			},
			Pagination: common.Pagination{
				Page:       req.Page,
				Limit:      req.Limit,
				Total:      total,
				TotalPages: totalPages,
			},
		},
		Data: userResponses,
	}

	c.JSON(http.StatusOK, response)
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
	user, err := userService.Register(req.Email, req.Username, req.Password, req.FullName)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with this email or username already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, common.NewErrorResponse(err.Error()))
		return
	}

	// Return user data (without password)
	response := auth.RegisterResponse{
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
	token, user, err := userService.Login(req.Email, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid credentials" || err.Error() == "account is deactivated" {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, common.NewErrorResponse(err.Error()))
		return
	}

	// Return response
	response := auth.LoginResponse{
		Token: token,
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

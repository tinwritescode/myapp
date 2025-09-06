package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinwritescode/myapp/internal/dto/common"
)

// BindJSON validates and binds JSON request body to the provided struct
// Returns true if binding was successful, false if there was an error
func BindJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewDetailedErrorResponse("Invalid request data", err.Error(), "VALIDATION_ERROR"))
		return false
	}
	return true
}

// BindQuery validates and binds query parameters to the provided struct
// Returns true if binding was successful, false if there was an error
func BindQuery(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid query parameters"))
		return false
	}
	return true
}

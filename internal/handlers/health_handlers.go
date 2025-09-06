package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Ping
// @Description Ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Ping handler - not implemented yet",
		"data":    nil,
	})
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get users
// @Description Get users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users [get]
func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "GetUsers handler - not implemented yet",
		"data":    []interface{}{},
	})
}

// @Summary Create user
// @Description Create user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Router /users [post]
func CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "CreateUser handler - not implemented yet",
		"data":    nil,
	})
}

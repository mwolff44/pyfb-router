package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController defines a health struct
type HealthController struct{}

// Status handles the server status request
func (h HealthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "Working!")
}

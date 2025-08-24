package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck handler para verificar o status da API
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"message":   "API funcionando normalmente",
			"timestamp": "2024-01-01T00:00:00Z",
			"version":   "1.0.0",
		})
	}
}

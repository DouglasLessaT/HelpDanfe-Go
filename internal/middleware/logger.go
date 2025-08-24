package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger middleware para logging de requisições HTTP
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"client_ip":    param.ClientIP,
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
			"method":       param.Method,
			"path":         param.Path,
			"protocol":     param.Request.Proto,
			"status_code":  param.StatusCode,
			"latency":      param.Latency,
			"user_agent":   param.Request.UserAgent(),
			"error":        param.ErrorMessage,
		}).Info("HTTP Request")

		return ""
	})
}

// CORS middleware para configurar CORS de forma segura
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter origem permitida das variáveis de ambiente
		allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "*" // Fallback para desenvolvimento
		}

		// Verificar se a origem da requisição é permitida
		origin := c.Request.Header.Get("Origin")
		if allowedOrigin != "*" && origin != "" {
			allowedOrigins := strings.Split(allowedOrigin, ",")
			originAllowed := false
			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(allowed) == origin {
					originAllowed = true
					break
				}
			}
			if !originAllowed {
				c.AbortWithStatus(403)
				return
			}
		}

		// Configurar headers CORS
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
		}
		
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Header("Access-Control-Max-Age", "86400") // 24 horas

		// Responder a requisições OPTIONS (preflight)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Auth middleware para autenticação (placeholder)
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementar autenticação aqui
		// Por exemplo, verificar JWT token
		c.Next()
	}
}

// RateLimit middleware para limitar taxa de requisições (placeholder)
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementar rate limiting aqui
		c.Next()
	}
}

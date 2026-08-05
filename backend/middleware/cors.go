package middleware

import (
	"net/http"
	"strings"

	"github.com/bryantaolong/user-system/config"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := "http://localhost:5173"
		if config.AppConfig != nil && config.AppConfig.CORS.AllowedOrigins != "" {
			allowedOrigins = config.AppConfig.CORS.AllowedOrigins
		}

		origin := c.GetHeader("Origin")
		origins := strings.Split(allowedOrigins, ",")
		allowed := false
		for _, o := range origins {
			if strings.TrimSpace(o) == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "3600")
		}

		if c.Request.Method == "OPTIONS" {
			if allowed {
				c.AbortWithStatus(204)
			} else {
				c.AbortWithStatus(http.StatusNoContent)
			}
			return
		}

		c.Next()
	}
}

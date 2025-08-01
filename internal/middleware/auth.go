package middleware

import (
	"net/http"
	"strings"

	"github.com/bryantaolong/system/internal/service"
	"github.com/gin-gonic/gin"
)

func AuthRequired(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供Token"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if !authService.ValidateToken(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Token无效"})
			return
		}
		username, _ := authService.GetCurrentUsername(token)
		c.Set("username", username)
		c.Next()
	}
}

func RoleRequired(authService *service.AuthService, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		roles, _ := authService.GetCurrentUserRoles(token)
		if !strings.Contains(roles, role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足"})
			return
		}
		c.Next()
	}
}

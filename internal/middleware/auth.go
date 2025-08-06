// Package middleware 提供基于 JWT + Redis 的统一认证与鉴权中间件。
package middleware

import (
	"context"
	"net/http"

	"github.com/bryantaolong/system/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// AuthRequired 验证请求头中的 JWT 并与 Redis 中的 token 做一致性校验。
func AuthRequired(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := jwt.GetTokenFromRequest(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
			return
		}
		if !jwt.ValidateToken(tokenStr) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Token无效"})
			return
		}

		claims, _ := jwt.ParseToken(tokenStr)
		redisToken, _ := rdb.Get(context.Background(), claims.Username).Result()
		if redisToken != tokenStr {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Token已失效"})
			return
		}

		// 可选：刷新 TTL
		_ = rdb.Expire(context.Background(), claims.Username, jwt.Expiration)

		c.Set(jwt.ContextKey, claims)
		c.Next()
	}
}

// RoleRequired 要求当前用户必须具备指定角色。
func RoleRequired(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(jwt.ContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供Token"})
			return
		}
		cc := claims.(*jwt.CustomClaims)
		for _, r := range cc.Roles {
			if r == requiredRole || r == jwt.RolePrefix+requiredRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足"})
	}
}

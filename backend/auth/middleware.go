package auth

import (
	"net/http"
	"strings"

	"github.com/bryan/user-system/pkg/jwt"
	"github.com/bryan/user-system/cache"
	"github.com/bryan/user-system/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(userRepo *UserRepository, redisSvc *redis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenStr := jwt.ExtractToken(authHeader)

		if tokenStr == "" {
			// 无 Token，放行（由具体路由权限控制）
			c.Next()
			return
		}

		// 解析 Token
		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			writeUnauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}

		// Redis 白名单验证
		redisToken := redisSvc.Get(c.Request.Context(), claims.Username)
		if redisToken == "" || redisToken != tokenStr {
			writeUnauthorized(c, "Token已失效，请重新登录")
			c.Abort()
			return
		}

		// 从数据库加载用户，验证用户状态
		user, err := userRepo.SelectByID(claims.UserID)
		if err != nil || user == nil || !user.IsEnabled() || !user.IsAccountNonLocked() {
			writeUnauthorized(c, "用户状态异常或不存在")
			c.Abort()
			return
		}

		// 将 claims 和 user 存入 context
		c.Set("claims", claims)
		c.Set("user", user)

		// 设置角色到 context
		roles := user.GetAuthorities()
		c.Set("roles", roles)

		c.Next()
	}
}

// RequireAuth 要求认证中间件
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists || claims == nil {
			writeUnauthorized(c, "未授权，请先登录")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireRole 要求特定角色中间件
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			writeForbidden(c, "权限不足，无法访问此资源")
			c.Abort()
			return
		}

		roleList, ok := roles.([]string)
		if !ok {
			writeForbidden(c, "权限不足，无法访问此资源")
			c.Abort()
			return
		}

		targetRole := role
		if !strings.HasPrefix(targetRole, "ROLE_") {
			targetRole = "ROLE_" + targetRole
		}

		found := false
		for _, r := range roleList {
			if r == targetRole {
				found = true
				break
			}
		}

		if !found {
			writeForbidden(c, "权限不足，无法访问此资源")
			c.Abort()
			return
		}

		c.Next()
	}
}

func writeUnauthorized(c *gin.Context, msg string) {
	response.ErrorWithHTTPStatus(c, http.StatusUnauthorized, response.StatusUnauthorized, msg)
}

func writeForbidden(c *gin.Context, msg string) {
	response.ErrorWithHTTPStatus(c, http.StatusForbidden, response.StatusForbidden, msg)
}

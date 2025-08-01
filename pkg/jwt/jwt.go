package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	SecretKey   = "BryanTaoLong2025!@#SuperSecretKeyJwtToken987" // 生产环境应从配置读取
	Expiration  = 24 * time.Hour                                 // 24小时
	TokenPrefix = "Bearer "                                      // Token前缀
	ContextKey  = "JWT_CLAIMS"                                   // Gin上下文中存储claims的key
	UserIdKey   = "sub"                                          // 用户ID在claims中的key
	UsernameKey = "username"                                     // 用户名在claims中的key
	RolesKey    = "roles"                                        // 角色在claims中的key
	RolePrefix  = "ROLE_"                                        // 角色前缀
)

// CustomClaims 自定义Claims结构
type CustomClaims struct {
	UserId   string   `json:"sub"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userId, username string, roles []string) (string, error) {
	claims := CustomClaims{
		UserId:   userId,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenFromRequest 从请求中获取Token
func GetTokenFromRequest(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, TokenPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[len(TokenPrefix):], nil
}

// GetCurrentUserId 获取当前用户ID
func GetCurrentUserId(c *gin.Context) (string, error) {
	claims, exists := c.Get(ContextKey)
	if !exists {
		return "", errors.New("claims not found in context")
	}

	customClaims, ok := claims.(*CustomClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}

	return customClaims.UserId, nil
}

// GetCurrentUsername 获取当前用户名
func GetCurrentUsername(c *gin.Context) (string, error) {
	claims, exists := c.Get(ContextKey)
	if !exists {
		return "", errors.New("claims not found in context")
	}

	customClaims, ok := claims.(*CustomClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}

	return customClaims.Username, nil
}

// GetCurrentUserRoles 获取当前用户角色
func GetCurrentUserRoles(c *gin.Context) ([]string, error) {
	claims, exists := c.Get(ContextKey)
	if !exists {
		return nil, errors.New("claims not found in context")
	}

	customClaims, ok := claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	return customClaims.Roles, nil
}

// ValidateToken 验证Token有效性
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}

// JwtMiddleware JWT中间件
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := GetTokenFromRequest(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 将claims存入上下文
		c.Set(ContextKey, claims)
		c.Next()
	}
}

// HasRole 检查用户是否拥有指定角色
func HasRole(c *gin.Context, requiredRole string) bool {
	roles, err := GetCurrentUserRoles(c)
	if err != nil {
		return false
	}

	for _, role := range roles {
		if role == requiredRole || role == RolePrefix+requiredRole {
			return true
		}
	}
	return false
}

// MustHaveRole 中间件：要求用户必须拥有指定角色
func MustHaveRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasRole(c, requiredRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

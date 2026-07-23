package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bryan/user-system/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Roles    string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID int64, username, roles string) (string, error) {
	cfg := config.AppConfig.JWT
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.ExpirationMs) * time.Millisecond)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWT.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ValidateToken 验证 Token 是否有效
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}

// GetUsernameFromToken 从 Token 中获取用户名
func GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

// GetRolesFromToken 从 Token 中获取角色列表
func GetRolesFromToken(tokenString string) []string {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil
	}
	if claims.Roles == "" {
		return nil
	}
	var roles []string
	for _, r := range strings.Split(claims.Roles, ",") {
		r = strings.TrimSpace(r)
		if r != "" {
			if !strings.HasPrefix(r, "ROLE_") {
				r = "ROLE_" + r
			}
			roles = append(roles, r)
		}
	}
	return roles
}

// GetUserIDFromToken 从 Token 中获取用户 ID
func GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ExtractToken 从 Authorization 头中提取 Token
func ExtractToken(authHeader string) string {
	cfg := config.AppConfig.JWT
	if authHeader == "" || !strings.HasPrefix(authHeader, cfg.TokenPrefix) {
		return ""
	}
	return strings.TrimPrefix(authHeader, cfg.TokenPrefix)
}

// GetRolesFromTokenByClaims 从 Claims 中获取角色列表
func GetRolesFromTokenByClaims(claims *Claims) []string {
	if claims == nil || claims.Roles == "" {
		return nil
	}
	var roles []string
	for _, r := range strings.Split(claims.Roles, ",") {
		r = strings.TrimSpace(r)
		if r != "" {
			if !strings.HasPrefix(r, "ROLE_") {
				r = "ROLE_" + r
			}
			roles = append(roles, r)
		}
	}
	return roles
}

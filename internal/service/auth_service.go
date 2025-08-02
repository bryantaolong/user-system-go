// Package service 提供用户认证与授权的核心业务逻辑。
// 所有与 JWT、Redis、数据库交互的操作均通过 pkg/jwt 和 go-redis 原生客户端完成。
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	http2 "github.com/bryantaolong/system/pkg/http"
	"github.com/bryantaolong/system/pkg/jwt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strings"
	"time"

	"github.com/bryantaolong/system/internal/model/entity"
	"github.com/bryantaolong/system/internal/model/request"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 负责用户注册、登录、登出、Token 刷新、权限判断等业务。
type AuthService struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewAuthService 创建并返回一个 AuthService 实例。
func NewAuthService(db *gorm.DB, rdb *redis.Client) *AuthService {
	return &AuthService{
		db:    db,
		redis: rdb,
	}
}

// Register 用户注册：检查用户名唯一性、密码加密、写入数据库。
func (s *AuthService) Register(registerReq request.RegisterRequest) (*entity.User, error) {
	var existingUser entity.User
	if err := s.db.Where("username = ?", registerReq.Username).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	user := entity.User{
		Username:          registerReq.Username,
		Password:          string(hashedPassword),
		PhoneNumber:       registerReq.PhoneNumber,
		Email:             registerReq.Email,
		Roles:             "ROLE_USER",
		PasswordResetTime: sql.NullTime{Time: time.Now(), Valid: true},
		CreateBy:          registerReq.Username,
		UpdateBy:          registerReq.Username,
		CreateTime:        time.Now(),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	return &user, nil
}

// Login 用户登录：验证密码、生成/复用 JWT、写 Redis、记录登录信息。
func (s *AuthService) Login(loginReq request.LoginRequest, r *http.Request) (string, error) {
	var user entity.User
	if err := s.db.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("用户名或密码错误")
		}
		return "", fmt.Errorf("查询用户失败: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		user.LoginFailCount++
		if user.LoginFailCount >= 5 {
			user.Status = 2
			user.AccountLockTime = sql.NullTime{Time: time.Now(), Valid: true}
			s.db.Save(&user)
			return "", fmt.Errorf("输入密码错误次数过多，账号锁定")
		}
		s.db.Save(&user)
		return "", fmt.Errorf("用户名或密码错误")
	}

	if user.Status == 1 {
		return "", fmt.Errorf("账号已被封禁")
	}

	if user.Status == 2 && user.AccountLockTime.Valid {
		if time.Since(user.AccountLockTime.Time) < time.Hour {
			return "", fmt.Errorf("账号已被锁定，请稍后再试")
		}
		user.Status = 0
		user.AccountLockTime = sql.NullTime{Valid: false}
	}

	// 1. 先看 Redis 是否有未过期 token
	existing, _ := s.redis.Get(context.Background(), user.Username).Result()
	if existing != "" && jwt.ValidateToken(existing) {
		_ = s.redis.Expire(context.Background(), user.Username, jwt.Expiration)
		return existing, nil
	}

	// 2. 生成新 token
	token, err := jwt.GenerateToken(fmt.Sprint(user.ID), user.Username, []string{user.Roles})
	if err != nil {
		return "", fmt.Errorf("生成Token失败: %v", err)
	}

	// 3. 存 Redis
	if err := s.redis.Set(context.Background(), user.Username, token, jwt.Expiration).Err(); err != nil {
		return "", fmt.Errorf("Token存储失败: %v", err)
	}

	user.LoginTime = sql.NullTime{Time: time.Now(), Valid: true}
	user.LoginIP = http2.GetClientIP(r)
	user.LoginFailCount = 0
	user.UpdateTime = sql.NullTime{Time: time.Now(), Valid: true}
	user.UpdateBy = user.Username

	if err := s.db.Save(&user).Error; err != nil {
		return "", fmt.Errorf("更新用户信息失败: %v", err)
	}

	token, err = jwt.GenerateToken(user.Username, user.Roles, []string{user.Roles})
	if err != nil {
		return "", fmt.Errorf("生成Token失败: %v", err)
	}

	if err := s.redis.Set(context.Background(), user.Username, token, 24*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("Token存储失败: %v", err)
	}

	return token, nil
}

// GetCurrentUser 根据 token 获取完整用户信息。
func (s *AuthService) GetCurrentUser(tokenString string) (*entity.User, error) {
	userID, err := s.GetCurrentUserID(tokenString)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return &user, nil
}

// GetCurrentUserID 从 token 中提取用户 ID。
func (s *AuthService) GetCurrentUserID(tokenString string) (string, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserId, nil
}

// GetCurrentUsername 从 token 中提取用户名。
func (s *AuthService) GetCurrentUsername(tokenString string) (string, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

// GetCurrentUserRoles 从 token 中提取角色字符串。
func (s *AuthService) GetCurrentUserRoles(tokenString string) (string, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return strings.Join(claims.Roles, ","), nil
}

// IsAdmin 判断用户是否为管理员。
func (s *AuthService) IsAdmin(tokenString string) (bool, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return false, err
	}
	for _, r := range claims.Roles {
		if r == "ROLE_ADMIN" {
			return true, nil
		}
	}
	return false, nil
}

// ValidateToken 验证 token 是否有效。
func (s *AuthService) ValidateToken(tokenString string) bool {
	return jwt.ValidateToken(tokenString)
}

// RefreshToken 刷新 token。
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return jwt.GenerateToken(claims.UserId, claims.Username, claims.Roles)
}

// Logout 删除 Redis 中的 token 实现登出。
func (s *AuthService) Logout(tokenString string) error {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return err
	}
	_, err = s.redis.Del(context.Background(), claims.Username).Result()
	return err
}

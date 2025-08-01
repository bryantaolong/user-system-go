package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	http2 "github.com/bryantaolong/system/pkg/http"
	"net/http"
	"strings"
	"time"

	"github.com/bryantaolong/system/internal/model/entity"
	"github.com/bryantaolong/system/internal/model/request"
	"github.com/bryantaolong/system/internal/service/redis"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db           *gorm.DB
	redisService *redis.RedisStringService
	jwtSecret    string
}

func NewAuthService(db *gorm.DB, redisService *redis.RedisStringService, jwtSecret string) *AuthService {
	return &AuthService{
		db:           db,
		redisService: redisService,
		jwtSecret:    jwtSecret,
	}
}

// Register 用户注册
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

// Login 用户登录
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

	existingToken, err := s.redisService.Get(context.Background(), user.Username)
	if err == nil && existingToken != "" {
		if _, err := s.validateToken(existingToken); err == nil {
			s.redisService.SetExpire(context.Background(), user.Username, 24*time.Hour)
			return existingToken, nil
		}
	}

	user.LoginTime = sql.NullTime{Time: time.Now(), Valid: true}
	user.LoginIP = http2.GetClientIP(r)
	user.LoginFailCount = 0
	user.UpdateTime = sql.NullTime{Time: time.Now(), Valid: true}
	user.UpdateBy = user.Username

	if err := s.db.Save(&user).Error; err != nil {
		return "", fmt.Errorf("更新用户信息失败: %v", err)
	}

	token, err := s.generateToken(&user)
	if err != nil {
		return "", fmt.Errorf("生成Token失败: %v", err)
	}

	if err := s.redisService.SetWithExpire(context.Background(), user.Username, token, 24*time.Hour); err != nil {
		return "", fmt.Errorf("Token存储失败: %v", err)
	}

	return token, nil
}

// GetCurrentUserID 获取当前用户ID
func (s *AuthService) GetCurrentUserID(tokenString string) (uint, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return 0, err
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("无效的Token声明")
	}
	return uint(userID), nil
}

// GetCurrentUsername 获取当前用户名
func (s *AuthService) GetCurrentUsername(tokenString string) (string, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return "", err
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("无效的Token声明")
	}
	return username, nil
}

// GetCurrentUser 获取当前用户完整信息
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

// GetCurrentUserRoles 获取当前用户的角色字符串
func (s *AuthService) GetCurrentUserRoles(tokenString string) (string, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return "", err
	}
	roles, ok := claims["roles"].(string)
	if !ok {
		return "", fmt.Errorf("roles 字段无效")
	}
	return roles, nil
}

// IsAdmin 判断用户是否为管理员
func (s *AuthService) IsAdmin(tokenString string) (bool, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return false, err
	}

	roles, ok := claims["roles"].(string)
	if !ok {
		return false, nil
	}
	return strings.Contains(roles, "ROLE_ADMIN"), nil
}

// ValidateToken 校验Token是否有效
func (s *AuthService) ValidateToken(tokenString string) bool {
	_, err := s.validateToken(tokenString)
	return err == nil
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	user, err := s.GetCurrentUser(tokenString)
	if err != nil {
		return "", err
	}
	return s.generateToken(user)
}

// Logout 用户登出
func (s *AuthService) Logout(tokenString string) error {
	username, err := s.GetCurrentUsername(tokenString)
	if err != nil {
		return err
	}
	if _, err := s.redisService.Delete(context.Background(), username); err != nil {
		return fmt.Errorf("清除Token失败: %v", err)
	}
	return nil
}

// generateToken 生成JWT Token
func (s *AuthService) generateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"roles":    user.Roles,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// validateToken 验证Token
func (s *AuthService) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("无效的Token")
}

package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/bryantaolong/user-system/config"
	"github.com/bryantaolong/user-system/model"
	pkgHttp "github.com/bryantaolong/user-system/pkg/http"
	"github.com/bryantaolong/user-system/pkg/jwt"
	"github.com/bryantaolong/user-system/cache"
	"github.com/bryantaolong/user-system/response"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service 认证服务
type Service struct {
	userRepo *UserRepository
	roleRepo *UserRoleRepository
	redis    *redis.RedisClient
	logger   *zap.Logger
}

// NewService 创建认证服务实例
func NewService(
	userRepo *UserRepository,
	roleRepo *UserRoleRepository,
	redisSvc *redis.RedisClient,
	logger *zap.Logger,
) *Service {
	return &Service{
		userRepo: userRepo,
		roleRepo: roleRepo,
		redis:    redisSvc,
		logger:   logger,
	}
}

// Register 用户注册
func (s *Service) Register(ctx context.Context, req *model.RegisterRequest) (*model.SysUser, error) {
	// 1. 检查用户名是否已存在
	existing, _ := s.userRepo.SelectByUsername(req.Username)
	if existing != nil {
		return nil, response.NewBusinessError("用户名已存在")
	}

	// 2. 查出默认角色
	defaultRole, err := s.roleRepo.SelectOneByIsDefaultTrue()
	if err != nil || defaultRole == nil {
		return nil, response.NewBusinessError("系统未配置默认角色")
	}

	// 3. 构建用户实体，密码加密
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBusinessError("密码加密失败")
	}

	now := time.Now()
	phone := req.Phone
	email := req.Email
	roles := defaultRole.RoleName

	user := &model.SysUser{
		Username:        req.Username,
		Password:        string(hashedPwd),
		Phone:           &phone,
		Email:           &email,
		Status:          model.UserStatusNormal,
		Roles:           &roles,
		PasswordResetAt: &now,
		LoginFailCount:  intPtr(0),
		Deleted:         0,
		Version:         0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// 4. 插入用户数据
	if err := s.userRepo.Insert(user); err != nil {
		return nil, response.NewBusinessError("插入数据库失败")
	}

	// 5. 回填审计字段（注册时无认证用户，操作人设为 "0"）
	operator := "0"
	user.CreatedBy = &operator
	user.UpdatedBy = &operator
	_ = s.userRepo.Update(user)

	s.logger.Info("用户注册成功", zap.Int64("id", user.ID), zap.String("username", user.Username))
	return user, nil
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *model.LoginRequest, userAgent, clientIP string) (string, error) {
	// 1. 验证用户凭证
	user, err := s.userRepo.SelectByUsername(req.Username)
	if err != nil || user == nil {
		s.logger.Warn("登录失败 - 用户不存在", zap.String("username", req.Username))
		return "", response.NewBusinessError("用户名或密码错误")
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		currentFailCount := 0
		if user.LoginFailCount != nil {
			currentFailCount = *user.LoginFailCount
		}
		user.LoginFailCount = intPtr(currentFailCount + 1)
		operator := fmt.Sprintf("%d", user.ID)
		if operator == "0" {
			operator = "0"
		}
		user.UpdatedBy = &operator
		user.Version++

		// 如果输入密码错误次数达到限额，则锁定账号
		if *user.LoginFailCount >= config.AppConfig.Security.LoginFailLimit {
			now := time.Now()
			user.Status = model.UserStatusLocked
			user.LockedAt = &now
			_ = s.userRepo.Update(user)
			s.logger.Warn("用户登录失败次数过多，已锁定", zap.String("username", user.Username))
			return "", response.NewBusinessError("输入密码错误次数过多，账号锁定")
		}
		_ = s.userRepo.Update(user)
		s.logger.Warn("用户登录密码错误", zap.String("username", user.Username), zap.Int("failCount", *user.LoginFailCount))
		return "", response.NewBusinessError("用户名或密码错误")
	}

	// 3. 检查现有 Token
	existingToken := s.redis.Get(ctx, user.Username)
	if existingToken != "" && jwt.ValidateToken(existingToken) {
		// 刷新 Redis 中的 Token 过期时间
		s.redis.SetExpire(ctx, user.Username, config.AppConfig.JWT.ExpirationMs/1000)
		return existingToken, nil
	}

	// 4. 更新用户登录信息
	now := time.Now()
	device := pkgHttp.GetClientOS(userAgent) + " / " + pkgHttp.GetClientAgent(userAgent)
	user.LastLoginAt = &now
	user.LastLoginIP = &clientIP
	user.LastLoginDevice = &device
	user.LoginFailCount = intPtr(0)
	operator := fmt.Sprintf("%d", user.ID)
	user.Version++
	user.UpdatedAt = now
	user.UpdatedBy = &operator
	_ = s.userRepo.Update(user)

	// 5. 生成新的 JWT Token
	roles := ""
	if user.Roles != nil {
		roles = *user.Roles
	}
	token, err := jwt.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		return "", response.NewBusinessError("Token 生成失败")
	}

	// 6. 存储到 Redis
	if !s.redis.SetWithExpire(ctx, user.Username, token, config.AppConfig.JWT.ExpirationMs/1000) {
		return "", response.NewBusinessError("Token 存储失败")
	}

	return token, nil
}

// GetCurrentUserID 获取当前登录用户的 ID
func (s *Service) GetCurrentUserID(ctx context.Context) (int64, error) {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return 0, response.NewUnauthorizedError("未授权")
	}
	return claims.UserID, nil
}

// GetCurrentUsername 获取当前登录用户的用户名
func (s *Service) GetCurrentUsername(ctx context.Context) (string, error) {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return "", response.NewUnauthorizedError("未授权")
	}
	return claims.Username, nil
}

// GetCurrentUser 获取当前登录用户的完整信息
func (s *Service) GetCurrentUser(ctx context.Context) (*model.SysUser, error) {
	userID, err := s.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.SelectByID(userID)
	if err != nil {
		return nil, response.NewResourceNotFoundError("用户不存在")
	}
	return user, nil
}

// IsAdmin 判断用户是否具有管理员权限
func (s *Service) IsAdmin(ctx context.Context) bool {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return false
	}
	for _, role := range jwt.GetRolesFromTokenByClaims(claims) {
		if role == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}

// ValidateToken 校验 JWT Token 是否有效
func (s *Service) ValidateToken(token string) bool {
	return jwt.ValidateToken(token)
}

// ValidateTokenWithStatus 校验 JWT Token 及账户状态
func (s *Service) ValidateTokenWithStatus(token string) (string, error) {
	if !jwt.ValidateToken(token) {
		return "Invalid token", nil
	}
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return "Invalid token", nil
	}

	user, err := s.userRepo.SelectByID(claims.UserID)
	if err != nil || user == nil {
		return "Invalid token", nil
	}

	if !user.IsEnabled() {
		return "Account disabled", nil
	}
	if !user.IsAccountNonLocked() {
		return "Account locked", nil
	}

	return "Validation passed", nil
}

// RefreshToken 刷新 JWT Token
func (s *Service) RefreshToken(ctx context.Context) (string, error) {
	username, err := s.GetCurrentUsername(ctx)
	if err != nil {
		return "", err
	}

	existingToken := s.redis.Get(ctx, username)
	if existingToken == "" || !jwt.ValidateToken(existingToken) {
		return "", response.NewUnauthorizedError("Token 无效或已过期")
	}

	claims, _ := jwt.ParseToken(existingToken)
	newToken, err := jwt.GenerateToken(claims.UserID, claims.Username, claims.Roles)
	if err != nil {
		return "", response.NewBusinessError("Token 生成失败")
	}

	if !s.redis.SetWithExpire(ctx, username, newToken, config.AppConfig.JWT.ExpirationMs/1000) {
		return "", response.NewBusinessError("Token 存储失败")
	}

	return newToken, nil
}

// ChangePassword 修改用户密码
func (s *Service) ChangePassword(ctx context.Context, oldPassword, newPassword string) (*model.SysUser, error) {
	user, err := s.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return nil, response.NewBusinessError("旧密码不正确")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBusinessError("密码加密失败")
	}

	now := time.Now()
	user.Password = string(hashedPwd)
	user.PasswordResetAt = &now
	operator := fmt.Sprintf("%d", user.ID)

	user.Version++
	user.UpdatedAt = now
	user.UpdatedBy = &operator

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewBusinessError("密码修改失败")
	}

	// 清除 Redis 中的旧 Token
	if !s.redis.Delete(ctx, user.Username) {
		s.logger.Warn("用户密码更新成功，但清除旧 Token 失败", zap.Int64("userId", user.ID))
	} else {
		s.logger.Info("用户密码更新成功，旧 Token 已清除", zap.Int64("userId", user.ID))
	}

	return user, nil
}

// Logout 退出登录
func (s *Service) Logout(ctx context.Context) error {
	username, err := s.GetCurrentUsername(ctx)
	if err != nil {
		return err
	}
	if !s.redis.Delete(ctx, username) {
		return response.NewBusinessError("Token 清除失败")
	}
	return nil
}

// DeleteAccount 注销用户
func (s *Service) DeleteAccount(ctx context.Context) (*model.SysUser, error) {
	user, err := s.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, response.NewResourceNotFoundError("用户状态异常，无法注销")
	}
	if err := s.userRepo.DeleteByID(user.ID, fmt.Sprintf("%d", user.ID)); err != nil {
		return nil, response.NewBusinessError("注销失败")
	}

	// 清除 Redis 中的 Token
	if !s.redis.Delete(ctx, user.Username) {
		s.logger.Warn("用户注销成功，但清除 Token 失败", zap.Int64("userId", user.ID))
	} else {
		s.logger.Info("用户注销成功，Token 已清除", zap.Int64("userId", user.ID))
	}
	return user, nil
}

// helper functions
func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }

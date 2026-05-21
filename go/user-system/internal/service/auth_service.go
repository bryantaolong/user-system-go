package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bryan/user-system/internal/config"
	"github.com/bryan/user-system/internal/model"
	pkgHttp "github.com/bryan/user-system/internal/pkg/http"
	"github.com/bryan/user-system/internal/pkg/jwt"
	"github.com/bryan/user-system/internal/pkg/redis"
	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.UserRoleRepository
	redisSvc *redis.RedisService
	logger   *zap.Logger
}

func NewAuthService(
	userRepo *repository.UserRepository,
	roleRepo *repository.UserRoleRepository,
	redisSvc *redis.RedisService,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		redisSvc: redisSvc,
		logger:   logger,
	}
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.SysUser, error) {
	// 1. 检查用户名是否已存在
	existing, _ := s.userRepo.SelectByUsername(req.Username)
	if existing != nil {
		return nil, response.NewBusinessException("用户名已存在")
	}

	// 2. 查出默认角色
	defaultRole, err := s.roleRepo.SelectOneByIsDefaultTrue()
	if err != nil || defaultRole == nil {
		return nil, response.NewBusinessException("系统未配置默认角色")
	}

	// 3. 构建用户实体，密码加密
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBusinessException("密码加密失败")
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
		Roles:           &roles,
		Status:          model.UserStatusNormal,
		Deleted:         0,
		Version:         0,
		PasswordResetAt: &now,
		LoginFailCount:  intPtr(0),
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       strPtr(req.Username),
		UpdatedBy:       strPtr(req.Username),
	}

	// 4. 插入用户数据
	if err := s.userRepo.Insert(user); err != nil {
		return nil, response.NewBusinessException("插入数据库失败")
	}

	s.logger.Info("用户注册成功", zap.Int64("id", user.ID), zap.String("username", user.Username))
	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest, userAgent, clientIP string) (string, error) {
	// 1. 验证用户凭证
	sysUser, err := s.userRepo.SelectByUsername(req.Username)
	if err != nil || sysUser == nil {
		s.logger.Warn("登录失败 - 用户不存在", zap.String("username", req.Username))
		return "", response.NewBusinessException("用户名或密码错误")
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(sysUser.Password), []byte(req.Password)); err != nil {
		currentFailCount := 0
		if sysUser.LoginFailCount != nil {
			currentFailCount = *sysUser.LoginFailCount
		}
		sysUser.LoginFailCount = intPtr(currentFailCount + 1)
		sysUser.UpdatedBy = strPtr(fmt.Sprintf("%d", sysUser.ID))
		sysUser.Version++

		// 如果输入密码错误次数达到限额，则锁定账号
		if *sysUser.LoginFailCount >= config.AppConfig.Security.LoginFailLimit {
			now := time.Now()
			sysUser.Status = model.UserStatusLocked
			sysUser.LockedAt = &now
			_ = s.userRepo.Update(sysUser)
			s.logger.Warn("用户登录失败次数过多，已锁定", zap.String("username", sysUser.Username))
			return "", response.NewBusinessException("输入密码错误次数过多，账号锁定")
		}
		_ = s.userRepo.Update(sysUser)
		s.logger.Warn("用户登录密码错误", zap.String("username", sysUser.Username), zap.Int("failCount", *sysUser.LoginFailCount))
		return "", response.NewBusinessException("用户名或密码错误")
	}

	// 3. 检查现有 Token
	existingToken := s.redisSvc.Get(ctx, sysUser.Username)
	if existingToken != "" && jwt.ValidateToken(existingToken) {
		// 刷新 Redis 中的 Token 过期时间
		s.redisSvc.SetExpire(ctx, sysUser.Username, config.AppConfig.JWT.ExpirationMs/1000)
		return existingToken, nil
	}

	// 4. 更新用户登录信息
	now := time.Now()
	device := pkgHttp.GetClientOS(userAgent) + " / " + pkgHttp.GetClientAgent(userAgent)
	sysUser.LastLoginAt = &now
	sysUser.LastLoginIP = &clientIP
	sysUser.LastLoginDevice = &device
	sysUser.LoginFailCount = intPtr(0)
	sysUser.UpdatedBy = strPtr(fmt.Sprintf("%d", sysUser.ID))
	sysUser.Version++
	_ = s.userRepo.Update(sysUser)

	// 5. 生成新的 JWT Token
	roles := ""
	if sysUser.Roles != nil {
		roles = *sysUser.Roles
	}
	token, err := jwt.GenerateToken(sysUser.ID, sysUser.Username, roles)
	if err != nil {
		return "", response.NewBusinessException("Token 生成失败")
	}

	// 6. 存储到 Redis
	if !s.redisSvc.SetWithExpire(ctx, sysUser.Username, token, config.AppConfig.JWT.ExpirationMs/1000) {
		return "", response.NewBusinessException("Token 存储失败")
	}

	return token, nil
}

// GetCurrentUserID 获取当前登录用户的 ID
func (s *AuthService) GetCurrentUserID(ctx context.Context) (int64, error) {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return 0, response.NewUnauthorizedException("未授权")
	}
	return claims.UserID, nil
}

// GetCurrentUsername 获取当前登录用户的用户名
func (s *AuthService) GetCurrentUsername(ctx context.Context) (string, error) {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return "", response.NewUnauthorizedException("未授权")
	}
	return claims.Username, nil
}

// GetCurrentUser 获取当前登录用户的完整信息
func (s *AuthService) GetCurrentUser(ctx context.Context) (*model.SysUser, error) {
	userID, err := s.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.SelectByID(userID)
	if err != nil {
		return nil, response.NewResourceNotFoundException("用户不存在")
	}
	return user, nil
}

// IsAdmin 判断用户是否具有管理员权限
func (s *AuthService) IsAdmin(ctx context.Context) bool {
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
func (s *AuthService) ValidateToken(token string) bool {
	return jwt.ValidateToken(token)
}

// ChangePassword 修改用户密码
func (s *AuthService) ChangePassword(ctx context.Context, oldPassword, newPassword string) (*model.SysUser, error) {
	user, err := s.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return nil, response.NewBusinessException("旧密码不正确")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBusinessException("密码加密失败")
	}

	now := time.Now()
	user.Password = string(hashedPwd)
	user.PasswordResetAt = &now
	user.Version++
	user.UpdatedBy = strPtr(fmt.Sprintf("%d", user.ID))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewBusinessException("密码修改失败")
	}

	// 清除 Redis 中的旧 Token
	if !s.redisSvc.Delete(ctx, user.Username) {
		s.logger.Warn("用户密码更新成功，但清除旧 Token 失败", zap.Int64("userId", user.ID))
	} else {
		s.logger.Info("用户密码更新成功，旧 Token 已清除", zap.Int64("userId", user.ID))
	}

	return user, nil
}

// Logout 退出登录
func (s *AuthService) Logout(ctx context.Context) error {
	username, err := s.GetCurrentUsername(ctx)
	if err != nil {
		return err
	}
	if !s.redisSvc.Delete(ctx, username) {
		return response.NewBusinessException("Token 清除失败")
	}
	return nil
}

// DeleteAccount 注销用户
func (s *AuthService) DeleteAccount(ctx context.Context) (*model.SysUser, error) {
	user, err := s.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, response.NewResourceNotFoundException("用户状态异常，无法注销")
	}
	if err := s.userRepo.DeleteByID(user.ID, fmt.Sprintf("%d", user.ID)); err != nil {
		return nil, response.NewBusinessException("注销失败")
	}
	s.logger.Info("用户注销成功", zap.Int64("userId", user.ID))
	return user, nil
}

// helper functions
func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }

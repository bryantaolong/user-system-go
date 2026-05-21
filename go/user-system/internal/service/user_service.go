package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/pkg/jwt"
	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    *repository.UserRepository
	roleSvc     *UserRoleService
	logger      *zap.Logger
}

func NewUserService(
	userRepo *repository.UserRepository,
	roleSvc *UserRoleService,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleSvc:  roleSvc,
		logger:   logger,
	}
}

// CreateUser 管理员创建用户
func (s *UserService) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.SysUser, error) {
	existing, _ := s.userRepo.SelectByUsername(req.Username)
	if existing != nil {
		return nil, response.NewBusinessException("用户名已存在")
	}

	roleIDs := req.RoleIDs
	if len(roleIDs) == 0 {
		defaultRole, err := s.roleSvc.GetDefaultRole(ctx)
		if err != nil || defaultRole == nil {
			return nil, response.NewBusinessException("系统未配置默认角色")
		}
		roleIDs = []int{defaultRole.ID}
	}

	roles, err := s.roleSvc.ListByIDs(ctx, roleIDs)
	if err != nil {
		return nil, response.NewBusinessException("角色查询失败")
	}
	if len(roles) != len(roleIDs) {
		return nil, response.NewBusinessException("部分角色不存在")
	}

	roleNames := make([]string, 0, len(roles))
	for _, r := range roles {
		roleNames = append(roleNames, r.RoleName)
	}
	rolesStr := strings.Join(roleNames, ",")

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBusinessException("密码加密失败")
	}

	operator := getCurrentOperator(ctx)
	user := &model.SysUser{
		Username:        req.Username,
		Password:        string(hashedPwd),
		Phone:          req.Phone,
		Email:          req.Email,
		Roles:          &rolesStr,
		Status:         model.UserStatusNormal,
		Deleted:        0,
		Version:        0,
		LoginFailCount: intPtr(0),
		CreatedBy:      &operator,
		UpdatedBy:      &operator,
	}

	if err := s.userRepo.Insert(user); err != nil {
		return nil, response.NewBusinessException("插入数据库失败")
	}

	s.logger.Info("管理员创建用户成功", zap.Int64("id", user.ID), zap.String("username", user.Username))
	return user, nil
}

// GetAllUsers 获取所有用户列表（分页）
func (s *UserService) GetAllUsers(ctx context.Context, pageNum, pageSize int) (*response.PageResult, error) {
	offset := (pageNum - 1) * pageSize
	users, err := s.userRepo.SelectPage(offset, pageSize, nil)
	if err != nil {
		return nil, response.NewBusinessException("查询用户列表失败")
	}
	total, _ := s.userRepo.Count(nil)

	rows := make([]interface{}, 0, len(users))
	for i := range users {
		rows = append(rows, users[i])
	}
	pr := response.NewPageResult(rows, total, pageNum, pageSize)
	return &pr, nil
}

// GetUserByID 根据用户 ID 获取用户信息
func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*model.SysUser, error) {
	user, err := s.userRepo.SelectByID(userID)
	if err != nil {
		return nil, response.NewResourceNotFoundException("用户不存在")
	}
	return user, nil
}

// GetUserByUsername 根据用户名获取用户信息
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	user, err := s.userRepo.SelectByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// QueryUsers 通用用户搜索
func (s *UserService) QueryUsers(ctx context.Context, req *model.UserSearchRequest, pageNum, pageSize int) (*response.PageResult, error) {
	offset := (pageNum - 1) * pageSize
	users, err := s.userRepo.SelectPage(offset, pageSize, req)
	if err != nil {
		return nil, response.NewBusinessException("查询用户失败")
	}
	total, _ := s.userRepo.Count(req)

	rows := make([]interface{}, 0, len(users))
	for i := range users {
		rows = append(rows, users[i])
	}
	pr := response.NewPageResult(rows, total, pageNum, pageSize)
	return &pr, nil
}

// UpdateUser 更新用户基础信息
func (s *UserService) UpdateUser(ctx context.Context, userID int64, req *model.UserUpdateRequest) (*model.SysUser, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Email != nil {
		user.Email = req.Email
	}

	user.Version++
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewBusinessException("更新失败，可能数据已变更")
	}

	s.logger.Info("用户信息更新成功", zap.Int64("userId", userID))
	return user, nil
}

// ChangeRoleByIds 修改用户角色
func (s *UserService) ChangeRoleByIds(ctx context.Context, userID int64, req *model.ChangeRoleRequest) (*model.SysUser, error) {
	roles, err := s.roleSvc.ListByIDs(ctx, req.RoleIDs)
	if err != nil {
		return nil, response.NewBusinessException("角色查询失败")
	}
	if len(roles) != len(req.RoleIDs) {
		return nil, response.NewBusinessException("部分角色不存在")
	}

	roleNames := make([]string, 0, len(roles))
	for _, r := range roles {
		roleNames = append(roleNames, r.RoleName)
	}
	rolesStr := strings.Join(roleNames, ",")

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Roles = &rolesStr
	user.Version++
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewBusinessException("角色更新失败")
	}
	return user, nil
}

// ResetPassword 重置用户密码（管理员）
func (s *UserService) ResetPassword(ctx context.Context, userID int64, newPassword string) (*model.SysUser, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBusinessException("密码加密失败")
	}

	user.Password = string(hashedPwd)
	user.Version++
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewBusinessException("密码重置失败")
	}

	s.logger.Info("用户密码强制修改成功", zap.Int64("userId", userID))
	return user, nil
}

// BlockUser 封禁指定用户
func (s *UserService) BlockUser(ctx context.Context, userID int64) (*model.SysUser, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Status = model.UserStatusBanned
	user.Version++
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewBusinessException("封禁失败")
	}

	s.logger.Info("用户封禁成功", zap.Int64("userId", userID))
	return user, nil
}

// UnblockUser 解封指定用户
func (s *UserService) UnblockUser(ctx context.Context, userID int64) (*model.SysUser, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Status = model.UserStatusNormal
	user.Version++
	user.UpdatedBy = strPtr(getCurrentOperator(ctx))

	if err := s.userRepo.Update(user); err != nil {
		return nil, response.NewBusinessException("解封失败")
	}

	s.logger.Info("用户解封成功", zap.Int64("userId", userID))
	return user, nil
}

// DeleteUser 删除用户（逻辑删除）
func (s *UserService) DeleteUser(ctx context.Context, userID int64) (int64, error) {
	operator := getCurrentOperator(ctx)
	if err := s.userRepo.DeleteByID(userID, operator); err != nil {
		return 0, response.NewResourceNotFoundException("用户不存在或已被删除")
	}
	s.logger.Info("用户删除成功 (逻辑删除)", zap.Int64("userId", userID))
	return userID, nil
}

// getCurrentOperator 获取当前操作者标识
func getCurrentOperator(ctx context.Context) string {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return "SYSTEM"
	}
	return fmt.Sprintf("%d", claims.UserID)
}

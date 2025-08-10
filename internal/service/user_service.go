package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	_ "strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bryantaolong/system/internal/model/entity"
	"github.com/bryantaolong/system/internal/model/request"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db          *gorm.DB
	authService *AuthService
}

func NewUserService(db *gorm.DB, authService *AuthService) *UserService {
	return &UserService{db: db, authService: authService}
}

// GetAllUsers 获取所有用户（分页）
func (s *UserService) GetAllUsers(ctx context.Context, page request.PageRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	offset := page.GetOffset()
	if err := s.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := s.db.WithContext(ctx).Limit(int(page.PageSize)).Offset(int(offset)).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	var user entity.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// SearchUsers 支持多条件搜索
func (s *UserService) SearchUsers(ctx context.Context, req request.UserSearchRequest, page request.PageRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := s.db.WithContext(ctx).Model(&entity.User{})
	query = s.buildSearchQuery(query, req)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := page.GetOffset()
	if err := query.Limit(int(page.PageSize)).Offset(int(offset)).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *UserService) buildSearchQuery(query *gorm.DB, req request.UserSearchRequest) *gorm.DB {
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Roles != "" {
		query = query.Where("roles LIKE ?", "%"+req.Roles+"%")
	}
	if req.Status != nil && *req.Status >= 0 {
		query = query.Where("status = ?", *req.Status)
	}
	if req.LoginFailCount != nil && *req.LoginFailCount >= 0 {
		query = query.Where("login_fail_count = ?", *req.LoginFailCount)
	}
	if req.Deleted != nil && *req.Deleted >= 0 {
		query = query.Where("deleted = ?", *req.Deleted)
	}
	if !req.CreateTimeStart.IsZero() && !req.CreateTimeEnd.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", req.CreateTimeStart, req.CreateTimeEnd)
	}
	if !req.UpdateTimeStart.IsZero() && !req.UpdateTimeEnd.IsZero() {
		query = query.Where("updated_at BETWEEN ? AND ?", req.UpdateTimeStart, req.UpdateTimeEnd)
	}
	return query
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, userID int64, req request.UserUpdateRequest) (*entity.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Username != "" && req.Username != user.Username {
		existing, _ := s.GetUserByUsername(ctx, req.Username)
		if existing != nil && existing.ID != userID {
			return nil, fmt.Errorf("用户名已存在")
		}
		user.Username = req.Username
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	token := extractTokenFromContext(ctx)
	operator, _ := s.authService.GetCurrentUsername(token)
	user.UpdatedBy = operator
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// ChangeRoleByIds 根据角色 ID 列表批量修改用户角色
func (s *UserService) ChangeRoleByIds(ctx context.Context, userID int64, req request.ChangeRoleRequest) (*entity.User, error) {
	// 1. 查询用户
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 查询对应的 UserRole
	var roles []entity.UserRole
	if err := s.db.WithContext(ctx).
		Where("id IN ?", req.RoleIds).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	// 3. 校验所有 id 都存在
	if len(roles) != len(req.RoleIds) {
		exist := make(map[int64]struct{}, len(roles))
		for _, r := range roles {
			exist[r.ID] = struct{}{}
		}
		missing := make([]int64, 0)
		for _, id := range req.RoleIds {
			if _, ok := exist[id]; !ok {
				missing = append(missing, id)
			}
		}
		return nil, fmt.Errorf("角色不存在：%v", missing)
	}

	// 4. 拼接 roleName -> 以英文逗号分隔
	names := make([]string, len(roles))
	for i, r := range roles {
		names[i] = r.RoleName
	}
	user.Roles = strings.Join(names, ",")

	// 5. 更新审计字段
	token := extractTokenFromContext(ctx)
	operator, _ := s.authService.GetCurrentUsername(token)
	user.UpdatedBy = operator
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	// 6. 事务保存
	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return user, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID int64, req request.ChangePasswordRequest) (*entity.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return nil, fmt.Errorf("旧密码不正确")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败")
	}

	user.Password = string(hashed)
	user.PasswordResetAt = sql.NullTime{Time: time.Now(), Valid: true}
	token := extractTokenFromContext(ctx)
	operator, _ := s.authService.GetCurrentUsername(token)
	user.UpdatedBy = operator
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// ChangePasswordForcefully 管理员强制修改密码
func (s *UserService) ChangePasswordForcefully(ctx context.Context, userID int64, newPassword string) (*entity.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败")
	}

	user.Password = string(hashed)
	user.PasswordResetAt = sql.NullTime{Time: time.Now(), Valid: true}
	token := extractTokenFromContext(ctx)
	operator, _ := s.authService.GetCurrentUsername(token)
	user.UpdatedBy = operator
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// BlockUser 封禁用户
func (s *UserService) BlockUser(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Status = 1
	token := extractTokenFromContext(ctx)
	operator, _ := s.authService.GetCurrentUsername(token)
	user.UpdatedBy = operator
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// UnblockUser 解封用户
func (s *UserService) UnblockUser(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Status = 0
	token := extractTokenFromContext(ctx)
	operator, _ := s.authService.GetCurrentUsername(token)
	user.UpdatedBy = operator
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser 逻辑删除用户
func (s *UserService) DeleteUser(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Deleted = 1
	token := extractTokenFromContext(ctx)
	operator, _ := s.authService.GetCurrentUsername(token)
	user.UpdatedBy = operator
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// extractTokenFromContext 从 gin.Context 中提取 token
func extractTokenFromContext(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		auth := ginCtx.GetHeader("Authorization")
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

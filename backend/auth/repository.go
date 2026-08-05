package auth

import (
	"github.com/bryantaolong/user-system/types"
	"gorm.io/gorm"
)

// UserRepository 用户数据仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Insert 创建用户
func (r *UserRepository) Insert(user *model.SysUser) error {
	return r.db.Create(user).Error
}

// SelectByID 根据 ID 查询用户
func (r *UserRepository) SelectByID(id int64) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Where("id = ? AND deleted = 0", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SelectByUsername 根据用户名查询用户
func (r *UserRepository) SelectByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Where("username = ? AND deleted = 0", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SelectPageByConditions 分页查询用户列表
func (r *UserRepository) SelectPageByConditions(offset, pageSize int, req *model.UserQueryRequest) ([]model.SysUser, error) {
	var users []model.SysUser
	query := r.db.Where("deleted = 0")
	query = r.applyQuery(query, req)
	err := query.Order("updated_at ASC").Offset(offset).Limit(pageSize).Find(&users).Error
	return users, err
}

// Count 统计用户数量
func (r *UserRepository) Count(req *model.UserQueryRequest) (int64, error) {
	var count int64
	query := r.db.Model(&model.SysUser{}).Where("deleted = 0")
	query = r.applyQuery(query, req)
	err := query.Count(&count).Error
	return count, err
}

// SelectByIDList 根据 ID 列表查询用户
func (r *UserRepository) SelectByIDList(ids []int64) ([]model.SysUser, error) {
	var users []model.SysUser
	err := r.db.Where("id IN ? AND deleted = 0", ids).Find(&users).Error
	return users, err
}

// Update 更新用户信息
func (r *UserRepository) Update(user *model.SysUser) error {
	return r.db.Save(user).Error
}

// DeleteByID 逻辑删除用户
func (r *UserRepository) DeleteByID(id int64, updatedBy string) error {
	now := gorm.Expr("NOW()")
	result := r.db.Model(&model.SysUser{}).
		Where("id = ? AND deleted = 0", id).
		Updates(map[string]interface{}{
			"deleted":    1,
			"updated_at": now,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetDB 返回底层 DB 实例（供其他模块使用）
func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *UserRepository) applyQuery(query *gorm.DB, req *model.UserQueryRequest) *gorm.DB {
	if req == nil {
		return query
	}
	if req.Username != nil && *req.Username != "" {
		query = query.Where("username LIKE ?", "%"+*req.Username+"%")
	}
	if req.Phone != nil && *req.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+*req.Phone+"%")
	}
	if req.Email != nil && *req.Email != "" {
		query = query.Where("email LIKE ?", "%"+*req.Email+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	return query
}

// UserRoleRepository 用户角色数据仓库
type UserRoleRepository struct {
	db *gorm.DB
}

// NewUserRoleRepository 创建角色仓库实例
func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

// SelectOneByIsDefaultTrue 查询默认角色
func (r *UserRoleRepository) SelectOneByIsDefaultTrue() (*model.UserRole, error) {
	var role model.UserRole
	err := r.db.Where("is_default = TRUE AND deleted = 0").First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// SelectAll 查询所有角色
func (r *UserRoleRepository) SelectAll() ([]model.UserRole, error) {
	var roles []model.UserRole
	err := r.db.Where("deleted = 0").Find(&roles).Error
	return roles, err
}

// SelectByIDList 根据 ID 列表查询角色
func (r *UserRoleRepository) SelectByIDList(ids []int) ([]model.UserRole, error) {
	var roles []model.UserRole
	err := r.db.Where("id IN ? AND deleted = 0", ids).Find(&roles).Error
	return roles, err
}

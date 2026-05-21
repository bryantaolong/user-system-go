package repository

import (
	"github.com/bryan/user-system/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Insert(user *model.SysUser) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) SelectByID(id int64) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Where("id = ? AND deleted = 0", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SelectByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Where("username = ? AND deleted = 0", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SelectPage(offset, pageSize int, req *model.UserSearchRequest) ([]model.SysUser, error) {
	var users []model.SysUser
	query := r.db.Where("deleted = 0")
	query = r.applySearch(query, req)
	err := query.Order("updated_at ASC").Offset(offset).Limit(pageSize).Find(&users).Error
	return users, err
}

func (r *UserRepository) Count(req *model.UserSearchRequest) (int64, error) {
	var count int64
	query := r.db.Model(&model.SysUser{}).Where("deleted = 0")
	query = r.applySearch(query, req)
	err := query.Count(&count).Error
	return count, err
}

func (r *UserRepository) SelectByIDList(ids []int64) ([]model.SysUser, error) {
	var users []model.SysUser
	err := r.db.Where("id IN ? AND deleted = 0", ids).Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *model.SysUser) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteByID(id int64, updatedBy string) error {
	now := gorm.Expr("NOW()")
	return r.db.Model(&model.SysUser{}).
		Where("id = ? AND deleted = 0", id).
		Updates(map[string]interface{}{
			"deleted":    1,
			"updated_at": now,
			"updated_by": updatedBy,
		}).Error
}

func (r *UserRepository) applySearch(query *gorm.DB, req *model.UserSearchRequest) *gorm.DB {
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

package repository

import (
	"github.com/bryan/user-system/internal/model"
	"gorm.io/gorm"
)

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

func (r *UserRoleRepository) SelectOneByIsDefaultTrue() (*model.UserRole, error) {
	var role model.UserRole
	err := r.db.Where("is_default = TRUE AND deleted = 0").First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *UserRoleRepository) SelectAll() ([]model.UserRole, error) {
	var roles []model.UserRole
	err := r.db.Where("deleted = 0").Find(&roles).Error
	return roles, err
}

func (r *UserRoleRepository) SelectByIDList(ids []int) ([]model.UserRole, error) {
	var roles []model.UserRole
	err := r.db.Where("id IN ? AND deleted = 0", ids).Find(&roles).Error
	return roles, err
}

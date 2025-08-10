package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/bryantaolong/system/internal/model/entity"
)

type UserRoleService struct {
	db *gorm.DB
}

func NewUserRoleService(db *gorm.DB) *UserRoleService {
	return &UserRoleService{db: db}
}

// ListAll 等价于 Java 的 listAll()
func (s *UserRoleService) ListAll(ctx context.Context) ([]entity.UserRole, error) {
	var roles []entity.UserRole
	if err := s.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}

	userRoles := make([]entity.UserRole, 0, len(roles))
	for _, r := range roles {
		userRoles = append(userRoles, entity.UserRole{
			ID:       r.ID,
			RoleName: r.RoleName,
		})
	}
	return userRoles, nil
}

// FindByIds 等价于 Java 的 findByIds(Collection<Long> ids)
func (s *UserRoleService) FindByIds(ctx context.Context, ids []int64) ([]entity.UserRole, error) {
	var roles []entity.UserRole
	if err := s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

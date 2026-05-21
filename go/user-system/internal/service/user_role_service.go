package service

import (
	"context"

	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/repository"
)

type UserRoleService struct {
	roleRepo *repository.UserRoleRepository
}

func NewUserRoleService(roleRepo *repository.UserRoleRepository) *UserRoleService {
	return &UserRoleService{roleRepo: roleRepo}
}

func (s *UserRoleService) ListAll(ctx context.Context) ([]model.UserRole, error) {
	return s.roleRepo.SelectAll()
}

func (s *UserRoleService) GetDefaultRole(ctx context.Context) (*model.UserRole, error) {
	return s.roleRepo.SelectOneByIsDefaultTrue()
}

func (s *UserRoleService) ListByIDs(ctx context.Context, ids []int) ([]model.UserRole, error) {
	return s.roleRepo.SelectByIDList(ids)
}

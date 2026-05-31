package repository

import (
	"github.com/bryan/user-system/internal/model"
	"gorm.io/gorm"
)

type UserProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (r *UserProfileRepository) Insert(profile *model.UserProfile) error {
	return r.db.Create(profile).Error
}

func (r *UserProfileRepository) SelectByUserID(userID int64) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.Where("user_id = ? AND deleted = 0", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *UserProfileRepository) SelectByRealName(realName string) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.Where("real_name = ? AND deleted = 0", realName).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *UserProfileRepository) Update(profile *model.UserProfile) error {
	return r.db.Save(profile).Error
}

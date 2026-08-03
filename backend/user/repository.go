package user

import (
	"github.com/bryantaolong/user-system/types"
	"gorm.io/gorm"
)

// ProfileRepository 用户资料数据仓库
type ProfileRepository struct {
	db *gorm.DB
}

// NewProfileRepository 创建用户资料仓库实例
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Insert 创建用户资料
func (r *ProfileRepository) Insert(profile *model.UserProfile) error {
	return r.db.Create(profile).Error
}

// SelectByUserID 根据用户 ID 查询资料
func (r *ProfileRepository) SelectByUserID(userID int64) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.Where("user_id = ? AND deleted = 0", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// SelectByRealName 根据真实姓名查询资料
func (r *ProfileRepository) SelectByRealName(realName string) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.Where("real_name = ? AND deleted = 0", realName).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Update 更新用户资料
func (r *ProfileRepository) Update(profile *model.UserProfile) error {
	return r.db.Save(profile).Error
}

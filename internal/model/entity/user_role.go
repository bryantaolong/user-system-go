package entity

import (
	"time"
)

type UserRole struct {
	ID        int64     `json:"id" db:"id"`
	RoleName  string    `json:"role_name" db:"role_name"`
	IsDefault bool      `json:"is_default" db:"is_default"`
	Deleted   int       `json:"-" db:"deleted"`       // 软删除标记不暴露给前端
	Version   int       `json:"version" db:"version"` // 乐观锁版本号
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_ta"`
	CreatedBy string    `json:"createdBy" db:"created_by"`
	UpdatedBy string    `json:"updatedBy" db:"updated_by"`
}

// BeforeCreate 创建前的钩子函数，可用于设置默认值等
func (u *UserRole) BeforeCreate() {
	u.CreatedAt = time.Now()
	if u.Deleted == 0 {
		u.Deleted = 0 // 默认未删除
	}
	if u.Version == 0 {
		u.Version = 0 // 默认版本号
	}
}

package entity

import (
	"database/sql"
	"time"
)

// UserProfile 用户资料实体结构体
type UserProfile struct {
	UserID    int64        `json:"userId" db:"user_id"`
	RealName  string       `json:"realName" db:"real_name"`
	Gender    int          `json:"gender" db:"gender"`
	Birthday  sql.NullTime `json:"birthday" db:"birthday"`
	Avatar    string       `json:"avatar" db:"avatar"`
	Deleted   int          `json:"-" db:"deleted"`       // 软删除标记不暴露给前端
	Version   int          `json:"version" db:"version"` // 乐观锁版本号
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time    `json:"updatedAt" db:"updated_ta"`
	CreatedBy string       `json:"createdBy" db:"created_by"`
	UpdatedBy string       `json:"updatedBy" db:"updated_by"`
}

// TableName 返回表名
func (UserProfile) TableName() string {
	return "user_profile"
}

// BeforeUpdate 更新前的钩子函数，设置更新时间
func (up *UserProfile) BeforeUpdate() {
	up.CreatedAt = time.Now()
	up.UpdatedAt = time.Now()
	if up.Deleted == 0 {
		up.Deleted = 0 // 默认未删除
	}
	if up.Version == 0 {
		up.Version = 0 // 默认版本号
	}
}

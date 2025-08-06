package entity

import (
	"database/sql"
	"time"
)

// UserProfile 用户资料实体结构体
type UserProfile struct {
	UserID     int64        `json:"userId" db:"user_id"`
	RealName   string       `json:"realName" db:"real_name"`
	Gender     int          `json:"gender" db:"gender"`
	Birthday   sql.NullTime `json:"birthday" db:"birthday"`
	Avatar     string       `json:"avatar" db:"avatar"`
	Deleted    int          `json:"-" db:"deleted"`       // 软删除标记不暴露给前端
	Version    int          `json:"version" db:"version"` // 乐观锁版本号
	CreateTime time.Time    `json:"createTime" db:"create_time"`
	UpdateTime time.Time    `json:"updateTime" db:"update_time"`
	CreateBy   string       `json:"createBy" db:"create_by"`
	UpdateBy   string       `json:"updateBy" db:"update_by"`
}

// TableName 返回表名
func (UserProfile) TableName() string {
	return "user_profile"
}

// BeforeUpdate 更新前的钩子函数，设置更新时间
func (up *UserProfile) BeforeUpdate() {
	up.CreateTime = time.Now()
	up.UpdateTime = time.Now()
	if up.Deleted == 0 {
		up.Deleted = 0 // 默认未删除
	}
	if up.Version == 0 {
		up.Version = 0 // 默认版本号
	}
}

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
	UpdateTime time.Time    `json:"updateTime" db:"update_time"`
	UpdateBy   string       `json:"updateBy" db:"update_by"`
}

// TableName 返回表名
func (UserProfile) TableName() string {
	return "user_profile"
}

// BeforeUpdate 更新前的钩子函数，设置更新时间
func (up *UserProfile) BeforeUpdate() {
	up.UpdateTime = time.Now()
}

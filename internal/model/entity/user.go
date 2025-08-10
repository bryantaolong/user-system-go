package entity

import (
	"database/sql"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户实体结构体
type User struct {
	ID              int64        `json:"id" db:"id"`
	Username        string       `json:"username" db:"username"`
	Password        string       `json:"-" db:"password"` // 密码不序列化到JSON
	Phone           string       `json:"phone" db:"phone"`
	Email           string       `json:"email" db:"email"`
	Status          int          `json:"status" db:"status"` // 状态（0-正常，1-封禁，2-锁定）
	Roles           string       `json:"roles" db:"roles"`   // 角色标识，多个用英文逗号分隔
	LastLoginAt     sql.NullTime `json:"LastLoginAt" db:"last_login_at"`
	LastLoginIP     string       `json:"loginIp" db:"login_ip"`
	PasswordResetAt sql.NullTime `json:"passwordResetTime" db:"password_reset_at"`
	LoginFailCount  int          `json:"loginFailCount" db:"login_fail_count"`
	LockedAt        sql.NullTime `json:"lockedAt" db:"locked_at"`
	Deleted         int          `json:"-" db:"deleted"`       // 软删除标记不暴露给前端
	Version         int          `json:"version" db:"version"` // 乐观锁版本号
	CreatedAt       time.Time    `json:"createAt" db:"created_at"`
	UpdatedAt       sql.NullTime `json:"updatedAt" db:"updated_ta"`
	CreatedBy       string       `json:"createdBy" db:"created_by"`
	UpdatedBy       string       `json:"updatedBy" db:"updated_by"`
}

// TableName 返回表名
func (User) TableName() string {
	return "user"
}

// CheckPassword 检查密码是否匹配
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HashPassword 对密码进行哈希处理
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// IsAccountNonLocked 检查账户是否未锁定
func (u *User) IsAccountNonLocked() bool {
	if u.Status == 0 {
		return true
	}
	if u.Status == 2 && u.LockedAt.Valid {
		return time.Now().After(u.LockedAt.Time.Add(time.Hour))
	}
	return false
}

// IsEnabled 检查账户是否启用
func (u *User) IsEnabled() bool {
	return u.Status != 1 && u.Deleted == 0
}

// GetAuthorities 获取用户权限列表
func (u *User) GetAuthorities() []string {
	if u.Roles == "" {
		return []string{}
	}
	return strings.Split(u.Roles, ",")
}

// BeforeCreate 创建前的钩子函数，可用于设置默认值等
func (u *User) BeforeCreate() {
	u.CreatedAt = time.Now()
	if u.Status == 0 {
		u.Status = 0 // 默认状态为正常
	}
	if u.Deleted == 0 {
		u.Deleted = 0 // 默认未删除
	}
	if u.Version == 0 {
		u.Version = 0 // 默认版本号
	}
}

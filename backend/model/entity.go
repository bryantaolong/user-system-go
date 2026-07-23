package model

import "time"

// UserStatus 用户状态枚举
type UserStatus int

const (
	UserStatusNormal UserStatus = 0 // 正常
	UserStatusBanned UserStatus = 1 // 封禁
	UserStatusLocked UserStatus = 2 // 锁定
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusNormal:
		return "NORMAL"
	case UserStatusBanned:
		return "BANNED"
	case UserStatusLocked:
		return "LOCKED"
	default:
		return "UNKNOWN"
	}
}

func UserStatusFromCode(code int) UserStatus {
	return UserStatus(code)
}

// Gender 性别枚举
type Gender int

const (
	GenderFemale Gender = 0 // 女
	GenderMale   Gender = 1 // 男
)

func (g Gender) String() string {
	switch g {
	case GenderFemale:
		return "FEMALE"
	case GenderMale:
		return "MALE"
	default:
		return "UNKNOWN"
	}
}

func GenderFromCode(code int) Gender {
	return Gender(code)
}

// SysUser 用户实体
type SysUser struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username         string     `gorm:"column:username;not null" json:"username"`
	Password         string     `gorm:"column:password;not null" json:"-"`
	Phone            *string    `gorm:"column:phone" json:"phone"`
	Email            *string    `gorm:"column:email" json:"email"`
	Status           UserStatus `gorm:"column:status;default:0" json:"status"`
	Roles            *string    `gorm:"column:roles" json:"roles"`
	LastLoginAt      *time.Time `gorm:"column:last_login_at" json:"lastLoginAt"`
	LastLoginIP      *string    `gorm:"column:last_login_ip" json:"lastLoginIp"`
	LastLoginDevice  *string    `gorm:"column:last_login_device" json:"lastLoginDevice"`
	PasswordResetAt  *time.Time `gorm:"column:password_reset_at" json:"passwordResetAt"`
	LoginFailCount   *int       `gorm:"column:login_fail_count;default:0" json:"loginFailCount"`
	LockedAt         *time.Time `gorm:"column:locked_at" json:"lockedAt"`
	Deleted          int        `gorm:"column:deleted;default:0;not null" json:"deleted"`
	Version          int        `gorm:"column:version;default:0;not null" json:"version"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CreatedBy        *string    `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy        *string    `gorm:"column:updated_by" json:"updatedBy"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

// IsAccountNonLocked 判断账号是否未被锁定
func (u *SysUser) IsAccountNonLocked() bool {
	if u.Status == UserStatusNormal {
		return true
	}
	if u.Status == UserStatusLocked && u.LockedAt != nil {
		return time.Now().After(u.LockedAt.Add(time.Hour))
	}
	return false
}

// IsEnabled 判断账号是否启用
func (u *SysUser) IsEnabled() bool {
	return u.Status != UserStatusBanned && u.Deleted == 0
}

// GetAuthorities 获取用户权限列表
func (u *SysUser) GetAuthorities() []string {
	if u.Roles == nil || *u.Roles == "" {
		return []string{}
	}
	roles := splitRoles(*u.Roles)
	for i, r := range roles {
		if len(r) > 0 && r[:5] != "ROLE_" {
			roles[i] = "ROLE_" + r
		}
	}
	return roles
}

// UserProfile 用户资料实体
type UserProfile struct {
	UserID    int64      `gorm:"column:user_id;primaryKey" json:"userId"`
	RealName  *string    `gorm:"column:real_name" json:"realName"`
	Gender    *Gender    `gorm:"column:gender" json:"gender"`
	Birthday  *time.Time `gorm:"column:birthday" json:"birthday"`
	Avatar    *string    `gorm:"column:avatar" json:"avatar"`
	Deleted   int        `gorm:"column:deleted;default:0;not null" json:"deleted"`
	Version   int        `gorm:"column:version;default:0;not null" json:"version"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CreatedBy *string    `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy *string    `gorm:"column:updated_by" json:"updatedBy"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}

// UserRole 用户角色实体
type UserRole struct {
	ID         int       `gorm:"column:id;primaryKey" json:"id"`
	RoleName   string    `gorm:"column:role_name;not null" json:"roleName"`
	IsDefault  bool      `gorm:"column:is_default;default:false;not null" json:"isDefault"`
	Deleted    int       `gorm:"column:deleted;default:0;not null" json:"deleted"`
	Version    int       `gorm:"column:version;default:0;not null" json:"version"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CreatedBy  *string   `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy  *string   `gorm:"column:updated_by" json:"updatedBy"`
}

func (UserRole) TableName() string {
	return "user_role"
}

// helper
func splitRoles(s string) []string {
	var result []string
	for _, r := range splitByComma(s) {
		r = trimSpace(r)
		if r != "" {
			result = append(result, r)
		}
	}
	return result
}

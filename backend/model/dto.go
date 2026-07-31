package model

import "time"

// ========== 请求 DTO ==========

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Phone    string `json:"phone" binding:"omitempty,regexp=^1[3-9]\\d{9}$"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// UserCreateRequest 管理员创建用户请求
type UserCreateRequest struct {
	Username string  `json:"username" binding:"required,min=2,max=20"`
	Password string  `json:"password" binding:"required,min=6,max=100"`
	Phone    *string `json:"phone" binding:"omitempty,regexp=^1[3-9]\\d{9}$"`
	Email    *string `json:"email" binding:"omitempty,email"`
	RoleIDs  []int   `json:"roleIds" binding:"omitempty"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Phone    *string `json:"phone" binding:"omitempty,regexp=^1[3-9]\\d{9}$"`
	Email    *string `json:"email" binding:"omitempty,email"`
	RealName *string `json:"realName" binding:"omitempty,min=2,max=20"`
	Gender   *int    `json:"gender" binding:"omitempty"`
	Birthday *string `json:"birthday" binding:"omitempty"` // ISO8601 date
	Avatar   *string `json:"avatar" binding:"omitempty,max=500,regexp=^(https?://.*)?$"`
}

// ChangeRoleRequest 修改角色请求
type ChangeRoleRequest struct {
	RoleIDs []int `json:"roleIds" binding:"required,min=1"`
}

// UserQueryRequest 用户查询请求
type UserQueryRequest struct {
	Username *string `json:"username" form:"username" binding:"omitempty,regexp=^[a-zA-Z0-9_\\u4e00-\\u9fa5]+$"`
	Phone    *string `json:"phone" form:"phone" binding:"omitempty,regexp=^1[3-9]\\d{9}$"`
	Email    *string `json:"email" form:"email" binding:"omitempty,email,max=100"`
	Status   *int    `json:"status" form:"status" binding:"required"`
}

// UserExportRequest 用户导出请求
type UserExportRequest struct {
	FileName *string `json:"fileName" form:"fileName"`
}

// ========== VO ==========

// UserVO 用户展示 VO
type UserVO struct {
	ID              int64      `json:"id"`
	Username        string     `json:"username"`
	Email           *string    `json:"email"`
	Phone           *string    `json:"phone"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"createdAt"`
	LastLoginAt     *time.Time `json:"lastLoginAt"`
	LastLoginIP     *string    `json:"lastLoginIp"`
	LastLoginDevice *string    `json:"lastLoginDevice"`
	Roles           *string    `json:"roles"`
}

// UserProfileVO 用户资料 VO
type UserProfileVO struct {
	UserID   int64      `json:"userId"`
	Username string     `json:"username"`
	Phone    *string    `json:"phone"`
	Email    *string    `json:"email"`
	RealName *string    `json:"realName"`
	Gender   *int       `json:"gender"`
	Birthday *time.Time `json:"birthday"`
	Avatar   *string    `json:"avatar"`
}

// UserRoleOptionVO 角色选项 VO
type UserRoleOptionVO struct {
	ID       int    `json:"id"`
	RoleName string `json:"roleName"`
}

// UserExportVO 用户导出 VO
type UserExportVO struct {
	ID              int64      `json:"id"`
	Username        string     `json:"username"`
	Phone           *string    `json:"phone"`
	Email           *string    `json:"email"`
	Status          string     `json:"status"`
	Roles           *string    `json:"roles"`
	LastLoginAt     *time.Time `json:"lastLoginAt"`
	LastLoginIP     *string    `json:"lastLoginIp"`
	LastLoginDevice *string    `json:"lastLoginDevice"`
	PasswordResetAt *time.Time `json:"passwordResetAt"`
	Deleted         string     `json:"deleted"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	CreatedBy       *string    `json:"createdBy"`
	UpdatedBy       *string    `json:"updatedBy"`
}

// ========== 转换器 ==========

func ToUserVO(user *SysUser) *UserVO {
	if user == nil {
		return nil
	}
	return &UserVO{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		Status:          user.Status.String(),
		CreatedAt:       user.CreatedAt,
		LastLoginAt:     user.LastLoginAt,
		LastLoginIP:     user.LastLoginIP,
		LastLoginDevice: user.LastLoginDevice,
		Roles:           user.Roles,
	}
}

func ToUserProfileVO(user *SysUser, profile *UserProfile) *UserProfileVO {
	if user == nil {
		return nil
	}
	vo := &UserProfileVO{
		UserID:   user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		Email:    user.Email,
	}
	if profile != nil {
		vo.RealName = profile.RealName
		if profile.Gender != nil {
			g := int(*profile.Gender)
			vo.Gender = &g
		}
		vo.Birthday = profile.Birthday
		vo.Avatar = profile.Avatar
	}
	return vo
}

func ToExportVO(user *SysUser) *UserExportVO {
	if user == nil {
		return nil
	}
	deleted := "未删除"
	if user.Deleted != 0 {
		deleted = "已删除"
	}
	return &UserExportVO{
		ID:              user.ID,
		Username:        user.Username,
		Phone:           user.Phone,
		Email:           user.Email,
		Status:          convertStatus(user.Status),
		Roles:           user.Roles,
		LastLoginAt:     user.LastLoginAt,
		LastLoginIP:     user.LastLoginIP,
		LastLoginDevice: user.LastLoginDevice,
		PasswordResetAt: user.PasswordResetAt,
		Deleted:         deleted,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		CreatedBy:       user.CreatedBy,
		UpdatedBy:       user.UpdatedBy,
	}
}

func convertStatus(status UserStatus) string {
	switch status {
	case UserStatusNormal:
		return "正常"
	case UserStatusBanned:
		return "封禁"
	case UserStatusLocked:
		return "锁定"
	default:
		return "未知"
	}
}

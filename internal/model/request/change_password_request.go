package request

// ChangePasswordRequest 密码修改请求结构体
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6"` // 旧密码
	NewPassword string `json:"newPassword" binding:"required,min=6"` // 新密码
}

// ChangePasswordRequestValidationMessages 密码修改请求验证消息
var ChangePasswordRequestValidationMessages = map[string]string{
	"OldPassword.required": "旧密码不能为空",
	"OldPassword.min":      "密码至少6位",
	"NewPassword.required": "新密码不能为空",
	"NewPassword.min":      "密码至少6位",
}

package request

// UserUpdateRequest 用户更新请求结构体
type UserUpdateRequest struct {
	Username string `json:"username" binding:"omitempty,min=2,max=20"`
	Phone    string `json:"phone" binding:"omitempty,phone"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// UserUpdateRequestValidationMessages 用户更新请求验证消息
var UserUpdateRequestValidationMessages = map[string]string{
	"Username.min": "用户名长度应在2-20个字符之间",
	"Username.max": "用户名长度应在2-20个字符之间",
	"Phone.phone":  "手机号格式不正确",
	"Email.email":  "邮箱格式不正确",
}

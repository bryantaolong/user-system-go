package request

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=20"` // 用户名
	Password string `json:"password" binding:"required,min=6"`        // 密码
}

// LoginRequestValidationMessages 登录请求验证消息
var LoginRequestValidationMessages = map[string]string{
	"Username.required": "用户名不能为空",
	"Username.min":      "用户名长度应在2-20个字符之间",
	"Username.max":      "用户名长度应在2-20个字符之间",
	"Password.required": "密码不能为空",
	"Password.min":      "密码至少6位",
}

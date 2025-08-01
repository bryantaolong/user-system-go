package request

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=2,max=20"`                      // 用户名
	Password    string `json:"password" binding:"required,min=6"`                             // 密码
	PhoneNumber string `json:"phoneNumber,omitempty" binding:"omitempty,startswith=1,len=11"` // 电话号码
	Email       string `json:"email,omitempty" binding:"omitempty,email"`                     // 邮箱地址
}

// RegisterRequestValidationMessages 注册请求验证消息
var RegisterRequestValidationMessages = map[string]string{
	"Username.required":      "用户名不能为空",
	"Username.min":           "用户名长度应在2-20个字符之间",
	"Username.max":           "用户名长度应在2-20个字符之间",
	"Password.required":      "密码不能为空",
	"Password.min":           "密码至少6位",
	"PhoneNumber.startswith": "电话号码格式不正确",
	"PhoneNumber.len":        "电话号码格式不正确",
	"Email.email":            "邮箱格式不正确",
}

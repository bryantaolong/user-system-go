package request

import (
	"reflect"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// UserSearchRequest 用户搜索请求结构体
type UserSearchRequest struct {
	Username        string    `form:"username" binding:"omitempty,min=2,max=20,usernameFormat"`
	Phone           string    `form:"phone" binding:"omitempty,phone"`
	Email           string    `form:"email" binding:"omitempty,email,max=100"`
	Status          *int      `form:"status" binding:"omitempty,min=0,max=1"`
	Roles           string    `form:"roles" binding:"omitempty,rolesFormat"`
	LastLoginAt     time.Time `form:"lastLoginAt" binding:"omitempty,ltnow"`
	LastLoginIp     string    `form:"lastLoginIp" binding:"omitempty,ip"`
	PasswordResetAt time.Time `form:"passwordResetAt" binding:"omitempty,ltnow"`
	LoginFailCount  *int      `form:"loginFailCount" binding:"omitempty,min=0,max=10"`
	LockedAt        time.Time `form:"lockedAt" binding:"omitempty,ltnow"`
	Deleted         *int      `form:"deleted" binding:"omitempty,min=0,max=1"`
	Version         *int      `form:"version" binding:"omitempty,min=0"`
	CreatedAt       time.Time `form:"createdAT" binding:"omitempty,ltnow"`
	CreateTimeStart time.Time `form:"createTimeStart" binding:"omitempty,ltnow"`
	CreateTimeEnd   time.Time `form:"createTimeEnd" binding:"omitempty,ltnow"`
	UpdatedAt       time.Time `form:"updatedAt" binding:"omitempty,ltnow"`
	UpdateTimeStart time.Time `form:"updateTimeStart" binding:"omitempty,ltnow"`
	UpdateTimeEnd   time.Time `form:"updateTimeEnd" binding:"omitempty,ltnow"`
	CreatedBy       string    `form:"createdBy" binding:"omitempty,max=50"`
	UpdatedBy       string    `form:"updatedBy" binding:"omitempty,max=50"`
}

// UserSearchRequestValidationMessages 用户搜索请求验证消息
var UserSearchRequestValidationMessages = map[string]string{
	"Username.min":            "用户名长度应在2-20个字符之间",
	"Username.max":            "用户名长度应在2-20个字符之间",
	"Username.usernameFormat": "用户名只能包含中文、字母、数字和下划线",
	"Phone.phone":             "手机号格式不正确",
	"Email.email":             "邮箱格式不正确",
	"Email.max":               "邮箱长度不能超过100个字符",
	"Status.min":              "状态不合法",
	"Status.max":              "状态不合法",
	"Roles.rolesFormat":       "角色格式不正确",
	"LastLoginAt.ltnow":       "登录时间不能是未来时间",
	"LastLoginIp.ip":          "IP地址格式不正确",
	"PasswordResetAt.ltnow":   "密码重置时间不能是未来时间",
	"LoginFailCount.min":      "登录失败次数不能为负数",
	"LoginFailCount.max":      "登录失败次数超过最大值",
	"LockedAt.ltnow":          "账户锁定时间不能是未来时间",
	"Deleted.min":             "删除标记不合法",
	"Deleted.max":             "删除标记不合法",
	"Version.min":             "版本号不能为负数",
	"CreatedAt.ltnow":         "创建时间不能是未来时间",
	"CreateTimeStart.ltnow":   "开始时间不能是未来时间",
	"CreateTimeEnd.ltnow":     "结束时间不能是未来时间",
	"UpdatedAt.ltnow":         "更新时间不能是未来时间",
	"UpdateTimeStart.ltnow":   "开始时间不能是未来时间",
	"UpdateTimeEnd.ltnow":     "结束时间不能是未来时间",
	"CreatedBy.max":           "创建人名称过长",
	"UpdatedBy.max":           "更新人名称过长",
}

// ValidateUserSearchRequest 自定义验证方法
func ValidateUserSearchRequest(sl validator.StructLevel) {
	req := sl.Current().Interface().(UserSearchRequest)

	// 验证时间范围
	if !req.CreateTimeStart.IsZero() && !req.CreateTimeEnd.IsZero() {
		if req.CreateTimeEnd.Before(req.CreateTimeStart) {
			sl.ReportError(req.CreateTimeEnd, "CreateTimeEnd", "createTimeEnd", "createTimeRange", "创建时间范围不合法")
		}
		if req.CreateTimeEnd.Before(time.Now().AddDate(-10, 0, 0)) {
			sl.ReportError(req.CreateTimeEnd, "CreateTimeEnd", "createTimeEnd", "createTimeTooOld", "结束时间不能早于10年前")
		}
		if req.CreateTimeStart.After(time.Now().AddDate(0, 0, 1)) {
			sl.ReportError(req.CreateTimeStart, "CreateTimeStart", "createTimeStart", "createTimeFuture", "开始时间不能晚于明天")
		}
	}

	if !req.UpdateTimeStart.IsZero() && !req.UpdateTimeEnd.IsZero() {
		if req.UpdateTimeEnd.Before(req.UpdateTimeStart) {
			sl.ReportError(req.UpdateTimeEnd, "UpdateTimeEnd", "updateTimeEnd", "updateTimeRange", "更新时间范围不合法")
		}
		if req.UpdateTimeEnd.Before(time.Now().AddDate(-10, 0, 0)) {
			sl.ReportError(req.UpdateTimeEnd, "UpdateTimeEnd", "updateTimeEnd", "updateTimeTooOld", "结束时间不能早于10年前")
		}
		if req.UpdateTimeStart.After(time.Now().AddDate(0, 0, 1)) {
			sl.ReportError(req.UpdateTimeStart, "UpdateTimeStart", "updateTimeStart", "updateTimeFuture", "开始时间不能晚于明天")
		}
	}

	// 验证必须同时指定开始和结束时间
	if (!req.CreateTimeStart.IsZero() && req.CreateTimeEnd.IsZero()) ||
		(req.CreateTimeStart.IsZero() && !req.CreateTimeEnd.IsZero()) {
		sl.ReportError(req.CreateTimeStart, "CreateTimeStart", "createTimeStart", "timeRangeComplete", "必须同时指定开始和结束时间")
	}

	if (!req.UpdateTimeStart.IsZero() && req.UpdateTimeEnd.IsZero()) ||
		(req.UpdateTimeStart.IsZero() && !req.UpdateTimeEnd.IsZero()) {
		sl.ReportError(req.UpdateTimeStart, "UpdateTimeStart", "updateTimeStart", "timeRangeComplete", "必须同时指定开始和结束时间")
	}
}

// 注册自定义验证器
func RegisterUserSearchValidators(v *validator.Validate) {
	// 用户名格式验证
	_ = v.RegisterValidation("usernameFormat", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		matched, _ := regexp.MatchString("^[a-zA-Z0-9_\\p{Han}]+$", username)
		return matched
	})

	// 角色格式验证
	_ = v.RegisterValidation("rolesFormat", func(fl validator.FieldLevel) bool {
		roles := fl.Field().String()
		if roles == "" {
			return true
		}
		matched, _ := regexp.MatchString("^[a-zA-Z0-9_,-]*$", roles)
		return matched
	})

	// 手机号验证
	_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if phone == "" {
			return true
		}
		matched, _ := regexp.MatchString("^1[3-9]\\d{9}$", phone)
		return matched
	})

	// IP地址验证
	_ = v.RegisterValidation("ip", func(fl validator.FieldLevel) bool {
		ip := fl.Field().String()
		if ip == "" {
			return true
		}
		matched, _ := regexp.MatchString("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$", ip)
		return matched
	})

	// 时间不能是未来时间
	_ = v.RegisterValidation("ltnow", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		switch field.Kind() {
		case reflect.Struct:
			if t, ok := field.Interface().(time.Time); ok {
				return t.IsZero() || t.Before(time.Now())
			}
		}
		return true
	})

	// 注册结构体验证
	v.RegisterStructValidation(ValidateUserSearchRequest, UserSearchRequest{})
}

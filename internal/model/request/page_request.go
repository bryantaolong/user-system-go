package request

// PageRequest 分页请求结构体
type PageRequest struct {
	PageNum  int64 `json:"pageNum" form:"pageNum" binding:"min=1"`   // 页码
	PageSize int64 `json:"pageSize" form:"pageSize" binding:"min=1"` // 每页大小
}

// PageRequestValidationMessages 分页请求验证消息
var PageRequestValidationMessages = map[string]string{
	"PageNum.min":  "页码不能小于1",
	"PageSize.min": "每页大小不能小于1",
}

// DefaultPageRequest 默认分页请求
var DefaultPageRequest = PageRequest{
	PageNum:  1,
	PageSize: 10,
}

// GetOffset 计算偏移量
func (p *PageRequest) GetOffset() int64 {
	return (p.PageNum - 1) * p.PageSize
}

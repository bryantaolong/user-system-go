package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HttpStatus 业务状态码枚举
const (
	StatusSuccess       = 200
	StatusBadRequest    = 400
	StatusUnauthorized  = 401
	StatusForbidden     = 403
	StatusNotFound      = 404
	StatusConflict      = 409
	StatusInternalError = 500
)

// Result 统一响应结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageResult 分页结果
type PageResult struct {
	Rows     interface{} `json:"rows"`
	Total    int64       `json:"total"`
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
	Pages    int64       `json:"pages"`
}

func NewPageResult(rows interface{}, total int64, pageNum, pageSize int) PageResult {
	var pages int64
	if total == 0 || pageSize == 0 {
		pages = 0
	} else {
		pages = (total + int64(pageSize) - 1) / int64(pageSize)
	}
	return PageResult{
		Rows:     rows,
		Total:    total,
		PageNum:  pageNum,
		PageSize: pageSize,
		Pages:    pages,
	}
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code:    StatusSuccess,
		Message: "成功",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, msg string) {
	c.JSON(http.StatusOK, Result{
		Code:    httpStatus,
		Message: msg,
		Data:    nil,
	})
}

func ErrorWithHTTPStatus(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Result{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// BusinessException 业务异常
type BusinessException struct {
	Message string
}

func NewBusinessException(msg string) *BusinessException {
	return &BusinessException{Message: msg}
}
func (e *BusinessException) Error() string {
	return e.Message
}

// ResourceNotFoundException 资源不存在异常
type ResourceNotFoundException struct {
	Message string
}

func NewResourceNotFoundException(msg string) *ResourceNotFoundException {
	return &ResourceNotFoundException{Message: msg}
}
func (e *ResourceNotFoundException) Error() string {
	return e.Message
}

// UnauthorizedException 未授权异常
type UnauthorizedException struct {
	Message string
}

func NewUnauthorizedException(msg string) *UnauthorizedException {
	return &UnauthorizedException{Message: msg}
}
func (e *UnauthorizedException) Error() string {
	return e.Message
}

// OptimisticLockException 乐观锁冲突异常
type OptimisticLockException struct {
	Message string
}

func NewOptimisticLockException(msg string) *OptimisticLockException {
	return &OptimisticLockException{Message: msg}
}
func (e *OptimisticLockException) Error() string {
	return e.Message
}

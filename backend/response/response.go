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

// BusinessError 业务异常
type BusinessError struct {
	Message string
}

func NewBusinessError(msg string) *BusinessError {
	return &BusinessError{Message: msg}
}
func (e *BusinessError) Error() string {
	return e.Message
}

// ResourceNotFoundError 资源不存在异常
type ResourceNotFoundError struct {
	Message string
}

func NewResourceNotFoundError(msg string) *ResourceNotFoundError {
	return &ResourceNotFoundError{Message: msg}
}
func (e *ResourceNotFoundError) Error() string {
	return e.Message
}

// UnauthorizedError 未授权异常
type UnauthorizedError struct {
	Message string
}

func NewUnauthorizedError(msg string) *UnauthorizedError {
	return &UnauthorizedError{Message: msg}
}
func (e *UnauthorizedError) Error() string {
	return e.Message
}

// PersistenceError 持久化异常
type PersistenceError struct {
	Message string
}

func NewPersistenceError(msg string) *PersistenceError {
	return &PersistenceError{Message: msg}
}
func (e *PersistenceError) Error() string {
	return e.Message
}

// OptimisticLockError 乐观锁冲突异常
type OptimisticLockError struct {
	Message string
}

func NewOptimisticLockError(msg string) *OptimisticLockError {
	return &OptimisticLockError{Message: msg}
}
func (e *OptimisticLockError) Error() string {
	return e.Message
}

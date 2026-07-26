package middleware

import (
	"net/http"

	"github.com/bryantaolong/user-system/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 全局异常处理中间件
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理 gin 绑定错误等
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			logger.Error("请求处理异常", zap.Error(err))

			switch e := err.(type) {
			case *response.BusinessError:
				response.Error(c, response.StatusInternalError, e.Message)
			case *response.ResourceNotFoundError:
				response.Error(c, response.StatusNotFound, e.Message)
			case *response.UnauthorizedError:
				response.ErrorWithHTTPStatus(c, http.StatusUnauthorized, response.StatusUnauthorized, e.Message)
			default:
				response.Error(c, response.StatusInternalError, "服务繁忙，请稍后重试")
			}
		}
	}
}

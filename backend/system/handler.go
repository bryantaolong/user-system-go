package system

import (
	"net/http"
	"strconv"

	"github.com/bryan/user-system/model"
	"github.com/bryan/user-system/response"
	"github.com/gin-gonic/gin"
)

// Handler 系统模块处理器
type Handler struct {
	logSvc *LogService
}

// NewHandler 创建系统模块处理器实例
func NewHandler(logSvc *LogService) *Handler {
	return &Handler{logSvc: logSvc}
}

// ListLatestLogs 获取最新的应用日志内容
func (h *Handler) ListLatestLogs(c *gin.Context) {
	lines, _ := strconv.Atoi(c.DefaultQuery("lines", "200"))
	fileName := c.Query("file")

	logs, err := h.logSvc.ListLatestLogs(fileName, lines)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, logs)
}

// ListLogFiles 获取可用的日志文件列表
func (h *Handler) ListLogFiles(c *gin.Context) {
	files := h.logSvc.ListLogFiles()
	c.JSON(http.StatusOK, model.Result{
		Code:    200,
		Message: "成功",
		Data:    files,
	})
}

// handleError 统一错误处理
func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *response.BusinessError:
		response.Error(c, response.StatusInternalError, e.Message)
	case *response.ResourceNotFoundError:
		response.Error(c, response.StatusNotFound, e.Message)
	case *response.UnauthorizedError:
		response.ErrorWithHTTPStatus(c, http.StatusUnauthorized, response.StatusUnauthorized, e.Message)
	case *response.OptimisticLockError:
		response.Error(c, response.StatusConflict, e.Message)
	default:
		response.Error(c, response.StatusInternalError, "服务繁忙，请稍后重试")
	}
}

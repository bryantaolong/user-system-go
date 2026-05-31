package handler

import (
	"net/http"
	"strconv"

	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/service"
	"github.com/gin-gonic/gin"
)

// SystemLogHandler 系统日志处理器
type SystemLogHandler struct {
	logSvc *service.LogService
}

func NewSystemLogHandler(logSvc *service.LogService) *SystemLogHandler {
	return &SystemLogHandler{logSvc: logSvc}
}

// ListLatestLogs 获取最新的应用日志内容
func (h *SystemLogHandler) ListLatestLogs(c *gin.Context) {
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
func (h *SystemLogHandler) ListLogFiles(c *gin.Context) {
	files := h.logSvc.ListLogFiles()
	c.JSON(http.StatusOK, response.Result{
		Code:    200,
		Message: "成功",
		Data:    files,
	})
}

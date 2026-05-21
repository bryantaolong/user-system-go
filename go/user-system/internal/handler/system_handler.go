package handler

import (
	"net/http"
	"strconv"

	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/service"
	"github.com/gin-gonic/gin"
)

type UserExportHandler struct {
	exportSvc *service.UserExportService
}

func NewUserExportHandler(exportSvc *service.UserExportService) *UserExportHandler {
	return &UserExportHandler{exportSvc: exportSvc}
}

// ExportAllUsers 导出所有用户数据为 Excel
func (h *UserExportHandler) ExportAllUsers(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	_ = pageNum // 导出服务内部处理分页

	f, err := h.exportSvc.ExportAllUsers()
	if err != nil {
		handleError(c, err)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=users.xlsx")

	if err := f.Write(c.Writer); err != nil {
		response.Error(c, response.StatusInternalError, "导出文件写入失败")
	}
}

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

package service

import (
	"time"

	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/repository"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type UserExportService struct {
	userRepo *repository.UserRepository
	logger   *zap.Logger
}

func NewUserExportService(userRepo *repository.UserRepository, logger *zap.Logger) *UserExportService {
	return &UserExportService{userRepo: userRepo, logger: logger}
}

// ExportAllUsers 全量导出用户数据为 Excel
func (s *UserExportService) ExportAllUsers() (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "用户列表"
	f.SetSheetName("Sheet1", sheet)

	// 设置表头
	headers := []string{"ID", "用户名", "手机", "邮箱", "状态", "角色", "最后登录时间", "最后登录IP", "最后登录设备", "密码重置时间", "删除状态", "创建时间", "更新时间", "创建人", "更新人"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 分批查询写入
	pageNum := 1
	pageSize := 1000
	row := 2

	for {
		offset := (pageNum - 1) * pageSize
		users, err := s.userRepo.SelectPage(offset, pageSize, nil)
		if err != nil || len(users) == 0 {
			break
		}

		for _, user := range users {
			vo := model.ToExportVO(&user)
			cols := []interface{}{
				vo.ID, vo.Username, derefStr(vo.Phone), derefStr(vo.Email),
				vo.Status, derefStr(vo.Roles),
				formatTime(vo.LastLoginAt), derefStr(vo.LastLoginIP),
				derefStr(vo.LastLoginDevice), formatTime(vo.PasswordResetAt),
				vo.Deleted, vo.CreatedAt.Format("2006-01-02 15:04:05"),
				vo.UpdatedAt.Format("2006-01-02 15:04:05"),
				derefStr(vo.CreatedBy), derefStr(vo.UpdatedBy),
			}
			for i, v := range cols {
				cell, _ := excelize.CoordinatesToCellName(i+1, row)
				f.SetCellValue(sheet, cell, v)
			}
			row++
		}

		s.logger.Info("已导出数据", zap.Int("count", row-2))
		pageNum++
	}

	// 设置表头样式
	style, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#CCFFCC"}},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheet, "A1", "O1", style)

	s.logger.Info("导出完成", zap.Int("total", row-2))
	return f, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

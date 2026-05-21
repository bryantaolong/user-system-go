package service

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bryan/user-system/internal/config"
	"github.com/bryan/user-system/internal/pkg/response"
	"go.uber.org/zap"
)

type LogService struct {
	logger *zap.Logger
}

func NewLogService(logger *zap.Logger) *LogService {
	return &LogService{logger: logger}
}

// ListLatestLogs 读取指定日志文件最近 N 行
func (s *LogService) ListLatestLogs(fileName string, maxLines int) ([]string, error) {
	limit := maxLines
	if limit < 1 {
		limit = 1
	}
	if limit > 2000 {
		limit = 2000
	}

	path := s.resolveLogPath(fileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, response.NewBusinessException("日志文件不存在，请检查日志配置")
	}

	allLines, err := s.readAllLines(path)
	if err != nil {
		return nil, response.NewBusinessException("读取日志文件失败，请稍后重试")
	}
	if len(allLines) == 0 {
		return []string{}, nil
	}

	fromIndex := len(allLines) - limit
	if fromIndex < 0 {
		fromIndex = 0
	}
	return allLines[fromIndex:], nil
}

// ListLogFiles 列出日志目录下所有可用日志文件
func (s *LogService) ListLogFiles() []string {
	logDir := s.getLogsDirectory()
	defaultLogPath := s.resolveDefaultLogPath()

	entries, err := os.ReadDir(logDir)
	if err != nil {
		return []string{filepath.Base(defaultLogPath)}
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".log") || strings.HasSuffix(name, ".gz") {
			files = append(files, name)
		}
	}

	if len(files) == 0 {
		return []string{filepath.Base(defaultLogPath)}
	}

	sort.Strings(files)
	return files
}

func (s *LogService) resolveLogPath(fileName string) string {
	if fileName == "" {
		return s.resolveDefaultLogPath()
	}
	logDir := s.getLogsDirectory()
	path := filepath.Join(logDir, fileName)
	absPath, _ := filepath.Abs(path)
	absLogDir, _ := filepath.Abs(logDir)
	if !strings.HasPrefix(absPath, absLogDir) {
		panic("非法的日志文件路径")
	}
	return path
}

func (s *LogService) resolveDefaultLogPath() string {
	logFileName := "logs/user-system.log"
	if config.AppConfig != nil && config.AppConfig.Logging.File != "" {
		logFileName = config.AppConfig.Logging.File
	}
	path := logFileName
	if !filepath.IsAbs(path) {
		wd, _ := os.Getwd()
		path = filepath.Join(wd, path)
	}
	absPath, _ := filepath.Abs(path)
	return absPath
}

func (s *LogService) getLogsDirectory() string {
	defaultLogPath := s.resolveDefaultLogPath()
	dir := filepath.Dir(defaultLogPath)
	return dir
}

func (s *LogService) readAllLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if strings.HasSuffix(path, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		return readLines(gz)
	}

	return readLines(f)
}

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

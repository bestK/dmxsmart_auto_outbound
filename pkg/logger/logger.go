package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var Logger *slog.Logger

func Init() {
	// 确保logs目录存在
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	// 创建日志文件，按日期命名
	logFileName := filepath.Join(logsDir, time.Now().Format("2006-01-02")+".log")

	// 创建文件处理器，启用自动按小时切割
	fileHandler := handler.MustRotateFile(logFileName, handler.EveryHour)

	// 创建控制台处理器
	consoleHandler := handler.NewConsoleHandler(slog.AllLevels)

	// 创建logger并设置多个handler
	logger := slog.NewWithHandlers(consoleHandler, fileHandler)
	Logger = logger
}

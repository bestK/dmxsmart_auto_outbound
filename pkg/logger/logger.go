package logger

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var Logger *slog.Logger

func Init() {
	logger := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels))
	Logger = logger
}

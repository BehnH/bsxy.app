package logger

import "go.uber.org/zap"

func SLog() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
		}
	}(logger) // flushes buffer, if any
	return logger.Sugar()
}

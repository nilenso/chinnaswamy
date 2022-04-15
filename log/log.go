package log

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func InitLogger() error {
	plainLogger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	logger = plainLogger.Sugar()
	return nil
}

func SyncLogs() error {
	err := logger.Sync()
	if err != nil {
		return err
	}
	return nil
}

func Errorw(msg string, rest ...interface{}) {
	logger.Errorw(msg, rest...)
}

func Infow(msg string, rest ...interface{}) {
	logger.Infow(msg, rest...)
}

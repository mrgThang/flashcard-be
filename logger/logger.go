package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() error {
	var err error
	log, err = zap.NewProduction()
	return err
}

func Error(s string, fields ...zap.Field) {
	log.Error(s, fields...)
}

func Info(s string, fields ...zap.Field) {
	log.Info(s, fields...)
}

func Warn(s string, fields ...zap.Field) {
	log.Warn(s, fields...)
}

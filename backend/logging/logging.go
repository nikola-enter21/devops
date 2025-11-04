package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.SugaredLogger
)

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	base, err := config.Build()
	if err != nil {
		panic(err)
	}
	logger = base.Sugar()
}

func MustNewLogger() *zap.SugaredLogger {
	return logger
}

func Sync() {
	_ = logger.Sync()
}

package logging

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once   sync.Once
	logger *zap.SugaredLogger
)

func MustNewLogger() *zap.SugaredLogger {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		base, err := config.Build()
		if err != nil {
			panic(err)
		}
		logger = base.Sugar()
	})
	return logger
}

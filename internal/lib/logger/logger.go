package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger(env string) *zap.Logger {
	var logger *zap.Logger
	var cnf zap.Config

	switch env {
	case "local":
		cnf = zap.NewDevelopmentConfig()
	case "prod":
		cnf = zap.NewProductionConfig()
	default:
		panic("Ошибка чтения env")
	}

	cnf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cnf.OutputPaths = []string{"stdout", "local.log"}
	cnf.ErrorOutputPaths = []string{"stderr", "error.log"}

	var err error
	logger, err = cnf.Build()
	if err != nil {
		panic("Ошибка инициализации логгера: " + err.Error())
	}
	return logger
}

package main

import (
	"github.com/DrusGalkin/auth-grpc-service/internal/app"
	"github.com/DrusGalkin/auth-grpc-service/internal/config"
	log "github.com/DrusGalkin/auth-grpc-service/pkg/lib/logger"
	"go.uber.org/zap"
)

// Команда для старта
// go run cmd/auth/main.go --config=./config/local.yaml
func main() {
	// Конфиг
	cfg := config.MustLoadConfig()

	// Логгер
	logger := log.SetupLogger(cfg.Env)
	defer logger.Sync()
	logger.Info("Запуск приложения", zap.Any("Конфиг найден", cfg))

	// Приложение
	application := app.New(logger, cfg.GRPC.Port, cfg.DBUrl, cfg.Secret, cfg.AccessTime, cfg.RefreshTime)

	// gRPC сервер
	application.GRPCServer.MustRun()

	//// Остановка gRPC сервера
	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	//
	//sign := <-stop
	//logger.Info("Приложение остановлено", zap.Any("main", sign.String()))
	//
	//application.GRPCServer.Stop()
}

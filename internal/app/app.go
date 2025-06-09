package app

import (
	grpcapp "github.com/DrusGalkin/auth-grpc-service/internal/app/grpc"
	"github.com/DrusGalkin/auth-grpc-service/internal/domain/models"
	"github.com/DrusGalkin/auth-grpc-service/internal/services"
	storageDB "github.com/DrusGalkin/auth-grpc-service/internal/storage/mysql"
	"go.uber.org/zap"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	logger *zap.Logger,
	grpcPort int,
	storagePath string,
	secret []byte,
	tokenAccess time.Duration,
	tokenRefresh time.Duration,
) *App {
	storage := storageDB.New(storagePath)
	authService := services.New(
		logger,
		models.SecretApp{
			Secret: secret,
		},
		storage,
		storage,
		tokenAccess,
		tokenRefresh,
	)

	grpcApp := grpcapp.NewGRPCApp(logger, grpcPort, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}

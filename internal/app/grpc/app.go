package grpc

import (
	"fmt"
	authGRPC "github.com/DrusGalkin/auth-grpc-service/internal/grpc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	logger     *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewGRPCApp(logger *zap.Logger, port int, authService authGRPC.Auth) *App {
	gRPCServer := grpc.NewServer()
	authGRPC.Register(gRPCServer, authService)
	return &App{
		logger:     logger,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "grpcApp.run"
	log := a.logger.With(
		zap.String("op", op),
		zap.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Старт gRPC сервар")

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcApp.Stop"
	a.logger.With(zap.String("op", op)).
		Info(
			"Остановка gRPC сервера",
			zap.Int("port", a.port),
		)
	a.gRPCServer.GracefulStop()
}

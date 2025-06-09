package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/DrusGalkin/auth-grpc-service/internal/services"
	"github.com/DrusGalkin/auth-grpc-service/internal/storage"
	"github.com/DrusGalkin/auth-grpc-service/pkg/lib/jwt"
	pk "github.com/DrusGalkin/auth-protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (*jwt.VerifyResponse, error)

	Register(
		ctx context.Context,
		email string,
		username string,
		password string,
	) (int64, error)

	IsAdmin(
		ctx context.Context,
		userId int64,
	) (bool, error)

	Refresh(
		ctx context.Context,
		refreshToken string,
	) (*jwt.VerifyResponse, error)

	ValidToken(
		ctx context.Context,
		token string,
	) (*jwt.Claim, error)
}

type serverAPI struct {
	pk.UnimplementedAuthServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	pk.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *pk.LoginRequest) (*pk.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Неверный логин или пароль: %s, %s", req.GetEmail(), req.GetPassword()))
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Ошибка входа пользователя")
		}

		return nil, status.Error(codes.Internal, "Ошибка входа пользователя")
	}

	return &pk.LoginResponse{
		Access:  token.Access,
		Refresh: token.Refresh,
	}, err
}

func (s *serverAPI) Register(ctx context.Context, req *pk.RegisterRequest) (*pk.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" || req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Неверный логин или пароль: %s, %s, %s", req.GetEmail(), req.GetUsername(), req.GetPassword()))
	}

	uid, err := s.auth.Register(ctx, req.GetEmail(), req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			return nil, status.Error(codes.AlreadyExists, "Такой пользователь уже существует")
		}

		return nil, status.Error(codes.Internal, "Ошибка регистрации пользователя")
	}

	return &pk.RegisterResponse{
		UserId: uid,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *pk.IsAdminRequest) (*pk.IsAdminResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.Internal, "Невалидный id пользователя")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "Пользователь не найден")
		}
		return nil, status.Error(codes.Internal, "Ошибка поиска админа")
	}

	return &pk.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) Refresh(ctx context.Context, req *pk.RefreshRequest) (*pk.RefreshResponse, error) {
	tokens, err := s.auth.Refresh(ctx, req.GetRefresh())
	if err != nil {
		return nil, status.Error(codes.Internal, "Невалидный токен")
	}

	return &pk.RefreshResponse{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}
func (s *serverAPI) ValidToken(ctx context.Context, req *pk.ValidTokenRequest) (*pk.ValidTokenResponse, error) {
	claim, err := s.auth.ValidToken(ctx, req.Access)
	if err != nil {
		return nil, status.Error(codes.Internal, "Невалидный токен")
	}

	return &pk.ValidTokenResponse{
		UserId: int64(claim.UserID),
		Email:  claim.Email,
	}, nil
}

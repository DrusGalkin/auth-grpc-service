package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/DrusGalkin/auth-grpc-service/internal/domain/models"
	"github.com/DrusGalkin/auth-grpc-service/internal/lib/jwt"
	"github.com/DrusGalkin/auth-grpc-service/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth struct {
	log          *zap.Logger
	secret       models.SecretApp
	userProvider UserProvider
	userSaver    UserSaver
	tokenAccess  time.Duration
	tokenRefresh time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		username string,
		hashPassword []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, uid int64) (bool, error)
}

var (
	ErrInvalidCredentials = errors.New("Неверный логин или пароль")
)

func New(
	log *zap.Logger,
	secret models.SecretApp,
	userProvider UserProvider,
	userSaver UserSaver,
	tokenAccess time.Duration,
	tokenRefresh time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		secret:       secret,
		userProvider: userProvider,
		userSaver:    userSaver,
		tokenAccess:  tokenAccess,
		tokenRefresh: tokenRefresh,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (*jwt.VerifyResponse, error) {
	const op = "auth.Login"
	log := a.log.With(
		zap.String("op", op),
		zap.String("email", email),
	)

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("Пользователь не найден", zap.Error(err))

			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Warn("Ошибка получения пользователя", zap.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.HashPassword, []byte(password)); err != nil {
		log.Warn("Невалидные данные", zap.Error(err))

		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	tokens, err := jwt.NewTokens(user, a.secret, a.tokenAccess, a.tokenRefresh)
	if err != nil {
		log.Warn("Ошибка генерации токенов", zap.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (a *Auth) Register(
	ctx context.Context,
	email string,
	username string,
	password string,
) (int64, error) {
	const op = "auth.Register"

	log := a.log.With(
		zap.String("op", op),
		zap.String("email", email),
	)

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Ошибка генерации хеша пароля", zap.Error(err))

		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, username, hashPass)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			log.Warn("Пользователь уже существует", zap.Error(err))

			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExist)
		}

		log.Error("Ошибка создания пользователя", zap.Error(err))

		return 0, fmt.Errorf("%s, %w", op, err)
	}

	log.Info("Пользователя зарегистрирован")

	return id, nil
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userId int64,
) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		zap.String("op", op),
		zap.Int64("userId", userId),
	)

	admin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Проверка пользователя на админа", zap.Bool("admin", admin))

	return admin, nil
}

func (s *Auth) Refresh(ctx context.Context, refreshToken string) (*jwt.VerifyResponse, error) {

	tokens, err := jwt.RefreshToken(refreshToken, s.secret, s.tokenAccess, s.tokenRefresh)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *Auth) ValidToken(ctx context.Context, token string) (*jwt.Claim, error) {
	claim, err := jwt.ValidToken(token, s.secret)
	if err != nil {
		return nil, err
	}
	return claim, nil
}

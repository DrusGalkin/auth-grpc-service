package jwt

import (
	"fmt"
	"github.com/DrusGalkin/auth-grpc-service/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claim struct {
	UserID int
	Email  string
	jwt.RegisteredClaims
}

type VerifyResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func NewTokens(user models.User, app models.SecretApp, timeAccess time.Duration, timeRefresh time.Duration) (*VerifyResponse, error) {
	refresh, access, err := GenerateToken(user, app, timeAccess, timeRefresh)
	if err != nil {
		return nil, err
	}

	return &VerifyResponse{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func RefreshToken(refresh string, app models.SecretApp, timeAccess time.Duration, timeRefresh time.Duration) (*VerifyResponse, error) {
	claims, err := ValidToken(refresh, app)
	if err != nil {
		return nil, fmt.Errorf("Ошибка проверки токена: %w", err)
	}

	user := models.User{
		ID:    claims.UserID,
		Email: claims.Email,
	}

	return NewTokens(user, app, timeAccess, timeRefresh)
}

func ValidToken(someToken string, app models.SecretApp) (*Claim, error) {
	token, err := jwt.ParseWithClaims(someToken, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return app.Secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Ошибка парса, невалидный токен: %w", err)
	}

	if claim, ok := token.Claims.(*Claim); ok && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("Невалидный токен")

}

func GenerateToken(user models.User, app models.SecretApp, timeAccess time.Duration, timeRefresh time.Duration) (string, string, error) {
	claimsRefresh := Claim{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeRefresh)),
		},
	}

	claimsAccess := Claim{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeAccess)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	refresh, err := refreshToken.SignedString([]byte(app.Secret))
	if err != nil {
		return "", "", fmt.Errorf("Ошибка генерации Refresh токена: %w", err)
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	access, err := accessToken.SignedString([]byte(app.Secret))
	if err != nil {
		return "", "", fmt.Errorf("Ошибка генерации Access токена: %w", err)
	}

	return refresh, access, nil
}

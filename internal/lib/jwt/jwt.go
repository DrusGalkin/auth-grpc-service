package jwt

import (
	"github.com/DrusGalkin/auth-grpc-service/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, app models.SecretApp, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(duration).Unix(),
	})

	return token.SignedString([]byte(app.Secret))
}

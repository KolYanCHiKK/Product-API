package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTAuth struct {
	Secret string
}

func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{
		Secret: secret,
	}
}

func (j *JWTAuth) CreateJWT(sessionId uuid.UUID, phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sessionId": sessionId,
		"phone":     phone,
	})

	singedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return singedToken, nil
}

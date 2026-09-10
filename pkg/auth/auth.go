package auth

import (
	"app/product-api/pkg/logs"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Статья про работу с JWT в Golang: https://allcourses.io/blog/rabota-s-jwt-v-golang/

type JWTAuth struct {
	Secret string
}

type DecodedToken struct {
	SessionId uuid.UUID
	Phone     string
}

func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{
		Secret: secret,
	}
}

func NewDecodedToken(sessionId uuid.UUID, phone string) *DecodedToken {
	return &DecodedToken{
		SessionId: sessionId,
		Phone:     phone,
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

func (j *JWTAuth) DecodeToken(authTokenStr string) (bool, *DecodedToken) {
	token, err := jwt.Parse(authTokenStr, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("unexpected signing algorithm: %v", token.Header["alg"])
		}

		return []byte(j.Secret), nil
	})
	if err != nil {
		logs.AddErrLog(logrus.Fields{
			"Error": err,
		}, "Произошла ошибка декодирования токена")
		return false, nil

	}

	// Если токен валидный, то вернем структуру декодирования
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sessionId, err := uuid.Parse(claims["sessionId"].(string))
		if err != nil {
			return false, nil
		}
		return true, NewDecodedToken(sessionId, claims["phone"].(string))
	}

	// Иначе запишем ошибку
	logs.AddErrLog(logrus.Fields{
		"Error": "Token was not decoded",
	}, "Произошла ошибка декодирования токена")
	return false, nil
}

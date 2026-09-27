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
	UserId    int
	Phone     string
}

type Claims struct {
	SessionID uuid.UUID `json:"sessionId"`
	UserID    int       `json:"userId"`
	Phone     string    `json:"phone"`

	jwt.RegisteredClaims
}

func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{
		Secret: secret,
	}
}

func NewDecodedToken(sessionId uuid.UUID, userId int, phone string) *DecodedToken {
	return &DecodedToken{
		SessionId: sessionId,
		UserId:    userId,
		Phone:     phone,
	}
}

func (j *JWTAuth) CreateJWT(sessionId uuid.UUID, userId int, phone string) (string, error) {
	claims := Claims{
		SessionID: sessionId,
		UserID:    userId,
		Phone:     phone,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	singedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return singedToken, nil
}

func (j *JWTAuth) DecodeToken(authTokenStr string) (bool, *DecodedToken) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		authTokenStr,
		claims,
		func(token *jwt.Token) (any, error) {
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
	if decodedClaims, ok := token.Claims.(*Claims); ok && token.Valid {
		return true, NewDecodedToken(decodedClaims.SessionID, decodedClaims.UserID, decodedClaims.Phone)
	}

	// Иначе запишем ошибку
	logs.AddErrLog(logrus.Fields{
		"Error": "Token was not decoded",
	}, "Произошла ошибка декодирования токена")
	return false, nil
}

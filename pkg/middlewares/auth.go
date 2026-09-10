package middlewares

import (
	"app/product-api/pkg/auth"
	"context"
	"net/http"
	"strings"
)

type AuthMiddlewares struct {
	*auth.JWTAuth
}

const (
	PhoneKey     string = "phone"
	SessionIdKey string = "sessionId"
)

func NewAuthMiddlewares(j *auth.JWTAuth) *AuthMiddlewares {
	return &AuthMiddlewares{JWTAuth: j}
}

func WriteUnauthedResponse(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func (j *AuthMiddlewares) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost && (req.URL.Path == "/users/login" || req.URL.Path == "/users/confirm-code") {
			next.ServeHTTP(w, req)
			return
		}

		authToken := req.Header.Get("Authorization")
		if !strings.HasPrefix(authToken, "Bearer ") {
			WriteUnauthedResponse(w)
			return
		}

		authToken = strings.TrimPrefix(authToken, "Bearer ")
		isValid, decodedPayload := j.DecodeToken(authToken)
		if !isValid {
			WriteUnauthedResponse(w)
			return
		}

		ctx := context.WithValue(
			context.WithValue(req.Context(), PhoneKey, decodedPayload.Phone),
			SessionIdKey,
			decodedPayload.SessionId,
		)

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

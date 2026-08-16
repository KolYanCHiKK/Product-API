package middlewares

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Middleware func(next http.Handler) http.Handler

func CallMiddleware(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		if origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, req)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			return
		}

		next.ServeHTTP(w, req)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		wrapper := &WrapperWriter{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(wrapper, req)

		log.WithFields(log.Fields{
			"method": req.Method,
			"status": wrapper.Status,
			"path":   req.URL.Path,
		}).Info("Получен HTTP запрос")
	})
}

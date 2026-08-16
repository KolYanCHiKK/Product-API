package middlewares

import "net/http"

type WrapperWriter struct {
	http.ResponseWriter
	Status int
}

func (w *WrapperWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

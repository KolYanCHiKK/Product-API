package responce

import (
	"encoding/json"
	"net/http"
)

type BaseError struct {
	Status string   `json:"status"`
	Errors []string `json:"errors"`
}

func CreateResponse(w http.ResponseWriter, statusCode int, body any) error {
	AddBaseResponseHeaders(w, statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return err
	}
	return nil
}

func CreateErrResponse(w http.ResponseWriter, statusCode int, err ...string) {
	body := BaseError{
		Status: "Error",
		Errors: err,
	}
	AddBaseResponseHeaders(w, statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func AddBaseResponseHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

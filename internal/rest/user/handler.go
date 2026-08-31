package user

import (
	"app/product-api/configs"
	"app/product-api/pkg/responce"
	"app/product-api/pkg/validation"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Service *Service
	*configs.Config
}

func NewHandler(mux *http.ServeMux, service *Service, conf *configs.Config) {
	h := &Handler{
		Service: service,
		Config:  conf,
	}

	mux.HandleFunc("POST /users/login", h.Login())
	mux.HandleFunc("POST /users/confirm-code", h.ConfirmCode())
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body LoginRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		errs := validation.ValidateBody(
			body,
			validation.UserAuthValidate,
			RegisterPhoneValidateParameters,
		)
		if errs != nil {
			responce.CreateErrResponse(w, 400, errs...)
			return
		}

		session, httpCode, err := h.Service.LoginUser(body.Phone)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		resp := &LoginResponse{SessionId: session.SessionId}
		err = responce.CreateResponse(w, httpCode, resp)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

	}
}

func (h *Handler) ConfirmCode() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body ConfirmCodeRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
		}

		errs := validation.ValidateBody(
			body,
			validation.AuthConfirmValidate,
			func() (*validator.Validate, error) {
				return validator.New(), nil
			},
		)
		if errs != nil {
			responce.CreateErrResponse(w, 400, errs...)
			return
		}

		userParams, err := h.Service.ConfirmCode(body.SessionId, body.Code)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		token, err := h.Auth.JWT.CreateJWT(body.SessionId, userParams.Phone)
		if err != nil {
			responce.CreateErrResponse(w, 401, err.Error())
			return
		}

		resp := ConfirmCodeResponse{
			AuthParameters: &AuthParameters{
				Token:     token,
				TokenType: "Bearer",
			},
			User: &UserParameters{
				UserId:    userParams.UserId,
				Surname:   userParams.Surname,
				Name:      userParams.Name,
				Phone:     userParams.Phone,
				City:      userParams.City,
				Address:   userParams.Address,
				BirthDate: userParams.BirthDate,
				Email:     userParams.Email,
			},
		}

		err = responce.CreateResponse(w, 200, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
		}
	}
}

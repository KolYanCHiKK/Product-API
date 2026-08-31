package user

import (
	"app/product-api/pkg/validation"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserParameters struct {
	UserId    int        `json:"userId,omitempty"`
	Surname   string     `json:"surname,omitempty"`
	Name      string     `json:"name,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	City      string     `json:"city,omitempty"`
	Address   string     `json:"address,omitempty"`
	BirthDate *time.Time `json:"birthDate,omitempty"`
	Email     string     `json:"email,omitempty"`
}

type AuthParameters struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
}

type LoginRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
}

func RegisterPhoneValidateParameters() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("phone", validation.PhoneNumberValidate)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

type LoginResponse struct {
	SessionId uuid.UUID `json:"sessionId"`
}

type ConfirmCodeRequest struct {
	SessionId uuid.UUID `json:"sessionId" validate:"required,uuid"`
	Code      string    `json:"code" validate:"required,len=4"`
}

type ConfirmCodeResponse struct {
	AuthParameters *AuthParameters `json:"authParameters"`
	User           *UserParameters `json:"user"`
}

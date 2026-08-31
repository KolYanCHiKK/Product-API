package user

import (
	"app/product-api/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	*Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

type UserConfirmParam struct {
	UserId    int        `json:"userId"`
	Surname   string     `json:"surname"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	City      string     `json:"city"`
	Address   string     `json:"address"`
	BirthDate *time.Time `json:"birthDate"`
	Email     string     `json:"email"`
}

func (s *Service) LoginUser(phone string) (*Session, int, error) {
	phone, err := utils.MapPhoneToStandardPattern(phone)
	if err != nil {
		return nil, 0, err
	}

	user, err := s.FindUserToRegisterByPhone(phone)
	if err != nil {
		return nil, 0, err
	}

	if user != nil {
		session, err := s.CreateSession(user.UserId)
		if err != nil {
			return nil, 0, err
		}

		return session, 200, nil
	}

	session, err := s.CreateUser(phone)
	if err != nil {
		return nil, 0, err
	}

	return session, 201, err

}

func (s *Service) ConfirmCode(sessionId uuid.UUID, code string) (*UserConfirmParam, error) {
	user, err := s.ConfirmUserCode(sessionId, code)
	if err != nil {
		return nil, err
	}

	return &UserConfirmParam{
		UserId:    user.UserId,
		Surname:   user.Surname,
		Name:      user.Name,
		Phone:     user.Phone,
		City:      user.City,
		Address:   user.Address,
		BirthDate: user.BirthDate,
		Email:     user.Email,
	}, nil
}

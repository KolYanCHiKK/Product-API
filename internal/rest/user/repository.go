package user

import (
	"app/product-api/pkg/db"
	"app/product-api/pkg/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *db.Db
}

func NewRepository(db2 *db.Db) *Repository {
	return &Repository{db: db2}
}

func (r *Repository) CreateUser(phone string) (session *Session, err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		user := User{
			Phone: phone,
		}
		result := tx.Create(&user)
		if result.Error != nil {
			return result.Error
		}

		code, err := utils.GenerateRandomCode(4)
		if err != nil {
			return err
		}
		session = &Session{
			UserId: user.UserId,
			Code:   code,
		}
		result = tx.Create(session)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *Repository) CreateSession(userId int) (*Session, error) {
	code, err := utils.GenerateRandomCode(4)
	if err != nil {
		return nil, err
	}
	session := Session{
		UserId: userId,
		Code:   code,
	}

	result := r.db.Create(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

func (r *Repository) FindUserToRegisterByPhone(phone string) (*User, error) {
	var user User
	result := r.db.
		Where("phone = ?", phone).
		First(&user)

	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *Repository) ConfirmUserCode(sessionId uuid.UUID, code string) (*User, error) {
	var user User
	result := r.db.
		Select("u.user_id").
		Table("sessions s").
		Joins("LEFT JOIN users u ON u.user_id = s.user_id").
		Where("s.session_id = ? AND s.code = ? AND s.is_confirmed = false", sessionId, code).
		First(&user)

	if result.RowsAffected == 0 {
		return nil, errors.New("invalid code")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		session := Session{SessionId: sessionId}
		resultTransact := tx.
			Model(&session).
			Update("is_confirmed", true)

		if resultTransact.Error != nil {
			return resultTransact.Error
		}

		resultTransact = tx.
			Model(&user).
			Clauses(clause.Returning{}).
			Update("is_deleted", false)

		if resultTransact.Error != nil {
			return resultTransact.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

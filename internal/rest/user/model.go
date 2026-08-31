package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId    int        `gorm:"not null;primaryKey;autoIncrement"`
	Surname   string     `gorm:"type:varchar(100);index:user_name_surname_idx"`
	Name      string     `gorm:"type:varchar(100);index:user_name_surname_idx"`
	Phone     string     `gorm:"type:varchar(65);unique;index:user_phone_idx"`
	City      string     `gorm:"type:varchar(100);index:user_city_idx"`
	Address   string     `gorm:"type:varchar(350)"`
	BirthDate *time.Time `gorm:"type:date"`
	Email     string     `gorm:"not null;type:varchar(120);index:user_email_idx"`
	Password  string     `gorm:"not null;type:varchar(1000)"`
	IsDeleted bool       `gorm:"not null;type:boolean;default:true"`
	CreatedAt time.Time  `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

type Session struct {
	SessionId   uuid.UUID `gorm:"not null;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserId      int       `gorm:"not null"`
	Code        string    `gorm:"not null;type:char(4);check:code_format,code ~ '^[0-9]{4}$'"`
	IsConfirmed bool      `gorm:"not null;type:boolean;default:false"`
	CreatedAt   time.Time `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

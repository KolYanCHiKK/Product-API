package order

import (
	"time"

	"gorm.io/datatypes"
)

type Order struct {
	OrderId      int            `json:"orderId" gorm:"primaryKey;index:idx_order_id_user_id,unique;autoIncrement" db:"order_id"`
	UserId       int            `json:"userId" gorm:"index:idx_order_id_user_id,unique;constraint:OnDelete:SET NULL;" db:"user_id"`
	Date         datatypes.Date `json:"date" gorm:"not null;default:current_date" db:"date"`
	IsPaid       bool           `json:"isPaid" gorm:"not null;default:false" db:"is_paid"`
	Descriptions string         `json:"descriptions,omitempty" gorm:"column:descriptions;type:varchar(2000)" db:"descriptions"`
	IsDeleted    bool           `gorm:"not null;type:boolean;default:false"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at;type:timestamptz;not null;default:now()" db:"created_at"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"column:updated_at;type:timestamptz;not null;default:now()" db:"updated_at"`

	OrderProducts []OrderProduct
}

type OrderProduct struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement" db:"id"`
	OrderId   int       `json:"orderId" gorm:"index:idx_order_id_product_id,unique;not null" db:"order_id"`
	ProductId int       `json:"productId" gorm:"index:idx_order_id_product_id,unique; not null" db:"product_id"`
	Price     float64   `json:"price" gorm:"column:price;type:decimal(50,2)" db:"price"`
	Quantity  int       `json:"quantity" gorm:"column:quantity;not null;default:0" db:"quantity"`
	IsDeleted bool      `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamptz;not null;default:now()" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;type:timestamptz;not null;default:now()" db:"updated_at"`
}

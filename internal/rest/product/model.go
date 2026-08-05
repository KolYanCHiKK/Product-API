package product

import (
	"time"

	"github.com/lib/pq"
)

type Product struct {
	ProductId    int            `json:"productId" gorm:"column:product_id;primaryKey;autoIncrement"`
	Name         string         `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Price        float64        `json:"price" gorm:"column:price;type:decimal(50,2);not null"`
	Quantity     int            `json:"quantity" gorm:"column:quantity;not null"`
	Descriptions string         `json:"descriptions,omitempty" gorm:"column:descriptions;type:varchar(2000)"`
	Images       pq.StringArray `json:"images,omitempty" gorm:"column:images;type:text[]"`
	CreateAt     time.Time      `json:"createAt" gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdateAt     time.Time      `json:"updateAt" gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

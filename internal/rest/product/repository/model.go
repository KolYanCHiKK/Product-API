package repository

import (
	"time"

	"github.com/lib/pq"
)

var ProductAllowedMap = map[string]string{
	"productId":    "product_id",
	"name":         "name",
	"price":        "price",
	"quantity":     "quantity",
	"descriptions": "descriptions",
	"images":       "images[1]",
	"createAt":     "created_at",
	"updateAt":     "updated_at",
}

type Product struct {
	ProductId    int            `json:"productId" gorm:"column:product_id;primaryKey;autoIncrement" db:"product_id"`
	Name         string         `json:"name" gorm:"column:name;type:varchar(255);not null" db:"name"`
	Price        float64        `json:"price" gorm:"column:price;type:decimal(50,2);not null" db:"price"`
	Quantity     int            `json:"quantity" gorm:"column:quantity;not null" db:"quantity"`
	Descriptions string         `json:"descriptions,omitempty" gorm:"column:descriptions;type:varchar(2000)" db:"descriptions"`
	Images       pq.StringArray `json:"images,omitempty" gorm:"column:images;type:text[]" db:"images"`
	CreateAt     time.Time      `json:"createAt" gorm:"column:created_at;type:timestamptz;not null;default:now()" db:"created_at"`
	UpdateAt     time.Time      `json:"updateAt" gorm:"column:updated_at;type:timestamptz;not null;default:now()" db:"updated_at"`
}

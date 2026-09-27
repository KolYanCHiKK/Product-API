package order

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderRepository struct {
	tx *gorm.DB
}

func NewOrderRepository(tx *gorm.DB) *OrderRepository {
	return &OrderRepository{
		tx: tx,
	}
}

type OrderProductRepository struct {
	tx *gorm.DB
}

func NewOrderProductRepository(tx *gorm.DB) *OrderProductRepository {
	return &OrderProductRepository{
		tx: tx,
	}
}

type MappedProductOrderParameters struct {
	Price float64 `gorm:"column:price"`
	Name  string  `gorm:"column:name"`
}

type MappedOrderProductWithName struct {
	Id           int            `gorm:"column:id"`
	OrderId      int            `gorm:"column:order_id"`
	UserId       int            `gorm:"column:user_id"`
	Date         datatypes.Date `gorm:"column:date"`
	IsPaid       bool           `gorm:"column:is_paid"`
	Descriptions string         `gorm:"column:descriptions"`
	Name         string         `gorm:"column:name"`
	ProductId    int            `gorm:"column:product_id"`
	Price        float64        `gorm:"column:price"`
	Quantity     int            `gorm:"column:quantity"`
}

func (r *OrderRepository) CreateOrder(ctx context.Context, userId int, descriptions string, isPaid bool) (*Order, error) {
	order := Order{
		UserId:       userId,
		IsPaid:       isPaid,
		Descriptions: descriptions,
	}
	result := r.tx.
		WithContext(ctx).
		Create(&order)

	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func (r *OrderProductRepository) CreateOrderProduct(ctx context.Context, orderId int, productId int, quantity int) (*OrderProduct, string, error) {
	orderProduct := OrderProduct{
		OrderId:   orderId,
		ProductId: productId,
		Quantity:  quantity,
	}

	var product MappedProductOrderParameters
	result := r.tx.
		WithContext(ctx).
		Table("products").
		Select("price, name").
		Where("product_id = ?", productId).
		Scan(&product)

	if result.Error != nil {
		return nil, "", result.Error
	}

	orderProduct.Price = product.Price
	result = r.tx.
		WithContext(ctx).
		Create(&orderProduct)

	if result.Error != nil {
		return nil, "", result.Error
	}

	return &orderProduct, product.Name, nil
}

func (r *OrderRepository) GetOrderById(ctx context.Context, orderId int) (*Order, error) {
	order := Order{
		OrderId: orderId,
	}

	result := r.tx.
		WithContext(ctx).
		Where("order_id = ? AND is_deleted = false", order.OrderId).
		First(&order)

	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func (r *OrderProductRepository) GetOrderProductByOrderId(ctx context.Context, orderId int) ([]MappedOrderProductWithName, error) {
	var orderProduct []MappedOrderProductWithName

	result := r.tx.
		WithContext(ctx).
		Debug().
		Table("order_products as op").
		Select("op.*, o.user_id, o.date, p.name").
		Joins("JOIN products p ON op.product_id = p.product_id").
		Joins("JOIN orders o ON op.order_id = o.order_id").
		Where("o.order_id = ? AND o.is_deleted = false AND op.is_deleted = false", orderId).
		Scan(&orderProduct)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return orderProduct, nil
}

func (r *OrderProductRepository) GetUserOrders(ctx context.Context, userId int) ([]MappedOrderProductWithName, error) {
	var orderProduct []MappedOrderProductWithName

	result := r.tx.
		WithContext(ctx).
		Debug().
		Table("order_products as op").
		Select("op.*, o.user_id, o.date, p.name").
		Joins("JOIN products p ON op.product_id = p.product_id").
		Joins("JOIN orders o ON op.order_id = o.order_id").
		Where("o.user_id = ? AND o.is_deleted = false AND op.is_deleted = false", userId).
		Scan(&orderProduct)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return orderProduct, nil
}

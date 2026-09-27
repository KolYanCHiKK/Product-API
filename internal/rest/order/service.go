package order

import (
	"app/product-api/pkg/middlewares"
	"app/product-api/pkg/uow"
	"app/product-api/pkg/utils"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	*uow.TransactionsManager
}

func NewService(tm *uow.TransactionsManager) *Service {
	return &Service{TransactionsManager: tm}
}

type MappedProductParameters struct {
	ProductId int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type MappedProductResponseParameters struct {
	*ProductResponseParameters
}

func (m *MappedProductResponseParameters) ToCommand() *ProductResponseParameters {
	return &ProductResponseParameters{
		ProductId: m.ProductId,
		Name:      m.Name,
		Price:     m.Price,
		Quantity:  m.Quantity,
	}
}

func (s *Service) CreateOrder(ctx context.Context, userId int, descriptions string, isPaid bool, products []MappedProductParameters) (*Order, []*MappedProductResponseParameters, error) {
	var createdOrder *Order
	var createdOrderProduct []*MappedProductResponseParameters

	err := s.TransactionsManager.Execute(ctx, func(tx *gorm.DB) error {
		orderRepo := NewOrderRepository(tx)
		orderProductRepo := NewOrderProductRepository(tx)

		order, err := orderRepo.CreateOrder(ctx, userId, descriptions, isPaid)
		if err != nil {
			return err
		}

		createdOrder = order

		for _, value := range products {
			orderProduct, name, err := orderProductRepo.CreateOrderProduct(ctx, order.OrderId, value.ProductId, value.Quantity)
			if err != nil {
				return err
			}

			createdOrderProduct = append(createdOrderProduct, &MappedProductResponseParameters{
				ProductResponseParameters: &ProductResponseParameters{
					ProductId: orderProduct.ProductId,
					Name:      name,
					Price:     orderProduct.Price,
					Quantity:  orderProduct.Quantity},
			})

		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return createdOrder, createdOrderProduct, nil
}

func (s *Service) GetOrderById(ctx context.Context, orderId int) (*GetOrderById, error) {
	var result *GetOrderById

	err := s.Execute(ctx, func(tx *gorm.DB) error {
		orderProductRepo := NewOrderProductRepository(tx)

		orderProduct, err := orderProductRepo.GetOrderProductByOrderId(ctx, orderId)
		if err != nil {
			return err
		}
		if orderProduct == nil {
			return nil
		}

		var products []ProductResponseParameters
		for _, value := range orderProduct {
			products = append(products, ProductResponseParameters{
				ProductId: value.ProductId,
				Name:      value.Name,
				Price:     value.Price,
				Quantity:  value.Quantity,
			})
		}

		result = &GetOrderById{
			OrderId:      orderProduct[0].OrderId,
			UserId:       orderProduct[0].UserId,
			Date:         utils.FormatDate(time.Time(orderProduct[0].Date)),
			IsPaid:       orderProduct[0].IsPaid,
			Descriptions: orderProduct[0].Descriptions,
			Products:     products,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil

}

func (s *Service) GetUserOrders(ctx context.Context) ([]GetUserOrdersResponse, error) {
	var result []GetUserOrdersResponse

	err := s.Execute(ctx, func(tx *gorm.DB) error {
		orderProductRepo := NewOrderProductRepository(tx)

		orderProduct, err := orderProductRepo.GetUserOrders(ctx, ctx.Value(middlewares.UserId).(int))
		fmt.Println(orderProduct)
		if err != nil {
			return err
		}
		if orderProduct == nil {
			return nil
		}

		i := 0
		for i < len(orderProduct) {
			targetOrderId := orderProduct[i].OrderId

			userOrderWithProducts := GetUserOrdersResponse{
				OrderId:      targetOrderId,
				Date:         utils.FormatDate(time.Time(orderProduct[i].Date)),
				IsPaid:       orderProduct[i].IsPaid,
				Descriptions: orderProduct[i].Descriptions,
			}

			var j int
			var products []ProductResponseParameters
			for j < len(orderProduct) {
				if orderProduct[j].OrderId == targetOrderId {
					products = append(products, ProductResponseParameters{
						ProductId: orderProduct[j].ProductId,
						Name:      orderProduct[j].Name,
						Price:     orderProduct[j].Price,
						Quantity:  orderProduct[j].Quantity,
					})

					orderProduct = append(orderProduct[:j], orderProduct[j+1:]...)
				} else {
					j++
				}
			}

			userOrderWithProducts.Products = products
			result = append(result, userOrderWithProducts)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil

}

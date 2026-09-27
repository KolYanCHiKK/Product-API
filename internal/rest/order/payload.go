package order

type ProductRequestParameters struct {
	ProductId int `json:"productId" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}

type ProductResponseParameters struct {
	ProductId int     `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type CreateOrderRequest struct {
	Description string                     `json:"description"`
	IsPaid      bool                       `json:"isPaid"`
	Products    []ProductRequestParameters `json:"products" validate:"required,min=1,dive"`
}

type CreateOrderResponse struct {
	OrderId      int    `json:"orderId"`
	UserId       int    `json:"userId"`
	Date         string `json:"date"`
	IsPaid       bool   `json:"isPaid"`
	Descriptions string `json:"descriptions"`
	Products     []*ProductResponseParameters
}

type GetOrderById struct {
	OrderId      int    `json:"orderId"`
	UserId       int    `json:"userId"`
	Date         string `json:"date"`
	IsPaid       bool   `json:"isPaid"`
	Descriptions string `json:"descriptions"`
	Products     []ProductResponseParameters
}

type GetUserOrdersResponse struct {
	OrderId      int    `json:"orderId"`
	Date         string `json:"date"`
	IsPaid       bool   `json:"isPaid"`
	Descriptions string `json:"descriptions"`
	Products     []ProductResponseParameters
}

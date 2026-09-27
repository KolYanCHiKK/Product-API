package order

import (
	"app/product-api/pkg/middlewares"
	"app/product-api/pkg/responce"
	"app/product-api/pkg/uow"
	"app/product-api/pkg/utils"
	"app/product-api/pkg/validation"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Deps *HandlerDeps
}

type HandlerDeps struct {
	*Service
	*uow.TransactionsManager
}

func NewHandler(mux *http.ServeMux, deps *HandlerDeps) {
	h := &Handler{Deps: deps}

	mux.HandleFunc("POST /orders", h.CreateOrder())
	mux.HandleFunc("GET /orders/{id}", h.GetOrderById())
	mux.HandleFunc("GET /my-orders", h.GetUserOrders())
}

func (h *Handler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CreateOrderRequest
		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		errs := validation.ValidateBody(
			body,
			validation.OrderCreateValidate,
			func() (*validator.Validate, error) {
				validate := validator.New()
				return validate, nil
			},
		)

		if errs != nil {
			responce.CreateErrResponse(w, 400, errs...)
			return
		}

		var products []MappedProductParameters
		for _, value := range body.Products {
			products = append(products, MappedProductParameters(value))
		}

		order, mappedResponse, err := h.Deps.Service.CreateOrder(
			r.Context(), r.Context().Value(middlewares.UserId).(int),
			body.Description, body.IsPaid, products,
		)

		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		var respProducts []*ProductResponseParameters
		for _, value := range mappedResponse {
			respProducts = append(respProducts, value.ToCommand())
		}

		resp := CreateOrderResponse{
			OrderId:      order.OrderId,
			UserId:       order.UserId,
			Date:         utils.FormatDate(time.Time(order.Date)),
			IsPaid:       order.IsPaid,
			Descriptions: order.Descriptions,
			Products:     respProducts,
		}

		err = responce.CreateResponse(w, 201, resp)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}
	}
}

func (h *Handler) GetOrderById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			responce.CreateErrResponse(w, 400, "Отсуствует идентификатор заказа")
			return
		}

		body, err := h.Deps.GetOrderById(r.Context(), orderId)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		if body == nil {
			responce.CreateErrResponse(w, 404, "Заказ не найден")
			return
		}

		err = responce.CreateResponse(w, 200, body)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}
	}
}

func (h *Handler) GetUserOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := h.Deps.GetUserOrders(r.Context())
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		if body == nil {
			responce.CreateErrResponse(w, 404, "Нет найденных заказов")
			return
		}

		err = responce.CreateResponse(w, 200, body)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}
	}
}

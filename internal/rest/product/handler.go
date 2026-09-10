package product

import (
	"app/product-api/pkg/responce"
	"app/product-api/pkg/validation"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *Service
}

func NewHandler(router *http.ServeMux, service *Service) {
	h := &Handler{Service: service}

	router.HandleFunc("GET /products", h.GetProducts())
	router.HandleFunc("GET /products/{id}", h.GetProductById())
	router.HandleFunc("POST /products", h.CreateProduct())
	router.HandleFunc("PUT /products/{id}", h.UpdateProduct())
	router.HandleFunc("PATCH /products/{id}", h.PatchProduct())
	router.HandleFunc("DELETE /products", h.DeleteProduct())
}

func (h *Handler) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		if query.Get("page") == "" || query.Get("pageSize") == "" {
			responce.CreateErrResponse(w, 400, "Missing pagination in request")
			return
		}
		page, err1 := strconv.Atoi(query.Get("page"))
		pageSize, err2 := strconv.Atoi(query.Get("pageSize"))
		if err1 != nil || err2 != nil {
			responce.CreateErrResponse(w, 400, "Parameters page and pageSize should be int")
			return
		}

		resp, err := h.Service.GetAllProduct(req.Context(), page, pageSize, query.Get("orderBy"))
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
		err = responce.CreateResponse(w, 200, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
	}
}

func (h *Handler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		productId := req.PathValue("id")
		resp, err := h.Service.GetProductById(req.Context(), productId)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
		if resp == nil {
			responce.CreateErrResponse(w, 404, "Resource is not found")
			return
		}

		err = responce.CreateResponse(w, 200, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
	}
}

func (h *Handler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body CreateProductRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		errs := validation.ValidateBody(
			body,
			validation.ProductRequestValidate,
			RegisterCreateProductRequestValidateParameters,
		)
		if errs != nil {
			responce.CreateErrResponse(w, 400, errs...)
			return
		}

		resp, err := h.Service.CreateProduct(req.Context(), &body)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}

		err = responce.CreateResponse(w, 201, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
	}
}

func (h *Handler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body UpdateProductRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		errs := validation.ValidateBody(
			body,
			validation.ProductRequestValidate,
			RegisterUpdateProductRequestValidateParameters,
		)
		if errs != nil {
			responce.CreateErrResponse(w, 400, errs...)
			return
		}
		productId, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			responce.CreateErrResponse(w, 400, "Id is not int")
			return
		}

		resp, err := h.Service.PutProduct(
			req.Context(),
			productId,
			body,
		)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
		if resp == nil {
			responce.CreateErrResponse(w, 404, "Product not found")
			return
		}

		err = responce.CreateResponse(w, 200, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
	}
}

func (h *Handler) PatchProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body PatchProductRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		errs := validation.ValidateBody(
			body,
			validation.ProductRequestValidate,
			RegisterPatchProductRequestValidateParameters,
		)
		if errs != nil {
			responce.CreateErrResponse(w, 400, errs...)
			return
		}
		productId, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			responce.CreateErrResponse(w, 400, "Id is not int")
			return
		}

		resp, err := h.Service.PatchProduct(
			req.Context(),
			map[string]any{
				"name":         body.Name,
				"price":        body.Price,
				"quantity":     body.Quantity,
				"descriptions": body.Descriptions,
				"images":       body.Images,
			},
			productId,
		)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
		if resp == nil {
			responce.CreateErrResponse(w, 404, "Product not found")
			return
		}

		err = responce.CreateResponse(w, 200, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
	}
}

func (h *Handler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body DeleteProductRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			responce.CreateErrResponse(w, 400, err.Error())
			return
		}

		errs := validation.ValidateBody(
			body,
			validation.ProductDeleteValidate,
			RegisterDeleteProductRequestValidateParameters,
		)
		if errs != nil {
			responce.CreateErrResponse(w, 400, errs...)
			return
		}

		resp, err := h.Service.DeleteProduct(
			req.Context(),
			map[string]any{
				"productId":    body.ProductId,
				"name":         body.Name,
				"price":        body.Price,
				"quantity":     body.Quantity,
				"descriptions": body.Descriptions,
				"images":       body.Images,
			},
		)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
			return
		}
		if resp == nil {
			responce.CreateErrResponse(w, 404, "Product not found")
			return
		}

		err = responce.CreateResponse(w, 200, resp)
		if err != nil {
			responce.CreateErrResponse(w, 200, err.Error())
			return
		}
	}
}

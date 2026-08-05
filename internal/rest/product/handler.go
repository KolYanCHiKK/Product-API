package product

import (
	"app/product-api/configs"
	"app/product-api/pkg/responce"
	"app/product-api/pkg/utils"
	"app/product-api/pkg/validation"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lib/pq"
)

type Handler struct {
	dep *HandlerDependent
}

type HandlerDependent struct {
	Config *configs.Config
}

func NewHandler(router *http.ServeMux, conf *configs.Config) {
	dep := HandlerDependent{conf}
	h := &Handler{dep: &dep}

	router.HandleFunc("GET /products", h.getProducts())
	router.HandleFunc("POST /products", h.CreateProduct())
	router.HandleFunc("PUT /products", h.UpdateProduct())
	router.HandleFunc("DELETE /products", h.DeleteProduct())
}

func (h *Handler) getProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := GetAllProductResponse{Products: []Product{
			{
				ProductId:    1,
				Name:         "Картофель",
				Price:        67.99,
				Quantity:     1500,
				Descriptions: "Картофель Greenteam",
				Images: pq.StringArray{
					"https://picsum.photos/800/600",
					"https://picsum.photos/1200/800",
				},
				CreateAt: time.Now(),
				UpdateAt: time.Now(),
			},
		}}
		err := responce.CreateResponse(w, 200, resp)
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

		resp := Product{
			ProductId:    1,
			Name:         body.Name,
			Price:        *body.Price,
			Quantity:     *body.Quantity,
			Descriptions: body.Descriptions,
			Images:       utils.MapProductImage(body.Images),
			CreateAt:     time.Now(),
			UpdateAt:     time.Now(),
		}
		err = responce.CreateResponse(w, 201, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
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

		resp := Product{
			ProductId:    1,
			Name:         body.Name,
			Price:        *body.Price,
			Quantity:     *body.Quantity,
			Descriptions: body.Descriptions,
			Images:       utils.MapProductImage(body.Images),
			CreateAt:     time.Now(),
			UpdateAt:     time.Now(),
		}
		err = responce.CreateResponse(w, 200, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
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

		var resp DeleteProductResponse
		if body.Name == nil && body.Descriptions == nil && body.Images == nil {
			resp.Success = false
			resp.Rows = nil
		} else {
			resp = DeleteProductResponse{
				Success: true,
				Rows: []Product{
					{
						ProductId:    1,
						Name:         "cffs",
						Price:        1,
						Quantity:     1,
						Descriptions: "fd",
						Images:       utils.MapProductImage("1"),
						CreateAt:     time.Now(),
						UpdateAt:     time.Now(),
					},
				},
			}
		}

		err = responce.CreateResponse(w, 200, resp)
		if err != nil {
			responce.CreateErrResponse(w, 500, err.Error())
		}
	}
}

package product

import (
	"app/product-api/internal/rest/product/repository"
	"app/product-api/pkg/validation"

	"github.com/go-playground/validator/v10"
)

type GetAllProductResponse struct {
	Products []repository.Product `json:"products"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
	Total    int                  `json:"total"`
}

type GetProductByIdResponse struct {
	repository.Product
}

type CreateProductRequest struct {
	Name         string   `json:"name" validate:"required,min=3"`
	Price        *float64 `json:"price" validate:"required,gt=0.04,scale=2"`
	Quantity     *int     `json:"quantity" validate:"required"`
	Descriptions string   `json:"descriptions" validate:"max=2000"`
	Images       string   `json:"images" validate:"max=2000"`
}

func RegisterCreateProductRequestValidateParameters() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("scale", validation.ScaleValidate)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

type CreateProductResponse struct {
	repository.Product
}

type UpdateProductRequest struct {
	Name         string   `json:"name" validate:"required,min=3"`
	Price        *float64 `json:"price" validate:"required,gt=0.04,scale=2"`
	Quantity     *int     `json:"quantity" validate:"required"`
	Descriptions string   `json:"descriptions" validate:"max=2000"`
	Images       string   `json:"images" validate:"max=2000"`
}

func RegisterUpdateProductRequestValidateParameters() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("scale", validation.ScaleValidate)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

type UpdateProductResponse struct {
	repository.Product
}

type PatchProductRequest struct {
	Name         string   `json:"name" validate:"omitempty,min=3"`
	Price        *float64 `json:"price" validate:"omitempty,gt=0.04,scale=2"`
	Quantity     *int     `json:"quantity" validate:"omitempty"`
	Descriptions string   `json:"descriptions" validate:"max=2000"`
	Images       string   `json:"images" validate:"max=2000"`
}

func RegisterPatchProductRequestValidateParameters() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("scale", validation.ScaleValidate)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

type PatchProductResponse struct {
	repository.Product
}

type DeleteProductRequest struct {
	ProductId    *int     `json:"productId"`
	Name         string   `json:"name" validate:"omitempty,min=3"`
	Price        *float64 `json:"price" validate:"omitempty,gt=0.04,scale=2"`
	Quantity     *int     `json:"quantity" validate:"omitempty"`
	Descriptions string   `json:"descriptions" validate:"max=2000"`
	Images       string   `json:"images" validate:"max=2000"`
}

func RegisterDeleteProductRequestValidateParameters() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("scale", validation.ScaleValidate)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

type DeleteProductResponse struct {
	Success     bool                 `json:"success"`
	TotalDelete int                  `json:"totalDelete"`
	Rows        []repository.Product `json:"rows,omitempty"`
}

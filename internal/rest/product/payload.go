package product

import (
	"app/product-api/pkg/validation"

	"github.com/go-playground/validator/v10"
)

type GetAllProductResponse struct {
	Products []Product `json:"products"`
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
	Product
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
	Product
}

type DeleteProductRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=3"`
	Descriptions *string `json:"descriptions" validate:"omitempty,max=2000"`
	Images       *string `json:"images" validate:"omitempty,max=2000"`
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
	Success bool      `json:"success"`
	Rows    []Product `json:"rows,omitempty"`
}

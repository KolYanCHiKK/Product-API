package repository

import (
	"app/product-api/pkg/db"
	"fmt"
)

type ProductRepository struct {
	*db.Db
}

func NewRepository(db *db.Db) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) Create(productParameters *Product) (*Product, error) {
	var product Product
	result := r.
		Raw(
			createProduct,
			productParameters.Name,
			productParameters.Price,
			productParameters.Quantity,
			productParameters.Descriptions,
			productParameters.Images,
		).
		Scan(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *ProductRepository) GetAllRows(limit int, offset int, orderBy string) ([]Product, error) {
	var products []Product

	result := r.
		Debug().
		Raw(
			fmt.Sprintf(getAllProducts, orderBy),
			limit,
			offset,
		).
		Scan(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *ProductRepository) CountAllRaws() (int, error) {
	var totalCount int
	result := r.
		Raw(
			countAllProductsRows,
		).
		Scan(
			&totalCount,
		)

	if result.Error != nil {
		return 0, result.Error
	}
	return totalCount, nil
}

func (r *ProductRepository) GetProductById(productId int) (*Product, error) {
	var product Product
	result := r.
		Raw(
			getProductByID,
			productId,
		).
		Scan(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &product, nil
}

func (r *ProductRepository) PutProduct(productId int, name string, price float64, quantity int, descriptions string, image string) (*Product, error) {
	var product Product
	result := r.
		Raw(
			updateProduct,
			name,
			price,
			quantity,
			descriptions,
			image,
			productId,
		).
		Scan(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &product, nil
}

func (r *ProductRepository) PatchProduct(productId int, setRow string, setParams []any) (*Product, error) {
	var product Product
	setParams = append(setParams, productId)
	result := r.
		Raw(
			fmt.Sprintf(patchProduct, setRow, len(setParams)),
			setParams...,
		).
		Scan(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &product, nil
}

func (r *ProductRepository) DeleteProduct(deleteRow string, deleteParams []any) ([]Product, *int, error) {
	var products []Product
	result := r.
		Raw(
			fmt.Sprintf(deleteProduct, deleteRow),
			deleteParams...,
		).
		Scan(&products)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	fmt.Println(result.RowsAffected)

	total := int(result.RowsAffected)
	if total == 0 {
		return nil, &total, nil
	}

	return products, &total, nil
}

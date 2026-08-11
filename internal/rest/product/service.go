package product

import (
	"app/product-api/internal/rest/product/repository"
	"app/product-api/pkg/utils"
	"strconv"
	"sync"
)

type Service struct {
	Repository *repository.ProductRepository
}

func NewService(repo *repository.ProductRepository) *Service {
	return &Service{repo}
}

func (s *Service) CreateProduct(req *CreateProductRequest) (*CreateProductResponse, error) {
	productRow := &repository.Product{
		Name:         req.Name,
		Price:        *req.Price,
		Quantity:     *req.Quantity,
		Descriptions: req.Descriptions,
		Images:       utils.MapProductImage(req.Images),
	}

	result, err := s.Repository.Create(productRow)
	if err != nil {
		return nil, err
	}
	return &CreateProductResponse{*result}, nil
}

func (s *Service) GetAllProduct(page int, pageSize int, orderBy string) (*GetAllProductResponse, error) {
	orderStr, err := utils.BuildOrderBy(orderBy, repository.ProductAllowedMap)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	errCh := make(chan error)

	totalCh := make(chan int)
	wg.Add(1)
	go func(totalCh chan<- int, errCh chan<- error) {
		defer wg.Done()
		total, err := s.Repository.CountAllRaws()
		if err != nil {
			errCh <- err
			return
		}
		totalCh <- total
	}(totalCh, errCh)

	productCh := make(chan []repository.Product)
	wg.Add(1)
	go func(productCh chan<- []repository.Product, errCh chan<- error) {
		defer wg.Done()
		rows, err := s.Repository.GetAllRows(
			pageSize,
			(page-1)*pageSize,
			orderStr,
		)
		if err != nil {
			errCh <- err
			return
		}
		productCh <- rows
	}(productCh, errCh)

	go func() {
		wg.Wait()
		close(totalCh)
		close(productCh)
		close(errCh)
	}()

	var rows []repository.Product
	var total int
	// Добавить при изучении темы "Контекст" context.WithCancel, чтобы после получения в канале
	// errCh значения, все остальные горутины завершались
	for received := 0; received < 2; {
		select {
		case err = <-errCh:
			return nil, err
		case total = <-totalCh:
			received++
		case rows = <-productCh:
			received++
		}
	}

	return &GetAllProductResponse{
		Products: rows,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *Service) GetProductById(id string) (*GetProductByIdResponse, error) {
	productId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	product, err := s.Repository.GetProductById(productId)
	if product == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &GetProductByIdResponse{*product}, nil
}

func (s *Service) PutProduct(productId int, req UpdateProductRequest) (*UpdateProductResponse, error) {
	product, err := s.Repository.PutProduct(
		productId,
		req.Name,
		*req.Price,
		*req.Quantity,
		req.Descriptions,
		req.Images,
	)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}
	return &UpdateProductResponse{*product}, nil
}

func (s *Service) PatchProduct(parametersMap map[string]any, productId int) (*PatchProductResponse, error) {
	setStr, setParam, err := utils.BuildSetParams(parametersMap, repository.ProductAllowedMap)
	if err != nil {
		return nil, err
	}

	product, err := s.Repository.PatchProduct(productId, *setStr, setParam)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	return &PatchProductResponse{*product}, nil
}

func (s *Service) DeleteProduct(parametersMap map[string]any) (*DeleteProductResponse, error) {
	whereStr, whereParam, err := utils.BuildWhereFilter(parametersMap, repository.ProductAllowedMap)
	if err != nil {
		return nil, err
	}

	products, totalRow, err := s.Repository.DeleteProduct(*whereStr, whereParam)
	if err != nil {
		return nil, err
	}
	if *totalRow == 0 {
		return nil, nil
	}

	return &DeleteProductResponse{
		Success:     true,
		TotalDelete: *totalRow,
		Rows:        products,
	}, nil
}

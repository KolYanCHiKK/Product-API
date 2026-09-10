package product

import (
	"app/product-api/internal/rest/product/repository"
	"app/product-api/pkg/utils"
	"context"
	"strconv"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	Repository *repository.ProductRepository
}

func NewService(repo *repository.ProductRepository) *Service {
	return &Service{repo}
}

func (s *Service) CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductResponse, error) {
	productRow := &repository.Product{
		Name:         req.Name,
		Price:        *req.Price,
		Quantity:     *req.Quantity,
		Descriptions: req.Descriptions,
		Images:       utils.MapProductImage(req.Images),
	}

	result, err := s.Repository.Create(ctx, productRow)
	if err != nil {
		return nil, err
	}
	return &CreateProductResponse{*result}, nil
}

func (s *Service) GetAllProduct(ctx context.Context, page int, pageSize int, orderBy string) (*GetAllProductResponse, error) {
	orderStr, err := utils.BuildOrderBy(orderBy, repository.ProductAllowedMap)
	if err != nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	totalCh := make(chan int)
	productCh := make(chan []repository.Product)
	errCh := make(chan error)

	g.Go(func() error {
		total, err := s.Repository.CountAllRaws(ctx)
		if err != nil {
			return err
		}

		totalCh <- total
		return nil
	})
	g.Go(func() error {
		rows, err := s.Repository.GetAllRows(
			ctx,
			pageSize,
			(page-1)*pageSize,
			orderStr,
		)
		if err != nil {
			return err
		}

		productCh <- rows
		return nil
	})
	go func() {
		if err := g.Wait(); err != nil {
			errCh <- err
		}
		close(errCh)
		close(totalCh)
		close(productCh)
	}()

	var rows []repository.Product
	var total int
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

func (s *Service) GetProductById(ctx context.Context, id string) (*GetProductByIdResponse, error) {
	productId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	product, err := s.Repository.GetProductById(ctx, productId)
	if product == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &GetProductByIdResponse{*product}, nil
}

func (s *Service) PutProduct(ctx context.Context, productId int, req UpdateProductRequest) (*UpdateProductResponse, error) {
	product, err := s.Repository.PutProduct(
		ctx,
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

func (s *Service) PatchProduct(ctx context.Context, parametersMap map[string]any, productId int) (*PatchProductResponse, error) {
	setStr, setParam, err := utils.BuildSetParams(parametersMap, repository.ProductAllowedMap)
	if err != nil {
		return nil, err
	}

	product, err := s.Repository.PatchProduct(ctx, productId, *setStr, setParam)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	return &PatchProductResponse{*product}, nil
}

func (s *Service) DeleteProduct(ctx context.Context, parametersMap map[string]any) (*DeleteProductResponse, error) {
	whereStr, whereParam, err := utils.BuildWhereFilter(parametersMap, repository.ProductAllowedMap)
	if err != nil {
		return nil, err
	}

	products, totalRow, err := s.Repository.DeleteProduct(ctx, *whereStr, whereParam)
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

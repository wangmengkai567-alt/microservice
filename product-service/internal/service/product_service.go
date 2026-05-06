package service

import (
	"errors"
	"fmt"
	"product-service/internal/model"
	"product-service/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(name, description string, price float64, stock int, imageUrl string) (uint, error) {
	if name == "" {
		return 0, errors.New("product name cannot be empty")
	}
	if price <= 0 {
		return 0, errors.New("price must be greater than 0")
	}
	if stock < 0 {
		return 0, errors.New("stock cannot be negative")
	}

	product := &model.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageUrl:    imageUrl,
	}

	if err := s.repo.Create(product); err != nil {
		return 0, err
	}

	return product.ID, nil
}

func (s *ProductService) GetProduct(id uint) (*model.Product, error) {
	if id == 0 {
		return nil, errors.New("invalid product id")
	}

	return s.repo.FindByID(id)
}

func (s *ProductService) UpdateProduct(id uint, name, description string, price float64, stock int, imageUrl string) error {
	if id == 0 {
		return errors.New("invalid product id")
	}

	product, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	if name != "" {
		product.Name = name
	}
	if description != "" {
		product.Description = description
	}
	if price > 0 {
		product.Price = price
	}
	if stock >= 0 {
		product.Stock = stock
	}
	if imageUrl != "" {
		product.ImageUrl = imageUrl
	}

	return s.repo.Update(product)
}

func (s *ProductService) AdjustStock(id uint, delta int) error {
	if id == 0 {
		return errors.New("invalid product id")
	}
	if delta == 0 {
		return nil
	}

	// 使用原子操作直接在数据库层面调整库存，避免并发竞态条件
	if err := s.repo.AdjustStockAtomic(id, delta); err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return errors.New("product not found or insufficient stock")
		}
		return fmt.Errorf("failed to adjust stock: %w", err)
	}

	return nil
}

func (s *ProductService) DeleteProduct(id uint) error {
	if id == 0 {
		return errors.New("invalid product id")
	}

	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	return s.repo.Delete(id)
}

func (s *ProductService) ListProducts(page, pageSize int) ([]model.Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.List(page, pageSize)
}

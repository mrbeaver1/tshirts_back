package service

import (
	entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
	domain "github.com/mrbeaver1/tshirts_back/internal/domain/repository"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() []*entity.Product {
	return s.repo.FindAll()
}

func (s *ProductService) GetProductByID(id uint64) *entity.Product {
	return s.repo.FindOneById(id)
}

func (s *ProductService) CreateProduct(product *entity.Product) {
	s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(product *entity.Product) {
	s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint64) {
	s.repo.Delete(id)
}

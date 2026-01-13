package service

import entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"

// ProductUseCaseInterface defines the interface for product use cases
// This allows for easier testing and decoupling
type ProductUseCaseInterface interface {
	GetAllProducts() []*entity.Product
	GetProductByID(id uint64) *entity.Product
	CreateProduct(product *entity.Product)
	UpdateProduct(product *entity.Product)
	DeleteProduct(id uint64)
}

package container

import (
	domain "github.com/mrbeaver1/t_shirt/internal/domain/repository"
	infrastructure "github.com/mrbeaver1/t_shirt/internal/infrastructure/repository"
)

type Container struct {
	productRepository domain.ProductRepository
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) GetProductRepository() domain.ProductRepository {
	if c.productRepository == nil {
		c.productRepository = infrastructure.NewInMemoryProductRepository()
	}
	return c.productRepository
}

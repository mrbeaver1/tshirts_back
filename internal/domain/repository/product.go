package domain

import (
	entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
)

type ProductRepository interface {
	FindAll() []*entity.Product
	FindOneById(id uint64) *entity.Product
	Create(p *entity.Product)
	Update(p *entity.Product)
	Delete(id uint64)
}

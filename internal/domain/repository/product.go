package domain

import (
	domainEntity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
)

type ProductRepository interface {
	FindAll() []*domainEntity.Product
	FindOneById(id uint64) *domainEntity.Product
	Create(p *domainEntity.Product)
	Update(p *domainEntity.Product)
	Delete(id uint64)
}

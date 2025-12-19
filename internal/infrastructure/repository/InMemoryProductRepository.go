package infrastructure

import (
	"sync"

	domainEntity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
)

type InMemoryProductRepository struct {
	products map[uint64]*domainEntity.Product
	mutex    sync.RWMutex
	nextId   uint64
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[uint64]*domainEntity.Product),
		nextId:   1,
	}
}

func (r *InMemoryProductRepository) FindAll() []*domainEntity.Product {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	products := make([]*domainEntity.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}
	return products
}

func (r *InMemoryProductRepository) FindOneById(id uint64) *domainEntity.Product {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if product, exists := r.products[id]; exists {
		return product
	}
	return nil
}

func (r *InMemoryProductRepository) Create(p *domainEntity.Product) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	p.Id = r.nextId
	r.nextId++
	r.products[p.Id] = p
}

func (r *InMemoryProductRepository) Update(p *domainEntity.Product) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[p.Id]; exists {
		r.products[p.Id] = p
	}
}

func (r *InMemoryProductRepository) Delete(id uint64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.products, id)
}

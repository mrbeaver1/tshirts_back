package persistence

import (
	"sync"

	entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
)

type InMemoryProductRepository struct {
	products map[uint64]*entity.Product
	mutex    sync.RWMutex
	nextId   uint64
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[uint64]*entity.Product),
		nextId:   1,
	}
}

func (r *InMemoryProductRepository) FindAll() []*entity.Product {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	products := make([]*entity.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}
	return products
}

func (r *InMemoryProductRepository) FindOneById(id uint64) *entity.Product {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if product, exists := r.products[id]; exists {
		return product
	}
	return nil
}

func (r *InMemoryProductRepository) Create(p *entity.Product) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	p.Id = r.nextId
	r.nextId++
	r.products[p.Id] = p
}

func (r *InMemoryProductRepository) Update(p *entity.Product) {
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

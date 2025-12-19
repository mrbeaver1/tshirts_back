package domain

import "time"

type Product struct {
	Id        uint64
	Name      string
	ColorId   uint64
	SizeId    uint64
	Quantity  uint64
	Available uint64
	ParentId  *uint64
	// IsParent  bool.      Подумать над реализацией данного флага
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// NewProduct creates a new Product with the specified parameters
// If parentId is nil, the product is considered a parent product
func NewProduct(name string, colorId, sizeId uint64, parentId *uint64) *Product {
	now := time.Now().UTC()

	// isParent := parentId == nil

	return &Product{
		Name:      name,
		ColorId:   colorId,
		SizeId:    sizeId,
		Quantity:  0,
		Available: 0,
		ParentId:  parentId,
		// IsParent:  isParent,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

// NewProductWithParentId creates a product with a specific parent ID
func NewProductWithParentId(name string, colorId, sizeId uint64, parentId uint64) *Product {
	return NewProduct(name, colorId, sizeId, &parentId)
}

// NewParentProduct creates a product without a parent (parent product)
func NewParentProduct(name string, colorId, sizeId uint64) *Product {
	return NewProduct(name, colorId, sizeId, nil)
}

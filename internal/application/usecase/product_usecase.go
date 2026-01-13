package usecase

import (
	entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
	"github.com/mrbeaver1/tshirts_back/internal/domain/service"
	service_impl "github.com/mrbeaver1/tshirts_back/internal/domain/service"
	service_interface "github.com/mrbeaver1/tshirts_back/internal/domain/service"
)

type ProductUseCase struct {
	productService *service_impl.ProductService
}

// Ensure ProductUseCase implements the ProductUseCaseInterface
var _ service_interface.ProductUseCaseInterface = (*ProductUseCase)(nil)

func NewProductUseCase(productService *service.ProductService) *ProductUseCase {
	return &ProductUseCase{productService: productService}
}

func (uc *ProductUseCase) GetAllProducts() []*entity.Product {
	return uc.productService.GetAllProducts()
}

func (uc *ProductUseCase) GetProductByID(id uint64) *entity.Product {
	return uc.productService.GetProductByID(id)
}

func (uc *ProductUseCase) CreateProduct(product *entity.Product) {
	uc.productService.CreateProduct(product)
}

func (uc *ProductUseCase) UpdateProduct(product *entity.Product) {
	uc.productService.UpdateProduct(product)
}

func (uc *ProductUseCase) DeleteProduct(id uint64) {
	uc.productService.DeleteProduct(id)
}

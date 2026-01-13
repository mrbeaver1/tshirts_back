package server

import (
	"context"
	"testing"

	entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
	service "github.com/mrbeaver1/tshirts_back/internal/domain/service"
	v1 "github.com/mrbeaver1/tshirts_back/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockProductUseCase is a mock implementation of ProductUseCase for testing
type MockProductUseCase struct {
	products map[uint64]*entity.Product
	nextID   uint64
}

func NewMockProductUseCase() *MockProductUseCase {
	return &MockProductUseCase{
		products: make(map[uint64]*entity.Product),
		nextID:   1,
	}
}

// Ensure MockProductUseCase implements the ProductUseCaseInterface
var _ service.ProductUseCaseInterface = (*MockProductUseCase)(nil)

func (m *MockProductUseCase) GetAllProducts() []*entity.Product {
	products := make([]*entity.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products
}

func (m *MockProductUseCase) GetProductByID(id uint64) *entity.Product {
	return m.products[id]
}

func (m *MockProductUseCase) CreateProduct(product *entity.Product) {
	product.Id = m.nextID
	m.products[m.nextID] = product
	m.nextID++
}

func (m *MockProductUseCase) UpdateProduct(product *entity.Product) {
	if _, exists := m.products[product.Id]; exists {
		m.products[product.Id] = product
	}
}

func (m *MockProductUseCase) DeleteProduct(id uint64) {
	delete(m.products, id)
}

func TestProductServer_GetProducts(t *testing.T) {
	mockUseCase := NewMockProductUseCase()

	// Add some test products
	mockUseCase.CreateProduct(&entity.Product{Name: "Test Product 1", ColorId: 1, SizeId: 1, Quantity: 10, Available: 5})
	mockUseCase.CreateProduct(&entity.Product{Name: "Test Product 2", ColorId: 2, SizeId: 2, Quantity: 20, Available: 10})

	server := NewProductServer(mockUseCase)

	req := &v1.GetProductsRequest{}
	resp, err := server.GetProducts(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if len(resp.Data.Products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(resp.Data.Products))
	}
}

func TestProductServer_GetProduct(t *testing.T) {
	mockUseCase := NewMockProductUseCase()
	testProduct := &entity.Product{Name: "Test Product", ColorId: 1, SizeId: 1, Quantity: 10, Available: 5}
	mockUseCase.CreateProduct(testProduct)

	server := NewProductServer(mockUseCase)

	req := &v1.GetProductRequest{Id: 1}
	resp, err := server.GetProduct(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if len(resp.Data.Products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(resp.Data.Products))
	}

	if resp.Data.Products[0].Name != "Test Product" {
		t.Errorf("Expected product name 'Test Product', got '%s'", resp.Data.Products[0].Name)
	}
}

func TestProductServer_GetProduct_NotFound(t *testing.T) {
	mockUseCase := NewMockProductUseCase()
	server := NewProductServer(mockUseCase)

	req := &v1.GetProductRequest{Id: 999} // Non-existent ID
	_, err := server.GetProduct(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound error, got %v", status.Code(err))
	}
}

func TestProductServer_CreateProduct(t *testing.T) {
	mockUseCase := NewMockProductUseCase()
	server := NewProductServer(mockUseCase)

	req := &v1.CreateProductRequest{
		Product: &v1.Product{
			Name:      "New Product",
			ColorId:   1,
			SizeId:    1,
			Quantity:  5,
			Available: 3,
		},
	}

	resp, err := server.CreateProduct(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if len(resp.Data.Products) != 1 {
		t.Errorf("Expected 1 product in response, got %d", len(resp.Data.Products))
	}

	if resp.Data.Products[0].Name != "New Product" {
		t.Errorf("Expected product name 'New Product', got '%s'", resp.Data.Products[0].Name)
	}
}

func TestProductServer_UpdateProduct(t *testing.T) {
	mockUseCase := NewMockProductUseCase()
	testProduct := &entity.Product{Name: "Original Product", ColorId: 1, SizeId: 1, Quantity: 10, Available: 5}
	mockUseCase.CreateProduct(testProduct)

	server := NewProductServer(mockUseCase)

	req := &v1.UpdateProductRequest{
		Product: &v1.Product{
			Id:        1,
			Name:      "Updated Product",
			ColorId:   2,
			SizeId:    2,
			Quantity:  15,
			Available: 8,
		},
	}

	resp, err := server.UpdateProduct(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Data.Products[0].Name != "Updated Product" {
		t.Errorf("Expected product name 'Updated Product', got '%s'", resp.Data.Products[0].Name)
	}
}

func TestProductServer_DeleteProduct(t *testing.T) {
	mockUseCase := NewMockProductUseCase()
	testProduct := &entity.Product{Name: "To Delete Product", ColorId: 1, SizeId: 1, Quantity: 10, Available: 5}
	mockUseCase.CreateProduct(testProduct)

	server := NewProductServer(mockUseCase)

	req := &v1.DeleteProductRequest{Id: 1}
	resp, err := server.DeleteProduct(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	// Verify the product was deleted
	remainingProducts := mockUseCase.GetAllProducts()
	if len(remainingProducts) != 0 {
		t.Errorf("Expected 0 products after deletion, got %d", len(remainingProducts))
	}
}

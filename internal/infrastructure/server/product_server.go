package server

import (
	"context"
	"time"

	entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
	usecase "github.com/mrbeaver1/tshirts_back/internal/domain/service"
	v1 "github.com/mrbeaver1/tshirts_back/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProductServer implements the ProductService gRPC service
type ProductServer struct {
	v1.UnimplementedProductServiceServer
	productUseCase usecase.ProductUseCaseInterface
}

// NewProductServer creates a new ProductServer
func NewProductServer(productUseCase usecase.ProductUseCaseInterface) *ProductServer {
	return &ProductServer{
		productUseCase: productUseCase,
	}
}

// GetProducts returns all products
func (s *ProductServer) GetProducts(ctx context.Context, req *v1.GetProductsRequest) (*v1.GetProductsResponse, error) {
	products := s.productUseCase.GetAllProducts()

	pbProducts := make([]*v1.Product, len(products))
	for i, p := range products {
		pbProducts[i] = &v1.Product{
			Id:        p.Id,
			Name:      p.Name,
			ColorId:   p.ColorId,
			SizeId:    p.SizeId,
			Quantity:  p.Quantity,
			Available: p.Available,
			ParentId:  p.ParentId,
			CreatedAt: convertTimeToTimestamp(p.CreatedAt),
			UpdatedAt: convertTimeToTimestamp(p.UpdatedAt),
			DeletedAt: convertTimeToTimestamp(p.DeletedAt),
		}
	}

	return &v1.GetProductsResponse{
		Data: &v1.Data{
			Products: pbProducts,
		},
	}, nil
}

// GetProduct returns a product by ID
func (s *ProductServer) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.GetProductResponse, error) {
	product := s.productUseCase.GetProductByID(req.Id)
	if product == nil {
		return nil, status.Errorf(codes.NotFound, "product not found with id: %d", req.Id)
	}

	pbProduct := &v1.Product{
		Id:        product.Id,
		Name:      product.Name,
		ColorId:   product.ColorId,
		SizeId:    product.SizeId,
		Quantity:  product.Quantity,
		Available: product.Available,
		ParentId:  product.ParentId,
		CreatedAt: convertTimeToTimestamp(product.CreatedAt),
		UpdatedAt: convertTimeToTimestamp(product.UpdatedAt),
		DeletedAt: convertTimeToTimestamp(product.DeletedAt),
	}

	return &v1.GetProductResponse{
		Data: &v1.Data{
			Products: []*v1.Product{pbProduct},
		},
	}, nil
}

// CreateProduct creates a new product
func (s *ProductServer) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.CreateProductResponse, error) {
	if req.Product == nil {
		return nil, status.Error(codes.InvalidArgument, "product is required")
	}

	product := &entity.Product{
		Name:      req.Product.Name,
		ColorId:   req.Product.ColorId,
		SizeId:    req.Product.SizeId,
		Quantity:  req.Product.Quantity,
		Available: req.Product.Available,
		ParentId:  req.Product.ParentId,
	}

	s.productUseCase.CreateProduct(product)

	pbProduct := &v1.Product{
		Id:        product.Id,
		Name:      product.Name,
		ColorId:   product.ColorId,
		SizeId:    product.SizeId,
		Quantity:  product.Quantity,
		Available: product.Available,
		ParentId:  product.ParentId,
		CreatedAt: convertTimeToTimestamp(product.CreatedAt),
		UpdatedAt: convertTimeToTimestamp(product.UpdatedAt),
		DeletedAt: convertTimeToTimestamp(product.DeletedAt),
	}

	return &v1.CreateProductResponse{
		Data: &v1.Data{
			Products: []*v1.Product{pbProduct},
		},
	}, nil
}

// UpdateProduct updates an existing product
func (s *ProductServer) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.UpdateProductResponse, error) {
	if req.Product == nil {
		return nil, status.Error(codes.InvalidArgument, "product is required")
	}

	product := &entity.Product{
		Id:        req.Product.Id,
		Name:      req.Product.Name,
		ColorId:   req.Product.ColorId,
		SizeId:    req.Product.SizeId,
		Quantity:  req.Product.Quantity,
		Available: req.Product.Available,
		ParentId:  req.Product.ParentId,
	}

	s.productUseCase.UpdateProduct(product)

	pbProduct := &v1.Product{
		Id:        product.Id,
		Name:      product.Name,
		ColorId:   product.ColorId,
		SizeId:    product.SizeId,
		Quantity:  product.Quantity,
		Available: product.Available,
		ParentId:  product.ParentId,
		CreatedAt: convertTimeToTimestamp(product.CreatedAt),
		UpdatedAt: convertTimeToTimestamp(product.UpdatedAt),
		DeletedAt: convertTimeToTimestamp(product.DeletedAt),
	}

	return &v1.UpdateProductResponse{
		Data: &v1.Data{
			Products: []*v1.Product{pbProduct},
		},
	}, nil
}

// DeleteProduct deletes a product
func (s *ProductServer) DeleteProduct(ctx context.Context, req *v1.DeleteProductRequest) (*v1.DeleteProductResponse, error) {
	s.productUseCase.DeleteProduct(req.Id)

	return &v1.DeleteProductResponse{
		Data: &v1.Data{
			Products: []*v1.Product{},
		},
	}, nil
}

// Helper function to convert *time.Time to *timestamppb.Timestamp
func convertTimeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t != nil {
		return timestamppb.New(*t)
	}
	return nil
}

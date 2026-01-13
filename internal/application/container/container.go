package container

import (
	usecase "github.com/mrbeaver1/tshirts_back/internal/application/usecase"
	config "github.com/mrbeaver1/tshirts_back/internal/config"
	domain_repo "github.com/mrbeaver1/tshirts_back/internal/domain/repository"
	service "github.com/mrbeaver1/tshirts_back/internal/domain/service"
	"github.com/mrbeaver1/tshirts_back/internal/infrastructure/database"
	persistence "github.com/mrbeaver1/tshirts_back/internal/infrastructure/persistence"
)

type Container struct {
	cfg            *config.Config
	productRepo    domain_repo.ProductRepository
	productService *service.ProductService
	productUseCase *usecase.ProductUseCase
}

func NewContainer() *Container {
	cfg := config.LoadConfig()

	return &Container{
		cfg: cfg,
	}
}

func (c *Container) GetProductUseCase() *usecase.ProductUseCase {
	if c.productUseCase == nil {
		c.productUseCase = usecase.NewProductUseCase(c.GetProductService())
	}
	return c.productUseCase
}

func (c *Container) GetProductService() *service.ProductService {
	if c.productService == nil {
		c.productService = service.NewProductService(c.GetProductRepository())
	}
	return c.productService
}

func (c *Container) GetProductRepository() domain_repo.ProductRepository {
	if c.productRepo == nil {
		// Можно выбрать реализацию репозитория в зависимости от конфигурации
		// Для примера используем InMemory реализацию
		pool, _ := database.ConnectToPool()
		c.productRepo = persistence.NewPostgresProductRepository(pool)

		// Или использовать PostgreSQL реализацию
		// db, err := ConnectToPool()
		// if err != nil {
		//     panic(err)
		// }
		// c.productRepo = persistence.NewPostgresProductRepository(db)
	}
	return c.productRepo
}

package main

import (
	"fmt"

	container "github.com/mrbeaver1/tshirts_back/internal/application/container"
	entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
	"github.com/mrbeaver1/tshirts_back/internal/infrastructure/database"
)

func main() {
	// Запускаем миграции
	err := database.RunMigrations()
	if err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
	}

	// Создаем контейнер зависимостей
	cont := container.NewContainer()

	// Получаем use case для работы с продуктами
	productUseCase := cont.GetProductUseCase()

	// Пример создания дочернего продукта с parentId=7
	childProduct := entity.NewProductWithParentId("Child product", 1, 2, 7)
	fmt.Println("Child product:", childProduct)

	// Пример создания родительского продукта без parentId
	parentProduct := entity.NewParentProduct("Parent product", 1, 2)
	fmt.Println("Parent product:", parentProduct)

	// Пример использования use case для создания продукта
	productUseCase.CreateProduct(parentProduct)
	fmt.Println("Created parent product with ID:", parentProduct.Id)

	// Пример получения всех продуктов
	allProducts := productUseCase.GetAllProducts()
	fmt.Printf("Total products: %d\n", len(allProducts))
}

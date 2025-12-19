package main

import (
	"fmt"

	domainEntity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
)

func main() {
	// Пример создания дочернего продукта с parentId=7
	childProduct := domainEntity.NewProductWithParentId("Child product", 1, 2, 7)
	fmt.Println("Child product:", childProduct)

	// Пример создания родительского продукта без parentId
	parentProduct := domainEntity.NewParentProduct("Parent product", 1, 2)
	fmt.Println("Parent product:", parentProduct)
}

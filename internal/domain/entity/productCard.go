package entity

import "time"

type ProductCardId uint64

type ProductCard struct {
	id             ProductCardId
	nomenclatureId NomenclatureId
	title          string
	description    string
	quantity       int64
	sku            string
	barcode        string
	sizeId         SizeId
	createdAt      time.Time
	updatedAt      time.Time
	deletedAt      time.Time
}

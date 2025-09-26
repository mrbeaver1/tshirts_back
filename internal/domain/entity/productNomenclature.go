package entity

import (
	"slices"
	"time"
)

type NomenclatureId uint64
type NomenclatureStatus uint32

const (
	Draft     NomenclatureStatus = 1
	OnReview  NomenclatureStatus = 2
	NeedFix   NomenclatureStatus = 3
	Active    NomenclatureStatus = 4
	InArchive NomenclatureStatus = 5
)

type ProductNomenclature struct {
	id           NomenclatureId
	title        string
	status       NomenclatureStatus
	productCards []ProductCard
	colorId      ColorId
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    time.Time
}

func InitNewNomenclature(title string) *ProductNomenclature {
	return &ProductNomenclature{
		title:     title,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		status:    Draft,
	}
}

func (p *ProductNomenclature) UpdateColor(colorId int) {
	p.colorId = ColorId(colorId)
	p.updatedAt = time.Now()
}

func (p *ProductNomenclature) AddProductCard(productCard ProductCard) {
	found := slices.ContainsFunc(p.productCards, func(pc ProductCard) bool {
		return productCard.id == pc.id
	})

	if found == false {
		p.productCards = append(p.productCards, productCard)
	}
}

func (p *ProductNomenclature) AddProductCards(productCards []ProductCard) {
	for _, productCard := range productCards {
		p.AddProductCard(productCard)
	}
}

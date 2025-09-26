package entity

import "time"

type SizeId uint64
type RussianSizeStandard uint64
type InternationalSizeStandard string
type SizeDimensions struct {
	width  uint64
	height uint64
}

type Size struct {
	id                    SizeId
	russianStandard       RussianSizeStandard
	internationalStandard InternationalSizeStandard
	dimensions            SizeDimensions
	createdAt             time.Time
	updatedAt             time.Time
	deletedAt             time.Time
}

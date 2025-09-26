package entity

import "time"

type ColorId uint64
type HexCode string
type RGB struct {
	red   uint32
	green uint32
	blue  uint32
}

type Color struct {
	id        ColorId
	name      string
	hex       HexCode
	rgb       RGB
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

package repositories

import "gorm.io/gorm"

type Frame struct {
	gorm.Model
	PlayerID string
	FrameNumber int32
	Roll1 string
	Roll2 string
	Roll3 string
	Score int32
}
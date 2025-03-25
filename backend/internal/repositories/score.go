package repositories

import "gorm.io/gorm"

type Score struct {
	gorm.Model
	Score int32
}
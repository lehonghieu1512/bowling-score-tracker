package repositories

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Game Game `gorm:"foreignKey:GameID"`
	GameID string 
	Name string
}
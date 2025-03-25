package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameBowlingRepository struct {
	db *gorm.DB
}

type Game struct {
	 // Game db model
	 gorm.Model
}

func NewGameBowlingRepo(db *gorm.DB) *GameBowlingRepository {
	return &GameBowlingRepository{
		db: db,
	}
}

func (repo *GameBowlingRepository) RegisterPlayers(c context.Context, playerNames []string) (id string, err error) {
	id = uuid.New().String()
	err = repo.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		id = uuid.New().String()
		game := Game{}
		if err := tx.Create(&game).Error; err != nil {
			return err
		}

		var players []Player
		for _, name := range playerNames {
			players = append(players, Player{
				Name: name,
			})
		}
		if err := tx.Create(&players).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("could not register players: %w", err)
	}
	return id, nil
}
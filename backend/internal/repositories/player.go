package repositories

import (
	"bowling-score-tracker/internal/services"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	Game Game `gorm:"foreignKey:GameID"`
	GameID uint 
	Name string
}

// PlayerRepository manages database operations for Player entities.
type PlayerRepository struct {
	db *gorm.DB
}

// NewPlayerRepo initializes a new PlayerRepository with the given database connection.
func NewPlayerRepo(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

// GetPlayersByGameIDs retrieves all players associated with the specified game IDs.
func (repo *PlayerRepository) GetPlayersByGameIDs(gameIDs []uint) ([]services.Player, error) {
	var players []Player
	if len(gameIDs) == 0 {
		return nil, nil
	}

	err := repo.db.Where("game_id IN ?", gameIDs).Find(&players).Error
	if err != nil {
		return nil, err
	}
	var outputPlayers []services.Player
	for _, player := range players {
		outputPlayers = append(outputPlayers, services.Player{
			ID: player.ID,
			Name: player.Name,
		})
	}
	return outputPlayers, nil
}
package repositories

import (
	"bowling-score-tracker/internal/services"

	"gorm.io/gorm"
)

type Frame struct {
	gorm.Model
	PlayerID uint
	FrameNumber int32
	Roll1 string
	Roll2 string
	Roll3 string
	Score int32
}

type FrameRepository struct {
	db *gorm.DB
}

func NewFrameRepo(db *gorm.DB) *FrameRepository {
	return &FrameRepository{
		db: db,
	}
}

// GetFramesByPlayerIDs retrieves all frames associated with the specified player IDs.
func (repo *FrameRepository) GetFramesByPlayerIDs(playerIDs []uint) ([]services.Frame, error) {
	var frames []Frame
	if len(playerIDs) == 0 {
		return nil, nil // Return an empty slice if no player IDs are provided
	}

	err := repo.db.Where("player_id IN ?", playerIDs).Find(&frames).Error
	if err != nil {
		return nil, err
	}

	var outputFrames []services.Frame
	for _, frame := range frames {
		outputFrames = append(outputFrames, 
			services.Frame{
				ID: frame.ID,
				Roll1: frame.Roll1,
				Roll2: frame.Roll2,
				Roll3: frame.Roll3,
				Score: frame.Score,
			},
		)
	}

	return outputFrames, nil
}
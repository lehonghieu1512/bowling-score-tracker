package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreateFramesInput struct {
	GameID string
	Frames map[string]PlayerFrameScore // map user id and their scores
}

type PlayerFrameScore struct {
	FrameNumber int     `json:"frame_number" validate:"required,min=1,max=10"`
	Roll1       *string `json:"roll1" validate:"required"`
	Roll2       *string `json:"roll2"`
	Roll3       *string `json:"roll3"`
	Score 		int32
}

type GameBowlingRepository struct {
	db *gorm.DB
}

type Game struct {
	 // Game db model
	 gorm.Model
	 CurrentFrame int32
	 PlayerNumber int32
}

func NewGameBowlingRepo(db *gorm.DB) *GameBowlingRepository {
	return &GameBowlingRepository{
		db: db,
	}
}

func (repo *GameBowlingRepository) RegisterPlayers(c context.Context, playerNames []string) (id string, err error) {
	id = uuid.New().String()
	err = repo.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		game := Game{
			CurrentFrame: 0,
			PlayerNumber: int32(len(playerNames)),
		}
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

func (repo *GameBowlingRepository) CreateFrames(c context.Context, input CreateFramesInput) (err error) {
	var game Game
	err = repo.db.WithContext(c).Model(&Game{}).Where("id = ?", input.GameID).Find(&game).Error
	if err != nil {
		return fmt.Errorf("could not find game: %w", err)
	}
	err = repo.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var frames []Frame 
		for playerID, score := range input.Frames {
			frames = append(frames, Frame{
				PlayerID: playerID,
				Roll1: lo.FromPtr(score.Roll1),
				Roll2: lo.FromPtr(score.Roll2),
				Roll3: lo.FromPtr(score.Roll3),
				FrameNumber: game.CurrentFrame,
				Score: score.Score,
			})
		}

		if err := tx.Create(&frames).Error; err != nil {
			return err
		}
		
		game.CurrentFrame += 1
		if err := tx.Save(&game).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not create frame: %w", err)
	}
	return nil
}
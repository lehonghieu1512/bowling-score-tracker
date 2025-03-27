package repositories

import (
	"bowling-score-tracker/internal/services"
	"context"
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreateFramesInput struct {
	GameID string
	Frames map[string]PlayerFrameScore // map user id and their scores
}

type PlayerFrameScore struct {
	FrameNumber int     
	Roll1       *string 
	Roll2       *string
	Roll3       *string
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

func (repo *GameBowlingRepository) RegisterPlayers(c context.Context, playerNames []string) (id uint, err error) {
	var game Game
	if len(playerNames) == 0 {
		return 0, fmt.Errorf("player name should not be empty")
	}
	err = repo.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		game = Game{
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
				GameID: game.ID,
			})
		}
		if err := tx.Create(&players).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("could not register players: %w", err)
	}
	return game.ID, nil
}

func (repo *GameBowlingRepository) CreateFrames(c context.Context, input services.CreateFrameInput) (err error) {
	var game Game

	err = repo.db.WithContext(c).Model(&Game{}).Where("id = ?", input.GameID).First(&game).Error
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
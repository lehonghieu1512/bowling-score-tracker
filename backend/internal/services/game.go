package services

import (
	"bowling-score-tracker/internal/repositories"
	"context"
	"fmt"
)

type CreateFrameInput struct {
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

type GameRepository interface {
	RegisterPlayers(c context.Context, playerNames []string) (id string, err error)
	CreateFrames(c context.Context, input repositories.CreateFramesInput) (err error)
}

type GameBowlingService struct {
	gameRepo GameRepository
}

func NewGameBowlingService(gameRepo GameRepository) *GameBowlingService {
	return &GameBowlingService{
		gameRepo: gameRepo,
	}
}

func (service *GameBowlingService) RegisterPlayers(c context.Context, playerNames []string) (id string, err error) {
	return service.gameRepo.RegisterPlayers(c, playerNames)
}

func (service *GameBowlingService) CreateFrame(c context.Context, input CreateFrameInput) error {
	frameScore := map[string]repositories.PlayerFrameScore{}
	for k, v := range input.Frames {
		frameScore[k] = repositories.PlayerFrameScore(v)
		frame := repositories.PlayerFrameScore(v)
		score, err := CalculateFrameScore(v)
		if err != nil {
			return fmt.Errorf("could not calculate score: %w", err)
		}
		print(score, " ")
		frame.Score = int32(score)
		frameScore[k] = frame
	}
	fmt.Printf("%v+", frameScore)
	return service.gameRepo.CreateFrames(c, repositories.CreateFramesInput{
		GameID: input.GameID,
		Frames: frameScore,
	})
}
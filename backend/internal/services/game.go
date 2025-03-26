package services

import (
	"context"
	"fmt"
)


type GameInfo struct {
	Players []Player
	Frames []Frame
}

type GameRepository interface {
	RegisterPlayers(c context.Context, playerNames []string) (id uint, err error)
	CreateFrames(c context.Context, input CreateFrameInput) (err error)
}



type GameBowlingService struct {
	gameRepo GameRepository
	playerRepo PlayerRepository
	frameRepo FrameRepo
}

func NewGameBowlingService(gameRepo GameRepository, playerRepo PlayerRepository, frameRepo FrameRepo) *GameBowlingService {
	return &GameBowlingService{
		gameRepo: gameRepo,
		playerRepo: playerRepo,
		frameRepo: frameRepo,
	}
}

func (service *GameBowlingService) RegisterPlayers(c context.Context, playerNames []string) (id uint, err error) {
	return service.gameRepo.RegisterPlayers(c, playerNames)
}

func (service *GameBowlingService) CreateFrame(c context.Context, input CreateFrameInput) error {
	frameScore := map[uint]PlayerFrameScore{}
	for k, v := range input.Frames {
		
		score, err := CalculateFrameScore(v)
		if err != nil {
			return fmt.Errorf("could not calculate score: %w", err)
		}
		v.Score = int32(score)
		frameScore[k] = v
	}
	fmt.Printf("%v+", frameScore)
	return service.gameRepo.CreateFrames(c, CreateFrameInput{
		GameID: input.GameID,
		Frames: frameScore,
	})
}

func (service *GameBowlingService) GetGameInfo(c context.Context, gameID uint) (gameinfo GameInfo, err error) {
	players, err := service.playerRepo.GetPlayersByGameIDs([]uint{gameID})
	if err != nil {
		return GameInfo{}, fmt.Errorf("could not get players: %w", err)
	}

	var playerIDs []uint
	for _, player := range players {
		playerIDs = append(playerIDs, player.ID)
	}

	frames, err := service.frameRepo.GetFramesByPlayerIDs(playerIDs)
	if err != nil {
		return GameInfo{}, fmt.Errorf("could not get frames: %w", err)
	}

	return GameInfo{
		Frames: frames,
		Players: players,
	}, nil
}
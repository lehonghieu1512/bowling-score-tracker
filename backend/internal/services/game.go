package services

import "context"

type GameRepository interface {
	RegisterPlayers(c context.Context, playerNames []string) (id string, err error)
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
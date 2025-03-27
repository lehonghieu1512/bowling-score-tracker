package services

import (
	"context"
	"errors"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) RegisterPlayers(c context.Context, playerNames []string) (uint, error) {
	args := m.Called(c, playerNames)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockGameRepository) CreateFrames(c context.Context, input CreateFrameInput) error {
	args := m.Called(c, input)
	return args.Error(0)
}

type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) GetPlayersByGameIDs(gameIDs []uint) ([]Player, error) {
	args := m.Called(gameIDs)
	return args.Get(0).([]Player), args.Error(1)
}

type MockFrameRepo struct {
	mock.Mock
}

func (m *MockFrameRepo) GetFramesByPlayerIDs(playerIDs []uint) ([]Frame, error) {
	args := m.Called(playerIDs)
	return args.Get(0).([]Frame), args.Error(1)
}

func TestRegisterPlayers(t *testing.T) {
	mockGameRepo := new(MockGameRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	mockFrameRepo := new(MockFrameRepo)

	service := NewGameBowlingService(mockGameRepo, mockPlayerRepo, mockFrameRepo)

	ctx := context.Background()
	playerNames := []string{"Alice", "Bob"}
	mockGameRepo.On("RegisterPlayers", ctx, playerNames).Return(uint(1), nil).Once()

	id, err := service.RegisterPlayers(ctx, playerNames)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)
	mockGameRepo.AssertExpectations(t)

	// Error case
	mockGameRepo.On("RegisterPlayers", ctx, playerNames).Return(uint(0), errors.New("registration failed")).Once()
	_, err = service.RegisterPlayers(ctx, playerNames)
	assert.Error(t, err)
}

func TestCreateFrame(t *testing.T) {
	mockGameRepo := new(MockGameRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	mockFrameRepo := new(MockFrameRepo)

	service := NewGameBowlingService(mockGameRepo, mockPlayerRepo, mockFrameRepo)

	ctx := context.Background()
	input := CreateFrameInput{
		GameID: 1,
		Frames: map[uint]PlayerFrameScore{
			1: {Roll1: lo.ToPtr("5"), Roll2: lo.ToPtr("4"), Score: 9},
		},
	}

	mockGameRepo.On("CreateFrames", ctx, mock.Anything).Return(nil).Once()
	
	err := service.CreateFrame(ctx, input)
	assert.NoError(t, err)
	mockGameRepo.AssertExpectations(t)

	// Error case - failed score calculation
	mockGameRepo.On("CreateFrames", ctx, mock.Anything).Return(errors.New("frame creation failed")).Once()
	err = service.CreateFrame(ctx, input)
	assert.Error(t, err)
}

func TestGetGameInfo(t *testing.T) {
	mockGameRepo := new(MockGameRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	mockFrameRepo := new(MockFrameRepo)

	service := NewGameBowlingService(mockGameRepo, mockPlayerRepo, mockFrameRepo)

	ctx := context.Background()
	gameID := uint(1)
	players := []Player{{ID: 1, Name: "Alice"}}
	frames := []Frame{{ID: 1, PlayerID: 1, Score: 9}}

	mockPlayerRepo.On("GetPlayersByGameIDs", []uint{gameID}).Return(players, nil).Once()
	mockFrameRepo.On("GetFramesByPlayerIDs", []uint{1}).Return(frames, nil).Once()

	gameInfo, err := service.GetGameInfo(ctx, gameID)

	assert.NoError(t, err)
	assert.Equal(t, players, gameInfo.Players)
	assert.Equal(t, frames, gameInfo.Frames)
	mockPlayerRepo.AssertExpectations(t)
	mockFrameRepo.AssertExpectations(t)

	// Error case - failed to get players
	mockPlayerRepo.On("GetPlayersByGameIDs", []uint{gameID}).Return([]Player(nil), errors.New("failed to fetch players")).Once()
	_, err = service.GetGameInfo(ctx, gameID)
	assert.Error(t, err)

	// Error case - failed to get frames
	mockPlayerRepo.On("GetPlayersByGameIDs", []uint{gameID}).Return(players, nil).Once()
	mockFrameRepo.On("GetFramesByPlayerIDs", []uint{1}).Return([]Frame(nil), errors.New("failed to fetch frames")).Once()
	_, err = service.GetGameInfo(ctx, gameID)
	assert.Error(t, err)
}

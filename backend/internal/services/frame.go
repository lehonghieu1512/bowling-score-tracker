package services

import (
	"errors"
	"strconv"

	"github.com/samber/lo"
)

type Frame struct {
	ID uint
	PlayerID string
	FrameNumber int32
	Roll1 string
	Roll2 string
	Roll3 string
	Score int32
}

type CreateFrameInput struct {
	GameID uint
	Frames map[uint]PlayerFrameScore // map user id and their scores
}

type PlayerFrameScore struct {
	FrameNumber int     `json:"frame_number" validate:"required,min=1,max=10"`
	Roll1       *string `json:"roll1" validate:"required"`
	Roll2       *string `json:"roll2"`
	Roll3       *string `json:"roll3"`
	Score 		int32 
}


type FrameRepo interface {
	GetFramesByPlayerIDs(playerIDs []uint) ([]Frame, error)
}

func CalculateFrameScore(frame PlayerFrameScore) (int, error) {
	if frame.Roll1 == nil {
		return 0, errors.New("Roll1 is required")
	}

	roll1Score, err := parseRoll(lo.FromPtr(frame.Roll1))
	if err != nil {
		return 0, err
	}

	roll2Score := 0
	if frame.Roll2 != nil {
		roll2Score, err = parseRoll(lo.FromPtr(frame.Roll2))
		if err != nil {
			return 0, err
		}
	}

	roll3Score := 0
	if frame.Roll3 != nil {
		roll3Score, err = parseRoll(lo.FromPtr(frame.Roll3))
		if err != nil {
			return 0, err
		}
	}

	// Calculate the frame score
	frameScore := roll1Score + roll2Score + roll3Score

	// Adjust for spares
	if frame.Roll2 != nil && *frame.Roll2 == "/" {
		frameScore = 10 // A spare always totals 10 for the frame
	}

	return frameScore, nil
}

// parseRoll converts a roll input into its numeric pinfall value.
func parseRoll(roll string) (int, error) {
	switch roll {
	case "X":
		return 10, nil // Strike
	case "/":
		return 0, nil  // Spare; actual value is calculated in context
	default:
		pins, err := strconv.Atoi(roll)
		if err != nil || pins < 0 || pins > 9 {
			return 0, errors.New("invalid roll value")
		}
		return pins, nil
	}
}
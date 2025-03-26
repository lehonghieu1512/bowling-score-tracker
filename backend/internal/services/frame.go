package services

import (
	"errors"
	"strconv"

	"github.com/samber/lo"
)

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
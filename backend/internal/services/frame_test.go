package services

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestCalculateFrameScore(t *testing.T) {
	tests := []struct {
		name    string
		frame   PlayerFrameScore
		expected int
		expectErr bool
	}{
		{
			name: "Success",
			frame: PlayerFrameScore{
				Roll1: lo.ToPtr("3"),
				Roll2: lo.ToPtr("5"),
				Score: 0,
			},
			expected: 8,
			expectErr: false,
		},
		{
			name: "Strike",
			frame: PlayerFrameScore{
				Roll1: lo.ToPtr("X"),
				Score: 0,
			},
			expected: 10,
			expectErr: false,
		},
		{
			name: "Spare",
			frame: PlayerFrameScore{
				Roll1: lo.ToPtr("5"),
				Roll2: lo.ToPtr("/"),
				Score: 0,
			},
			expected: 10,
			expectErr: false,
		},
		{
			name: "Invalid Roll1",
			frame: PlayerFrameScore{
				Roll1: lo.ToPtr("Z"),
			},
			expected: 0,
			expectErr: true,
		},
		{
			name: "Roll1 missing",
			frame: PlayerFrameScore{
				Roll1: nil,
			},
			expected: 0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateFrameScore(tt.frame)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

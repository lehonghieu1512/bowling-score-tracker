package repositories

import (
	"bowling-score-tracker/internal/services"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetPlayersByGameIDs(t *testing.T) {
    // Initialize the mock database and GORM DB using the reusable setup function
    gormDB, mock, cleanup := setupMockDB(t)
    defer func() {
        db, _ := gormDB.DB()
        db.Close()
		cleanup()
    }()

    // Initialize the repository with the mock GORM DB
    repo := NewPlayerRepo(gormDB)

    // Define test cases
    tests := []struct {
        name       string
        gameIDs    []uint
        mockSetup  func()
        wantResult []services.Player
        wantErr    bool
    }{
        {
            name:    "Valid game IDs",
            gameIDs: []uint{1,2},
            mockSetup: func() {
                rows := sqlmock.NewRows([]string{"id", "name", "game_id"}).
                    AddRow(1, "Player One", 1).
                    AddRow(2, "Player Two", 2)
                mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `players` WHERE game_id IN (?,?) AND `players`.`deleted_at` IS NULL")).
                    WithArgs(1, 2).
                    WillReturnRows(rows)
            },
            wantResult: []services.Player{
                {ID: 1, Name: "Player One"},
                {ID: 2, Name: "Player Two"},
            },
            wantErr: false,
        },
        {
            name:    "Empty game IDs",
            gameIDs: []uint{},
            mockSetup: func() {
                // No expectations since the method should return immediately
            },
            wantResult: nil,
            wantErr:    false,
        },
        {
            name:    "Non-existent game IDs",
            gameIDs: []uint{999},
            mockSetup: func() {
                rows := sqlmock.NewRows([]string{"id", "name", "game_id"})
                mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `players` WHERE game_id IN (?) AND `players`.`deleted_at` IS NULL")).
                    WithArgs(999).
                    WillReturnRows(rows)
            },
            wantResult: nil,
            wantErr:    false,
        },
    }

    // Execute test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            result, err := repo.GetPlayersByGameIDs(tt.gameIDs)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResult, result)
            }
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}
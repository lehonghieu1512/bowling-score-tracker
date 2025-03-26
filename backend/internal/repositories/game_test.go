package repositories

import (
	"bowling-score-tracker/internal/services"
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
    // Create a new sqlmock database connection
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
	mock.ExpectQuery(regexp.QuoteMeta("select sqlite_version()")).
		WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.35.5"))
    // Configure GORM to use the sqlmock database connection
    gormDB, err := gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
    if err != nil {
        t.Fatalf("Failed to open gorm DB: %v", err)
    }

    // Return the GORM DB, sqlmock, and a cleanup function
    return gormDB, mock, func() {
        db.Close()
    }
}

func TestRegisterPlayers(t *testing.T) {
    // Initialize the mock database and repository


    // Define test cases
    tests := []struct {
        name        string
        playerNames []string
        wantErr     bool
		wantDBErr bool
    }{
        {
            name:        "Valid player names",
            playerNames: []string{"Alice", "Bob"},
            wantErr:     false,
        },
        {
            name:        "No player names provided",
            playerNames: []string{},
            wantErr:     true,
        },
		{
            name:        "db error",
            playerNames: []string{"Alice", "Bob"},
            wantErr:     true,
			wantDBErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()
			repo := NewGameBowlingRepo(db)
			if tt.wantDBErr {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `games` (`created_at`,`updated_at`,`deleted_at`,`current_frame`,`player_number`) VALUES (?,?,?,?,?) RETURNING `id`")).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnError(fmt.Errorf("test"))
				mock.ExpectRollback()
			} else if len(tt.playerNames) > 0 {
				mock.ExpectBegin()
				// Mock the INSERT INTO "games" SQL statement
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `games` (`created_at`,`updated_at`,`deleted_at`,`current_frame`,`player_number`) VALUES (?,?,?,?,?) RETURNING `id`")).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
                // Mock the INSERT INTO "players" SQL statement
                mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `players` (`created_at`,`updated_at`,`deleted_at`,`game_id`,`name`) VALUES (?,?,?,?,?),(?,?,?,?,?) RETURNING `id`")).
                    WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, tt.playerNames[0],
                        sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, tt.playerNames[1]).
						WillReturnRows(sqlmock.NewRows([]string{}))
				// Commit the transaction
				mock.ExpectCommit()
            }
            _, err := repo.RegisterPlayers(context.Background(), tt.playerNames)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            // Ensure all expectations were met
            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("there were unfulfilled expectations: %s", err)
            }
        })
    }
}

func TestCreateFrames(t *testing.T) {
    db, mock, cleanup := setupMockDB(t)
    defer cleanup()
    repo := NewGameBowlingRepo(db)

    tests := []struct {
        name    string
        input   services.CreateFrameInput
        setup   func()
        wantErr bool
    }{
        {
            name: "Valid frames input",
            input: services.CreateFrameInput{
                GameID: 1,
                Frames: map[uint]services.PlayerFrameScore{
                    1: {FrameNumber: 1, Roll1: lo.ToPtr("X"), Score: 10},
                    2: {FrameNumber: 1, Roll1: lo.ToPtr("5"), Roll2: lo.ToPtr("/"), Score: 10},
                },
            },
            setup: func() {
                // Expect game retrieval
                mock.ExpectQuery(regexp.QuoteMeta(
                    "SELECT * FROM `games` WHERE id = ? AND `games`.`deleted_at` IS NULL ORDER BY `games`.`id` LIMIT 1")).
                    WithArgs(1).
                    WillReturnRows(sqlmock.NewRows([]string{"id", "current_frame"}).AddRow(1, 1))

				mock.ExpectBegin()

                // Expect frame inserts
                mock.ExpectQuery(regexp.QuoteMeta(
                    "INSERT INTO `frames` (`created_at`,`updated_at`,`deleted_at`,`player_id`,`frame_number`,`roll1`,`roll2`,`roll3`,`score`) VALUES (?,?,?,?,?,?,?,?,?),(?,?,?,?,?,?,?,?,?) RETURNING `id`")).
                    WithArgs(
                        sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 1,"X", "", "", 10,
                        sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 2, 1, "5", "/", "", 10,
                    ).
					WillReturnRows(sqlmock.NewRows([]string{}))

                // Expect game update
                mock.ExpectExec(regexp.QuoteMeta(
                    "UPDATE `games` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`current_frame`=?,`player_number`=? WHERE `games`.`deleted_at` IS NULL AND `id` = ?")).
                    WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
                    WillReturnResult(sqlmock.NewResult(1, 1))

                // Expect transaction commit
                mock.ExpectCommit()
            },
            wantErr: false,
        },
        {
            name: "Invalid game ID",
            input: services.CreateFrameInput{
                GameID: 1,
                Frames: map[uint]services.PlayerFrameScore{
                    1: {FrameNumber: 1, Roll1: lo.ToPtr("X"), Score: 10},
                },
            },
            setup: func() {
                // Expect game retrieval with no results
                mock.ExpectQuery(regexp.QuoteMeta(
                   "SELECT * FROM `games` WHERE id = ? AND `games`.`deleted_at` IS NULL ORDER BY `games`.`id` LIMIT 1")).
                    WithArgs(1).
                    WillReturnError(sql.ErrNoRows)
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            err := repo.CreateFrames(context.Background(), tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("unmet mock expectations: %s", err)
            }
        })
    }
}
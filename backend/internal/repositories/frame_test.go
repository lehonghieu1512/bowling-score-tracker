package repositories

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestGetFramesByPlayerIDs tests the GetFramesByPlayerIDs method of FrameRepository using table-driven tests.
func TestGetFramesByPlayerIDs(t *testing.T) {
	// Initialize sqlmock database connection and mock object
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta("select sqlite_version()")).
		WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.35.5"))
	// Open gorm DB connection using the sqlmock database with sqlite driver
	gormDB, err := gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{})
	assert.NoError(t, err)

	// Create an instance of FrameRepository with the mocked gorm DB
	repo := NewFrameRepo(gormDB)

	// Define test cases
	tests := []struct {
		name      string
		playerIDs []uint
		wantErr   bool
		wantLen   int
		prepareMockData func(sqlmock.Sqlmock)
	}{
		{
			name:      "Valid player IDs",
			playerIDs: []uint{1, 2},
			wantErr: false,
			wantLen: 2,
			prepareMockData:  func(sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `frames` WHERE player_id IN (?,?) AND `frames`.`deleted_at` IS NULL")).
				WithArgs(1, 2).
				WillReturnRows(sqlmock.NewRows([]string{"id", "player_id", "frame_number", "roll1", "roll2", "roll3", "score"}).
				AddRow(1, 1, 1, "X", "", "", 10).
				AddRow(2, 2, 1, "5", "/", "", 10))
			},
		},
		{
			name:      "No player IDs provided",
			playerIDs: []uint{},
			wantErr:   true,
			wantLen:   0,
			prepareMockData:  func(sqlmock.Sqlmock) {

			},
		},
		{
			name:      "Non-existent player IDs",
			playerIDs: []uint{3, 4},
			wantErr:   false,
			wantLen:   0,
			prepareMockData:  func(sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `frames` WHERE player_id IN (?,?) AND `frames`.`deleted_at` IS NULL")).
				WithArgs(3, 4).
				WillReturnRows(sqlmock.NewRows([]string{"id", "player_id", "frame_number", "roll1", "roll2", "roll3", "score"}),)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if len(tt.playerIDs) > 2 {
			// Expect the query to be executed with the specified player IDs
			// mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `frames` WHERE player_id IN (?,?) AND `frames`.`deleted_at` IS NULL")).
			// 	WithArgs(tt.playerIDs[0], tt.playerIDs[1]).
			// 	WillReturnRows(tt.mockRows)
			// }
			// Call the method under test
			tt.prepareMockData(mock)
			frames, err := repo.GetFramesByPlayerIDs(tt.playerIDs)

			// Assert error expectation
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Assert the length of returned frames
			assert.Len(t, frames, tt.wantLen)

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
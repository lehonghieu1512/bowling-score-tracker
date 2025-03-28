package main

import (
	"bowling-score-tracker/internal/controllers"
	"bowling-score-tracker/internal/repositories"
	"bowling-score-tracker/internal/services"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func dbMigration() (*gorm.DB, error) {
	DB, err := gorm.Open(sqlite.Open("bowling.db"), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %w", err)
	}

	// Auto-migrate schema
	err = DB.AutoMigrate(&repositories.Player{}, &repositories.Frame{}, &repositories.Game{})
	if err != nil {
		return nil, fmt.Errorf("could not run auto migration: %w", err)
	}
	log.Println("Database connected and migrated!")
	return DB, nil
}

func getEchoServer(
	gameController *controllers.GameController,
) *echo.Echo {
	e := echo.New()
	// TODO: add here
	e.POST("/games", gameController.CreateGame)
	e.POST("/games/:gameID/frames", gameController.CreateFrame)
	e.GET("/games/:gameID", gameController.GetGameInfo)
	return e
}

func main() {
	db, err := dbMigration()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	gameBowlingRepo := repositories.NewGameBowlingRepo(db)
	frameRepo := repositories.NewFrameRepo(db)
	playerRepo := repositories.NewPlayerRepo(db)
	gameBowlingService := services.NewGameBowlingService(gameBowlingRepo, playerRepo, frameRepo)
	gameBowlingController := controllers.NewGameController(gameBowlingService)

	e := getEchoServer(gameBowlingController)
	e.Start(":8080")
}
package controllers

import (
	"bowling-score-tracker/internal/services"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type CreateGameRequest struct {
	PlayerNames []string `json:"player_names"`
}

type CreateFrameRequest struct {
	GameID string                        `json:"game_id" param:"gameID" validate:"required,uuid"`
	Scores map[string]PlayerFrameScore   `json:"scores" validate:"required"`
}

type PlayerFrameScore struct {
	FrameNumber int     `json:"frame_number" validate:"required,min=1,max=10"`
	Roll1       *string `json:"roll1" validate:"required"`
	Roll2       *string `json:"roll2"`
	Roll3       *string `json:"roll3"`
}

type GameService interface {
	RegisterPlayers(c context.Context, playerNames []string) (id string, err error)
	CreateFrame(c context.Context, input services.CreateFrameInput) error
}

type GameController struct {
	gameService GameService
}

func NewGameController(gameService GameService) *GameController {
	return &GameController{
		gameService: gameService,
	}
}

func (controller *GameController) CreateGame(c echo.Context) error {
		// Parse request body
		var req CreateGameRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	
		// Validate player count
		if len(req.PlayerNames) == 0 || len(req.PlayerNames) > 5 {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Number of players must be between 1 and 5"})
		}

		id , err := controller.gameService.RegisterPlayers(c.Request().Context(), req.PlayerNames)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not start game"})
		}

		return c.JSON(http.StatusOK, echo.Map{"session_id": id})
}

func (controller *GameController) CreateFrame(c echo.Context) error {
			// Parse request body
			var req CreateFrameRequest
			fmt.Println("du ma")
			if err := c.Bind(&req); err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
			}
			frameScore := map[string]services.PlayerFrameScore{}
			for k, v := range req.Scores {
				frameScore[k] = services.PlayerFrameScore{
					FrameNumber: v.FrameNumber,
					Roll1: v.Roll1,
					Roll2: v.Roll2,
					Roll3: v.Roll3,
				}
			}
			err := controller.gameService.CreateFrame(c.Request().Context(), services.CreateFrameInput{
				Frames: frameScore,
			})
			if err != nil {
				log.Error(err)
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create frame"})
			}
			return c.JSON(http.StatusNoContent, nil)
}
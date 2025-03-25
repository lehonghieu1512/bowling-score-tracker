package controllers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type CreateGameRequest struct {
	PlayerNames []string `json:"player_names"`
}


type GameService interface {
	RegisterPlayers(c context.Context, playerNames []string) (id string, err error)
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
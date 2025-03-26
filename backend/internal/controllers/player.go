package controllers

type Player struct {
	ID uint `json:"id"`
	GameID string `json:"game_id"`
	Name string `json:"name"`
}
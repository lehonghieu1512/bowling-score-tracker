package services

type Player struct {
	ID uint
	GameID string 
	Name string
}

type PlayerRepository interface {
	GetPlayersByGameIDs(gameIDs []uint) ([]Player, error)
}
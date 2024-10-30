package internal

import (
	"github.com/google/uuid"
	"qwixx/internal/player"
)

type GameID string

type Administrator struct {
	lobbies map[GameID][]player.Player
	runners []GameRunner
}

func NewAdministrator() *Administrator {
	return &Administrator{lobbies: make(map[GameID][]player.Player)}
}

func (a *Administrator) CreateGame(host player.Player) GameID {
	randomGameID := GameID(uuid.New().String())
	a.lobbies[randomGameID] = []player.Player{host}
	return randomGameID
}

func (a *Administrator) JoinGame(gameID GameID, newPlayer player.Player) {
	a.lobbies[gameID] = append(a.lobbies[gameID], newPlayer)
}

func (a *Administrator) StartGame(gameID GameID) {
	// TODO does this take the players and then remove them from the map....without deleting them altogether. need pointers maybe?
	players := a.lobbies[gameID]
	delete(a.lobbies, gameID)
	runner := NewGameRunner(players)
	a.runners = append(a.runners, runner)
	go runner.RunGame()
}

package qwixx

import (
	"github.com/google/uuid"
	"math/rand"
	"qwixx/board"
	"qwixx/player"
)

type playerWithID struct {
	ID     player.PlayerID
	player player.Player
}

type GameRunner interface {
	RunGame()
}

type gameRunnerImpl struct {
	playOrder []playerWithID
	boards    map[player.PlayerID]board.Board
	locks     map[board.RowColor]bool
}

func NewGameRunner(players []player.Player) GameRunner {
	playOrder := establishPlayOrder(players)
	return &gameRunnerImpl{
		playOrder: playOrder,
		boards:    initializeBoards(playOrder),
		locks:     make(map[board.RowColor]bool),
	}
}

func (gr *gameRunnerImpl) RunGame() {
	gr.notifyPlayersOfPlayOrder()

	turnCount := 0
	for !gr.isGameWon() {
		if turnCount > 1000 { // guard against eternal games
			break
		}
		currentPlayer := gr.playOrder[turnCount%len(gr.playOrder)]
		gr.runSingleTurn(currentPlayer)
		turnCount++
	}
	gr.endGame()
}

// establishPlayOrder establishes the play order of a game comprised of the given list of players
func establishPlayOrder(players []player.Player) []playerWithID {
	playOrder := make([]playerWithID, 0, len(players))
	for _, p := range players {
		playOrder = append(playOrder, playerWithID{ID: player.PlayerID(uuid.New().String()), player: p})
	}
	rand.Shuffle(len(playOrder), func(i, j int) {
		playOrder[i], playOrder[j] = playOrder[j], playOrder[i]
	})
	return playOrder
}

func initializeBoards(playOrder []playerWithID) map[player.PlayerID]board.Board {
	boards := make(map[player.PlayerID]board.Board, len(playOrder))
	for _, p := range playOrder {
		boards[p.ID] = board.NewGameBoard()
	}
	return boards
}

func (gr *gameRunnerImpl) notifyPlayersOfPlayOrder() {
	orderIDs := make([]player.PlayerID, 0, len(gr.playOrder))
	for _, p := range gr.playOrder {
		orderIDs = append(orderIDs, p.ID)
	}
	for _, p := range gr.playOrder {
		p.player.InformOfPlayOrder(orderIDs)
	}
}

// isGameWon determines if the currently running game is won
// a game is won if either
// - two rows are locked
// - a player has taken four penalties
// TODO implement penalties
func (gr *gameRunnerImpl) isGameWon() bool {
	return len(gr.locks) >= 2
}

func (gr *gameRunnerImpl) runSingleTurn(currentPlayer playerWithID) {

}

func (gr *gameRunnerImpl) endGame() {
	// TODO broadcast end game
}

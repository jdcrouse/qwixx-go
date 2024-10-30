package player

import (
	"qwixx/actions"
	"qwixx/board"
)

type PlayerID string

type Player interface {
	GetName() string
	InformOfPlayOrder(playerNames []string)
	PromptActivePlayerTurn(playerBoard board.Board, diceRoll actions.DiceRoll) actions.ActivePlayerTurn
	PromptInactivePlayerTurn(playerBoard board.Board, diceRoll actions.DiceRoll) actions.InactivePlayerTurn
	InformSuccessfulTurn(updatedBoard board.Board)
	InformOfOpponentMove(playerID PlayerID, move actions.Move)
	InformRowLocked(color actions.RowColor)
	InformWin()
	InformLoss(winnerID PlayerID)
}

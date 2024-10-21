package player

import (
	"qwixx/actions"
	"qwixx/board"
)

type PlayerID string

type Player interface {
	GetName() string
	InformOfID(playerID PlayerID)
	InformOfPlayOrder(playerIDs []PlayerID)
	PromptActivePlayerMove(playerBoard board.Board, diceRoll actions.DiceRoll) (whiteDiceMove *actions.Move, colorDiceMove *actions.Move, takePenalty bool)
	PromptInactivePlayerMove(playerBoard board.Board, diceRoll actions.WhiteDiceRoll) actions.Move
	InformOfOpponentMove(playerID PlayerID, move actions.Move)
	InformRowLocked(color actions.RowColor)
	InformGameOver(winnerID PlayerID)
}

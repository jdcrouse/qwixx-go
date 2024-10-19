package player

import "qwixx/board"

type PlayerID string

type Player interface {
	InformOfID(playerID PlayerID)
	InformOfPlayOrder(playerIDs []PlayerID)
	PromptMove(diceRoll board.DiceRoll) board.Move
	InformOfOpponentMove(playerID PlayerID, move board.Move)
	InformRowLocked(color board.RowColor)
	InformGameOver(winnerID PlayerID)
}

// TODO need penalties!! new type of move?

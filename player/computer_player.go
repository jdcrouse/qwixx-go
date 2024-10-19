package player

import (
	"fmt"
	"qwixx/board"
)

var _ Player = ComputerPlayer{}

type ComputerPlayer struct {
	ID             PlayerID
	board          board.Board
	opponentBoards map[PlayerID]board.Board
}

func NewComputerPlayer() Player {
	return &ComputerPlayer{}
}

func (c ComputerPlayer) InformOfID(playerID PlayerID) {
	c.ID = playerID
}

func (c ComputerPlayer) InformOfPlayOrder(playerIDs []PlayerID) {
	fmt.Printf("play order: %v", playerIDs)
	for _, playerID := range playerIDs {
		c.opponentBoards[playerID] = board.NewGameBoard()
	}
}

func (c ComputerPlayer) InformOfOpponentMove(playerID PlayerID, move board.Move) {
	fmt.Printf("player %v made move %v", playerID, move)
	_, _ = c.opponentBoards[playerID].MakeMove(move)
}

func (c ComputerPlayer) InformRowLocked(color board.RowColor) {
	fmt.Printf("row %v was locked", color)
	c.board.LockRow(color)
}

func (c ComputerPlayer) InformGameOver(winnerID PlayerID) {
	if winnerID == c.ID {
		fmt.Printf("PARTYYYY YOU WON")
	} else {
		fmt.Printf("BOO YOU LOST :( %v won the game", winnerID)
	}
}

func (c ComputerPlayer) PromptMove(diceRoll board.DiceRoll) board.Move {
	possibleMoves := board.DeterminePossibleMoves(diceRoll)
	for _, move := range possibleMoves {
		if ok, _ := c.board.IsMoveValid(move); ok {
			return move
		}
	}
	return board.Move{} // TODO handle penalty when implemented
}

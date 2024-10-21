package player

import (
	"fmt"
	"qwixx/actions"
	"qwixx/board"
	"qwixx/rule_checker"
)

var _ Player = ComputerPlayer{}

type ComputerPlayer struct {
	name string
	ID   PlayerID
}

func NewComputerPlayer(name string) Player {
	return &ComputerPlayer{
		name: name,
	}
}

func (c ComputerPlayer) GetName() string {
	return c.name
}

func (c ComputerPlayer) InformOfID(playerID PlayerID) {
	c.ID = playerID
}

func (c ComputerPlayer) InformOfPlayOrder(playerIDs []PlayerID) {
}

func (c ComputerPlayer) InformOfOpponentMove(playerID PlayerID, move actions.Move) {
	fmt.Printf("player %v made move %v\n", playerID, move)
}

func (c ComputerPlayer) InformRowLocked(color actions.RowColor) {
	fmt.Printf("row %v was locked\n", color)
}

func (c ComputerPlayer) InformGameOver(winnerID PlayerID) {
	if winnerID == c.ID {
		fmt.Printf("PARTYYYY YOU WON\n")
	} else {
		fmt.Printf("BOO YOU LOST :( %v won the game\n", winnerID)
	}
}

func (c ComputerPlayer) PromptActivePlayerMove(playerBoard board.Board, diceRoll actions.DiceRoll) (whiteDiceMove *actions.Move, colorDiceMove *actions.Move, takePenalty bool) {
	// TODO prompt these with two separate methods instead of one?
	possibleWhiteDiceMoves := rule_checker.DeterminePossibleWhiteDiceMoves(diceRoll)
	for _, move := range possibleWhiteDiceMoves {
		if ok, _ := playerBoard.IsMoveValid(move); ok {
			whiteDiceMove = &move
			break
		}
	}
	possibleColorDiceMoves := rule_checker.DeterminePossibleWhiteDiceMoves(diceRoll)
	for _, move := range possibleColorDiceMoves {
		if ok, _ := playerBoard.IsMoveValid(move); ok {
			colorDiceMove = &move
			break
		}
	}

	if whiteDiceMove == nil && colorDiceMove == nil {
		return whiteDiceMove, colorDiceMove, true
	}
	return whiteDiceMove, colorDiceMove, false
}

func (c ComputerPlayer) PromptInactivePlayerMove(playerBoard board.Board, diceRoll actions.WhiteDiceRoll) actions.Move {
	return actions.Move{} // TODO handle
}

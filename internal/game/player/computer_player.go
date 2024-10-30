package player

import (
	"fmt"
	"qwixx/internal/game/actions"
	"qwixx/internal/game/board"
	"qwixx/internal/game/rule_checker"
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

func (c ComputerPlayer) InformOfPlayOrder(playerNames []string) {
	playOrder := ""
	for idx, name := range playerNames {
		playOrder += fmt.Sprintf("  %v: %v", idx+1, name)
	}
	fmt.Printf("play order is:\n%v\n", playOrder)
}

func (c ComputerPlayer) PromptActivePlayerTurn(
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
) actions.ActivePlayerTurn {
	return actions.ActivePlayerTurn{
		WhiteDiceMove: pickFirstValidWhiteDiceMove(playerBoard, diceRoll),
		ColorDiceMove: pickFirstValidColorDiceMove(playerBoard, diceRoll),
	}
}

func (c ComputerPlayer) PromptInactivePlayerTurn(
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
) actions.InactivePlayerTurn {
	return actions.InactivePlayerTurn{
		WhiteDiceMove: pickFirstValidWhiteDiceMove(playerBoard, diceRoll),
	}
}

func pickFirstValidWhiteDiceMove(
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
) *actions.Move {
	possibleWhiteDiceMoves := rule_checker.DeterminePossibleWhiteDiceMoves(diceRoll)
	for _, move := range possibleWhiteDiceMoves {
		if ok, _ := playerBoard.IsMoveValid(move); ok {
			return &move
		}
	}
	return nil
}

func pickFirstValidColorDiceMove(
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
) *actions.Move {
	possibleWhiteDiceMoves := rule_checker.DeterminePossibleColorDiceMoves(diceRoll)
	for _, move := range possibleWhiteDiceMoves {
		if ok, _ := playerBoard.IsMoveValid(move); ok {
			return &move
		}
	}
	return nil
}

func (c ComputerPlayer) InformSuccessfulTurn(updatedBoard board.Board) {

}

func (c ComputerPlayer) InformOfOpponentMove(playerID PlayerID, move actions.Move) {}

func (c ComputerPlayer) InformRowLocked(color actions.RowColor) {
	fmt.Printf("row %v was locked\n", color)
}

func (c ComputerPlayer) InformWin() {
	fmt.Printf("PARTYYYY YOU WON\n")
}
func (c ComputerPlayer) InformLoss(winnerID PlayerID) {
	fmt.Printf("BOO YOU LOST :( %v won the game\n", winnerID)
}

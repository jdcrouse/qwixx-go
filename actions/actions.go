package actions

import (
	"fmt"
	"math/rand"
)

type RowColor int

const (
	RowColorRed RowColor = iota
	RowColorYellow
	RowColorGreen
	RowColorBlue
)

func (m RowColor) String() string {
	switch m {
	case RowColorRed:
		return "Red"
	case RowColorYellow:
		return "Yellow"
	case RowColorGreen:
		return "Green"
	case RowColorBlue:
		return "Blue"
	default:
		return ""
	}
}

// a Move represents crossing off the square with the given number on the row with the given color
type Move struct {
	RowColor   RowColor
	CellNumber int
}

func (m Move) String() string {
	return fmt.Sprintf("(%v %v)", m.RowColor.String(), m.CellNumber)
}

func NewMove(rowColor RowColor, cellNumber int) Move {
	return Move{RowColor: rowColor, CellNumber: cellNumber}
}

// DiceRoll represents the roll of the six Qwixx dice, where two are white
// and the other four are one of each row color from the Qwixx board (red, yellow, green, blue)
type DiceRoll struct {
	WhiteDiceRoll
	ColorDiceRoll
}

type WhiteDiceRoll struct {
	White1 int
	White2 int
}

type ColorDiceRoll struct {
	Red    int
	Blue   int
	Green  int
	Yellow int
}

func RollQwixxDice() DiceRoll {
	return DiceRoll{
		WhiteDiceRoll: WhiteDiceRoll{
			White1: rand.Intn(6) + 1,
			White2: rand.Intn(6) + 1,
		},
		ColorDiceRoll: ColorDiceRoll{
			Red:    rand.Intn(6) + 1,
			Yellow: rand.Intn(6) + 1,
			Green:  rand.Intn(6) + 1,
			Blue:   rand.Intn(6) + 1,
		},
	}

}

// ActivePlayerTurn represents the turn of an active player, where they can cross off a cell both with the sum of
// the white dice and the sum of one white die with one color die
// If both moves are nil, a penalty is taken
type ActivePlayerTurn struct {
	WhiteDiceMove *Move
	ColorDiceMove *Move
}

func (apt ActivePlayerTurn) String() string {
	whiteDiceMoveString := "nil"
	if apt.WhiteDiceMove != nil {
		whiteDiceMoveString = apt.WhiteDiceMove.String()
	}
	colorDiceMoveString := "nil"
	if apt.ColorDiceMove != nil {
		colorDiceMoveString = apt.ColorDiceMove.String()
	}
	return fmt.Sprintf("[W: %v, C: %v", whiteDiceMoveString, colorDiceMoveString)
}

// InactivePlayerTurn represents the turn of an inactive player, where they can optionally cross off a cell
// with the sum of the white dice
// If the move is nil, nothing happens and no penalty is taken.
type InactivePlayerTurn struct {
	WhiteDiceMove *Move
}

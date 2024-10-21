package board

import (
	"fmt"
	"qwixx/actions"
)

// Board represents the Qwixx board for a single player.
// The Qwixx board consists of four rows of eleven cells. The rows have a color: Red, Yellow, Green, Blue from top down.
// Cells can be empty or crossed off.
// Each row's cells are numbered from 2 to 12:
// - Red and Yellow rows are numbers in ascending order from 2-12.
// - Green ad Blue rows are numbered in descending order from 12-2.
type Board interface {
	Print() string
	IsMoveValid(move actions.Move) (ok bool, reason string)
	MakeMove(move actions.Move) error
	LockRow(color actions.RowColor)
	CalculateScore() int
}

type boardImpl struct {
	redRow    Row
	yellowRow Row
	greenRow  Row
	blueRow   Row
	locks     map[actions.RowColor]bool
}

func NewGameBoard() Board {
	return &boardImpl{
		redRow:    NewRedRow(),
		yellowRow: NewYellowRow(),
		greenRow:  NewGreenRow(),
		blueRow:   NewBlueRow(),
	}
}

func (b *boardImpl) Print() string {
	var textRepresentation string

	textRepresentation += "Red: "
	textRepresentation += b.redRow.Print()
	textRepresentation += "\n"

	textRepresentation += "Yellow: "
	textRepresentation += b.yellowRow.Print()
	textRepresentation += "\n"

	textRepresentation += "Green: "
	textRepresentation += b.greenRow.Print()
	textRepresentation += "\n"

	textRepresentation += "Blue: "
	textRepresentation += b.blueRow.Print()

	return textRepresentation
}

func (b *boardImpl) IsMoveValid(move actions.Move) (ok bool, reason string) {
	switch move.RowColor {
	case actions.RowColorRed:
		return b.redRow.IsMoveValid(move.CellNumber)
	case actions.RowColorYellow:
		return b.yellowRow.IsMoveValid(move.CellNumber)
	case actions.RowColorGreen:
		return b.greenRow.IsMoveValid(move.CellNumber)
	case actions.RowColorBlue:
		return b.blueRow.IsMoveValid(move.CellNumber)
	default:
		return false, fmt.Sprintf("invalid move row color: %d", move.RowColor)
	}
}

func (b *boardImpl) MakeMove(move actions.Move) error {
	switch move.RowColor {
	case actions.RowColorRed:
		return b.redRow.MakeMove(move.CellNumber)
	case actions.RowColorYellow:
		return b.yellowRow.MakeMove(move.CellNumber)
	case actions.RowColorGreen:
		return b.greenRow.MakeMove(move.CellNumber)
	case actions.RowColorBlue:
		return b.blueRow.MakeMove(move.CellNumber)
	default:
		return fmt.Errorf("invalid move row color: %d", move.RowColor)
	}
}

func (b *boardImpl) LockRow(color actions.RowColor) {
	b.locks[color] = true
}

func (b *boardImpl) CalculateScore() int {
	return b.redRow.CalculateScore() + b.yellowRow.CalculateScore() + b.greenRow.CalculateScore() + b.blueRow.CalculateScore()
}

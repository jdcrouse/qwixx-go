package board

import (
	"fmt"
)

// Board represents the Qwixx board for a single player.
// The Qwixx board consists of four rows of eleven cells. The rows have a color: Red, Yellow, Green, Blue from top down.
// Cells can be empty or crossed off.
// Each row's cells are numbered from 2 to 12:
// - Red and Yellow rows are numbers in ascending order from 2-12.
// - Green ad Blue rows are numbered in descending order from 12-2.
type Board interface {
	Print() string
	MakeMove(move Move) (ok bool, _ error)
}

type boardImpl struct {
	redRow    Row
	yellowRow Row
	greenRow  Row
	blueRow   Row
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

func (b *boardImpl) MakeMove(move Move) (ok bool, _ error) {
	switch move.rowColor {
	case Red:
		return b.redRow.MakeMove(move.cellNumber)
	case Yellow:
		return b.yellowRow.MakeMove(move.cellNumber)
	case Green:
		return b.greenRow.MakeMove(move.cellNumber)
	case Blue:
		return b.blueRow.MakeMove(move.cellNumber)
	default:
		return false, fmt.Errorf("invalid move row color: %d", move.rowColor)
	}
}

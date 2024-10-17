package board

import (
	"fmt"
)

// Row represents a row of eleven cells on the Qwixx board that can either be free or crossed off.
//
// Red and Yellow rows are "ascending", their cells are numbered from 2-12 [2 3 4 5 6 7 8 9 10 11 12]
// Green and Blue rows are "descending", their cells are numbered from 12-2 [12 11 10 9 8 7 6 5 4 3 2]
//
// Cells can only be crossed off from left to right. To cross off a cell in a row, the cell must be unoccupied and there must be no crossed off cells to its right.
type Row interface {
	// Print prints a text representation of the row's cells
	Print() string

	// MakeMove attempts to cross off the given cell number in this row,
	// mutating the row with the new state if the move is valid and returning an error if the move is invalid
	MakeMove(cellNumber int) (ok bool, err error)

	// IsLocked determines if this row is locked.
	// A row is locked for all players when any player has crossed off the rightmost cell in their row of that color.
	// Further cells cannot be crossed off once a row is locked.
	IsLocked() bool
}

type rowType int

const (
	RowTypeAscending rowType = iota
	RowTypeDescending
)

type rowImpl struct {
	rowType rowType
	cells   []int
	locked  bool
}

func NewRedRow() Row {
	return newRedRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
func NewYellowRow() Row {
	return newYellowRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
func NewGreenRow() Row {
	return newGreenRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
func NewBlueRow() Row {
	return newBlueRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func (r *rowImpl) IsLocked() bool {
	return r.locked
}

func (r *rowImpl) Print() string {
	var textRepresentation string
	for idx, value := range r.cells {
		cellNumber, _ := indexToCellNumber(r.rowType, idx)
		textRepresentation += printCell(cellNumber, value)
		if idx < len(r.cells)-1 {
			textRepresentation += " "
		}
	}
	return textRepresentation
}

// cellNumberToIndex turns a cell number into the index of a slice containing the value of the row's cells
func indexToCellNumber(rowType rowType, index int) (int, error) {
	if index < 0 || index > 10 {
		return -1, fmt.Errorf("invalid index %v", index)
	}
	switch rowType {
	case RowTypeAscending:
		// for a row [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12], 0->2, 1->3, ..., 10->12
		return index + 2, nil
	case RowTypeDescending:
		// for a row [12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2], 0->12, 1->11, ..., 10->2
		return 12 - index, nil
	}
	return -1, fmt.Errorf("invalid row type %v", rowType)
}

// printCell prints the given cell with its number and value like [12|X]
// value should only be 0 or 1 and corresponds to the crossed-off state of the cell
func printCell(cellNumber int, value int) string {
	return fmt.Sprintf("[%v|%v]", cellNumber, valueAsText(value))
}

func valueAsText(value int) string {
	if value == 1 {
		return "X"
	}
	return " "
}

// MakeMove crosses off the given cell in this row, returning the new row
func (r *rowImpl) MakeMove(cellNumber int) (ok bool, _ error) {
	err := isMoveValid(r.cells, r.rowType, r.locked, cellNumber)
	if err != nil {
		return false, err
	}
	index, err := cellNumberToIndex(r.rowType, cellNumber)
	if err != nil {
		return false, err
	}
	r.cells[index] = 1
	return true, nil
}

// isMoveValid determines if the cell of the given number for the given row and row type can be crossed off.
// Cells can only be crossed off from left to right.
// To cross off a cell in a row, the cell must be unoccupied and there must be no crossed off cells to its right.
func isMoveValid(cells []int, rowType rowType, isLocked bool, cellNumber int) error {
	if isLocked {
		return fmt.Errorf("row is locked")
	}

	if cellNumber < 2 || cellNumber > 12 {
		return fmt.Errorf("cell number must be between 2 and 12")
	}
	index, err := cellNumberToIndex(rowType, cellNumber)
	if err != nil {
		return err
	}
	if cells[index] == 1 {
		return fmt.Errorf("cell %v is already crossed off", cellNumber)
	}
	for _, value := range cells[index:] {
		if value == 1 {
			return fmt.Errorf("cell %v is to the left of already crossed off cells", cellNumber)
		}
	}
	return nil
}

// cellNumberToIndex turns a cell number into the index of a slice containing the value of the row's cells
// for a row [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12], cell 10 would correspond to index 8
func cellNumberToIndex(rowType rowType, cellNumber int) (int, error) {
	if cellNumber < 2 || cellNumber > 12 {
		return -1, fmt.Errorf("invalid cell number %v", cellNumber)
	}
	switch rowType {
	case RowTypeAscending:
		// for a row [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12], 2->0, 3->1, ..., 12->10
		return cellNumber - 2, nil
	case RowTypeDescending:
		// for a row [12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2], 12->0, 11->1, ..., 2->10
		return 12 - cellNumber, nil
	}
	return -1, fmt.Errorf("invalid row type %v", rowType)
}

func newRedRowFromCells(cells []int) Row {
	return newRowFromCells(RowTypeAscending, cells)
}
func newYellowRowFromCells(cells []int) Row {
	return newRowFromCells(RowTypeAscending, cells)
}
func newGreenRowFromCells(cells []int) Row {
	return newRowFromCells(RowTypeDescending, cells)
}
func newBlueRowFromCells(cells []int) Row {
	return newRowFromCells(RowTypeDescending, cells)
}
func newRowFromCells(rowType rowType, cells []int) Row {
	return &rowImpl{rowType: rowType, cells: cells, locked: false}
}

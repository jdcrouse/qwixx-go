package board

import (
	"errors"
	"fmt"
	"strconv"
)

// scoreTable represents the scores a row receives based on the number of crossed off cells at the end of the game
// each index corresponds to the score if that many cells are crossed off
// index 0 means 0 crossed off cells, 1 means 1 crossed off,
// all the way up to 12 since locking a row gives an extra cell (the lock)
var scoreTable = []int{0, 1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 66, 78}

// Row represents a row of 11 cells on the Qwixx board that can either be free or crossed off.
//
// Red and Yellow rows are "ascending", their cells are numbered from 2-12 [2 3 4 5 6 7 8 9 10 11 12]
// Green and Blue rows are "descending", their cells are numbered from 12-2 [12 11 10 9 8 7 6 5 4 3 2]
//
// Cells can only be crossed off from left to right. To cross off a cell in a row, the cell must be unoccupied and there must be no crossed off cells to its right.
type Row interface {
	// Print prints a text representation of the row's cells
	Print() string

	// Copy copies this row into a new struct to prevent mutation of the original
	Copy() Row

	// IsMoveValid determines if the given move is valid for this row, returning the reason as an error if it was not valid
	IsMoveValid(cellNumber int) (ok bool, reason string)

	// MakeMove attempts to cross off the given cell number in this row,
	// mutating the row with the new state if the move is valid and returning an error if the move is invalid
	MakeMove(cellNumber int) error

	// IsCellMarked determines if the given cell number has been crossed off in this row
	IsCellMarked(cellNumber int) bool

	// IsLocked determines if this row is locked.
	// A row is locked for all players when any player has crossed off the rightmost cell in their row of that color.
	// Further cells cannot be crossed off once a row is locked.
	IsLocked() bool
	// TODO there is a difference between a row that is locked, and a row that WAS LOCKED ON THIS BOARD
	// the former just means the row cant be played on anymore, the latter influences the score of this row because you
	// get to cross off an extra cell for the lock

	// CalculateScore determines the score of this row based on the number of cells that are crossed off
	CalculateScore() int
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

func (r *rowImpl) Print() string {
	var textRepresentation string
	for idx, value := range r.cells {
		cellNumber, _ := indexToCellNumber(r.rowType, idx)
		textRepresentation += printCell(cellNumber, value)
		if idx < len(r.cells)-1 {
			textRepresentation += " "
		} else {
			textRepresentation += printLockCell(value)
		}
	}
	return textRepresentation
}

func (r *rowImpl) Copy() Row {
	newCells := make([]int, len(r.cells))
	copy(newCells, r.cells)
	return &rowImpl{
		rowType: r.rowType,
		cells:   newCells,
		locked:  r.locked,
	}
}

func (r *rowImpl) IsMoveValid(cellNumber int) (ok bool, reason string) {
	return isMoveValid(r.cells, r.rowType, r.locked, cellNumber)
}

// MakeMove crosses off the given cell in this row, returning the new row
func (r *rowImpl) MakeMove(cellNumber int) error {
	ok, reason := isMoveValid(r.cells, r.rowType, r.locked, cellNumber)
	if !ok {
		return errors.New(reason)
	}
	index, err := cellNumberToIndex(r.rowType, cellNumber)
	if err != nil {
		return err
	}
	r.cells[index] = 1
	return nil
}

func (r *rowImpl) IsCellMarked(cellNumber int) bool {
	index, err := cellNumberToIndex(r.rowType, cellNumber)
	if err != nil {
		return false
	}
	return r.cells[index] == 1
}

func (r *rowImpl) IsLocked() bool {
	return r.locked
}

func (r *rowImpl) CalculateScore() int {
	// TODO include locked row? probably should add a twelfth cell
	crossOffCellCount := 0
	for _, value := range r.cells {
		if value == 1 {
			crossOffCellCount++
		}
	}
	// if the last cell is crossed off, that means this row was locked by this player so they also get to cross off the lock cell
	if r.cells[10] == 1 {
		crossOffCellCount++
	}
	return scoreTable[crossOffCellCount]
}

// isMoveValid determines if the cell of the given number for the given row and row type can be crossed off.
// Cells can only be crossed off from left to right.
// To cross off a cell in a row, the cell must be unoccupied and there must be no crossed off cells to its right.
func isMoveValid(cells []int, rowType rowType, isLocked bool, cellNumber int) (ok bool, reason string) {
	if isLocked {
		return false, "row is locked"
	}

	moveIndex, err := cellNumberToIndex(rowType, cellNumber)
	if err != nil {
		return false, err.Error()
	}

	// cell cannot be crossed off if it is already crossed off

	if cells[moveIndex] == 1 {
		return false, fmt.Sprintf("cell %v is already crossed off", cellNumber)
	}

	countCrossedOff := 0
	countCrossedOffToRightOfIndex := 0
	for idx, value := range cells {
		if value == 1 {
			countCrossedOff++
			if idx > moveIndex {
				countCrossedOffToRightOfIndex++
			}
		}
	}

	// cell cannot be crossed off if there are crossed off cells to its right
	if countCrossedOffToRightOfIndex > 0 {
		return false, fmt.Sprintf("cell %v is to the left of already crossed off cells", cellNumber)
	}

	// 5 other cells in row must be crossed off in order to cross off rightmost cell
	if moveIndex == 10 && countCrossedOff < 5 {
		return false, "cannot cross off rightmost cell of row unless 5 cells have been crossed off in that row"
	}

	return true, ""
}

// cellNumberToIndex turns a cell number into the index of a slice containing the value of the row's cells
// for ascending row [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12], 0->2, 1->3, ..., 10->12
// for descending row [12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2], 0->12, 1->11, ..., 10->2
func indexToCellNumber(rowType rowType, index int) (int, error) {
	if index < 0 || index > 10 {
		return -1, fmt.Errorf("invalid index: %v. must be between 0 and 10", index)
	}
	switch rowType {
	case RowTypeAscending:
		return index + 2, nil
	case RowTypeDescending:
		return 12 - index, nil
	}
	return -1, fmt.Errorf("invalid row type %v", rowType)
}

// cellNumberToIndex turns a cell number into the index of a slice containing the value of the row's cells
// for a row [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12], 2->0, 3->1, ..., 12->10
// for a row [12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2], 12->0, 11->1, ..., 2->10
func cellNumberToIndex(rowType rowType, cellNumber int) (int, error) {
	if cellNumber < 2 || cellNumber > 12 {
		return -1, fmt.Errorf("invalid cell number: %v. must be between 2 and 12", cellNumber)
	}
	switch rowType {
	case RowTypeAscending:
		return cellNumber - 2, nil
	case RowTypeDescending:
		return 12 - cellNumber, nil
	}
	return -1, fmt.Errorf("invalid row type %v", rowType)
}

// printCell prints the given cell with its number and value like [12|X]
// value should only be 0 or 1 and corresponds to the crossed-off state of the cell
func printCell(cellNumber int, value int) string {
	cellNumberText := strconv.Itoa(cellNumber)
	return fmt.Sprintf("[%v|%v]", cellNumberText, valueAsText(value))
}

func printLockCell(value int) string {
	return fmt.Sprintf(" [L|%v]", valueAsText(value))
}

func valueAsText(value int) string {
	if value == 1 {
		return "X"
	}
	return " "
}

func NewRedRow() Row {
	return newRedRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false)
}
func NewYellowRow() Row {
	return newYellowRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false)
}
func NewGreenRow() Row {
	return newGreenRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false)
}
func NewBlueRow() Row {
	return newBlueRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false)
}

func newRedRowFromCells(cells []int, locked bool) Row {
	return newAscendingRowFromCells(cells, locked)
}
func newYellowRowFromCells(cells []int, locked bool) Row {
	return newAscendingRowFromCells(cells, locked)
}
func newGreenRowFromCells(cells []int, locked bool) Row {
	return newDescendingRowFromCells(cells, locked)
}
func newBlueRowFromCells(cells []int, locked bool) Row {
	return newDescendingRowFromCells(cells, locked)
}
func newAscendingRowFromCells(cells []int, locked bool) Row {
	return newRowFromCells(RowTypeAscending, cells, locked)
}
func newDescendingRowFromCells(cells []int, locked bool) Row {
	return newRowFromCells(RowTypeDescending, cells, locked)
}
func newRowFromCells(rowType rowType, cells []int, locked bool) Row {
	return &rowImpl{rowType: rowType, cells: cells, locked: locked}
}

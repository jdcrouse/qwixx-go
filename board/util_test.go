package board

import (
	"errors"
	"qwixx/actions"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyActivePlayerTurn(t *testing.T) {
	type testCase struct {
		name           string
		inputBoard     Board
		inputTurn      actions.ActivePlayerTurn
		expectedBoard  Board
		expectedError  error
		expectedReason string
	}

	testCases := []testCase{
		{
			name:       "Apply white dice move only",
			inputBoard: NewGameBoard(),
			inputTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{RowColor: actions.RowColorRed, CellNumber: 5},
			},
			expectedBoard: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			},
		},
		{
			name:       "Apply color dice move only",
			inputBoard: NewGameBoard(),
			inputTurn: actions.ActivePlayerTurn{
				ColorDiceMove: &actions.Move{RowColor: actions.RowColorBlue, CellNumber: 8},
			},
			expectedBoard: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}, false),
			},
		},
		{
			name:       "Apply both white and color dice moves",
			inputBoard: NewGameBoard(),
			inputTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{RowColor: actions.RowColorYellow, CellNumber: 7},
				ColorDiceMove: &actions.Move{RowColor: actions.RowColorGreen, CellNumber: 10},
			},
			expectedBoard: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			},
		},
		{
			name: "Error on invalid white dice move",
			inputBoard: &boardImpl{
				redRow:    newRedRowFromCells([]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			},
			inputTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{RowColor: actions.RowColorRed, CellNumber: 2},
			},
			expectedBoard: &boardImpl{
				redRow:    newRedRowFromCells([]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			},
			expectedError:  errors.New("invalid white dice move"),
			expectedReason: "white dice move is already marked",
		},
		{
			name: "Error on invalid color dice move",
			inputBoard: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			},
			inputTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{RowColor: actions.RowColorYellow, CellNumber: 7},
				ColorDiceMove: &actions.Move{RowColor: actions.RowColorBlue, CellNumber: 8},
			},
			expectedError:  errors.New("invalid color dice move"),
			expectedReason: "color dice move is already marked",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updatedBoard, err := ApplyActivePlayerTurn(tc.inputBoard, tc.inputTurn)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedReason, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedBoard, updatedBoard)
			}
		})
	}
}

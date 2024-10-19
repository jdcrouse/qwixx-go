package board

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoardImpl_Print(t *testing.T) {
	type testCase struct {
		name                         string
		input                        Board
		expectedStringRepresentation string
	}
	testCases := []testCase{
		{
			name:  "base case, default starting board",
			input: NewGameBoard(),
			expectedStringRepresentation: `Red: [2| ] [3| ] [4| ] [5| ] [6| ] [7| ] [8| ] [9| ] [10| ] [11| ] [12| ]
Yellow: [2| ] [3| ] [4| ] [5| ] [6| ] [7| ] [8| ] [9| ] [10| ] [11| ] [12| ]
Green: [12| ] [11| ] [10| ] [9| ] [8| ] [7| ] [6| ] [5| ] [4| ] [3| ] [2| ]
Blue: [12| ] [11| ] [10| ] [9| ] [8| ] [7| ] [6| ] [5| ] [4| ] [3| ] [2| ]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedStringRepresentation, tc.input.Print())
		})
	}
}

func TestBoardImpl_MakeMove(t *testing.T) {
	type testCase struct {
		name                   string
		inputGameBoard         Board
		inputMove              Move
		expectedOk             bool
		expectedErr            error
		expectedGameBoardState Board
	}
	testCases := []testCase{
		{
			name:           "brand new game, valid move should leave board with a cell crossed off",
			inputGameBoard: NewGameBoard(),
			inputMove: Move{
				rowColor:   RowColorRed,
				cellNumber: 3,
			},
			expectedOk: true,
			expectedGameBoardState: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
				yellowRow: NewYellowRow(),
				greenRow:  NewGreenRow(),
				blueRow:   NewBlueRow(),
			},
		},
		{
			name:           "brand new game, invalid move should leave board unchanged",
			inputGameBoard: NewGameBoard(),
			inputMove: Move{
				rowColor:   RowColorRed,
				cellNumber: 1,
			},
			expectedOk:             false,
			expectedErr:            errors.New("cell number must be between 2 and 12"),
			expectedGameBoardState: NewGameBoard(),
		},
		{
			name: "mid game, valid move should leave board with a cell crossed off",
			inputGameBoard: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0}, false),
			},
			inputMove: Move{
				rowColor:   RowColorGreen,
				cellNumber: 4,
			},
			expectedOk: true,
			expectedGameBoardState: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0}, false),
			},
		},
		{
			name: "mid game, invalid move should leave board unchanged",
			inputGameBoard: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0}, false),
			},
			inputMove: Move{
				rowColor:   RowColorGreen,
				cellNumber: 10,
			},
			expectedOk:  false,
			expectedErr: errors.New("cell 10 is to the left of already crossed off cells"),
			expectedGameBoardState: &boardImpl{
				redRow:    newRedRowFromCells([]int{0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0}, false),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := tc.inputGameBoard.MakeMove(tc.inputMove)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.Equal(t, tc.expectedOk, ok)
				require.Equal(t, tc.expectedGameBoardState, tc.inputGameBoard)
			}

		})
	}
}

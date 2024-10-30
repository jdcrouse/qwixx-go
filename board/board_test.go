package board

import (
	"errors"
	"qwixx/actions"
	"testing"

	"github.com/stretchr/testify/require"
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
			expectedStringRepresentation: `Red: [2| ] [3| ] [4| ] [5| ] [6| ] [7| ] [8| ] [9| ] [10| ] [11| ] [12| ] [L| ]
Yellow: [2| ] [3| ] [4| ] [5| ] [6| ] [7| ] [8| ] [9| ] [10| ] [11| ] [12| ] [L| ]
Green: [12| ] [11| ] [10| ] [9| ] [8| ] [7| ] [6| ] [5| ] [4| ] [3| ] [2| ] [L| ]
Blue: [12| ] [11| ] [10| ] [9| ] [8| ] [7| ] [6| ] [5| ] [4| ] [3| ] [2| ] [L| ]`,
		},
		{
			name: "board with some marked cells",
			input: &boardImpl{
				redRow:    newRedRowFromCells([]int{1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0}, false),
				yellowRow: newYellowRowFromCells([]int{0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			},
			expectedStringRepresentation: `Red: [2|X] [3|X] [4| ] [5|X] [6| ] [7| ] [8| ] [9| ] [10| ] [11| ] [12| ] [L| ]
Yellow: [2| ] [3| ] [4|X] [5|X] [6|X] [7| ] [8| ] [9| ] [10| ] [11| ] [12| ] [L| ]
Green: [12| ] [11| ] [10| ] [9| ] [8|X] [7|X] [6|X] [5| ] [4| ] [3| ] [2| ] [L| ]
Blue: [12|X] [11|X] [10|X] [9| ] [8| ] [7| ] [6| ] [5| ] [4| ] [3| ] [2| ] [L| ]`,
		},
		{
			name: "board with locked rows",
			input: &boardImpl{
				redRow:    newRedRowFromCells([]int{1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1}, true),
				yellowRow: newYellowRowFromCells([]int{1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, false),
				greenRow:  newGreenRowFromCells([]int{0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0}, false),
				blueRow:   newBlueRowFromCells([]int{1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1}, true),
			},
			expectedStringRepresentation: `Red: [2|X] [3|X] [4|X] [5|X] [6|X] [7| ] [8| ] [9| ] [10| ] [11| ] [12|X] [L|X]
Yellow: [2|X] [3|X] [4|X] [5|X] [6|X] [7|X] [8| ] [9| ] [10| ] [11| ] [12| ] [L| ]
Green: [12| ] [11| ] [10| ] [9| ] [8|X] [7|X] [6|X] [5|X] [4|X] [3| ] [2| ] [L| ]
Blue: [12|X] [11|X] [10|X] [9|X] [8|X] [7| ] [6| ] [5| ] [4| ] [3| ] [2|X] [L|X]`,
		},
		{
			name: "board with all cells marked except lock",
			input: &boardImpl{
				redRow:    newRedRowFromCells([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, false),
				yellowRow: newYellowRowFromCells([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, false),
				greenRow:  newGreenRowFromCells([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, false),
				blueRow:   newBlueRowFromCells([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, false),
			},
			expectedStringRepresentation: `Red: [2|X] [3|X] [4|X] [5|X] [6|X] [7|X] [8|X] [9|X] [10|X] [11|X] [12| ] [L| ]
Yellow: [2|X] [3|X] [4|X] [5|X] [6|X] [7|X] [8|X] [9|X] [10|X] [11|X] [12| ] [L| ]
Green: [12|X] [11|X] [10|X] [9|X] [8|X] [7|X] [6|X] [5|X] [4|X] [3|X] [2| ] [L| ]
Blue: [12|X] [11|X] [10|X] [9|X] [8|X] [7|X] [6|X] [5|X] [4|X] [3|X] [2| ] [L| ]`,
		},
		{
			name: "board with mixed states",
			input: &boardImpl{
				redRow:    newRedRowFromCells([]int{1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1}, false),
				yellowRow: newYellowRowFromCells([]int{0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0}, true),
				greenRow:  newGreenRowFromCells([]int{1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0}, false),
				blueRow:   newBlueRowFromCells([]int{0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1}, true),
			},
			expectedStringRepresentation: `Red: [2|X] [3|X] [4|X] [5| ] [6| ] [7|X] [8|X] [9| ] [10|X] [11| ] [12|X] [L|X]
Yellow: [2| ] [3|X] [4|X] [5|X] [6| ] [7| ] [8|X] [9|X] [10|X] [11|X] [12| ] [L| ]
Green: [12|X] [11| ] [10|X] [9|X] [8|X] [7| ] [6| ] [5|X] [4|X] [3|X] [2| ] [L| ]
Blue: [12| ] [11|X] [10|X] [9|X] [8|X] [7|X] [6| ] [5| ] [4|X] [3|X] [2|X] [L|X]`,
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
		inputMove              actions.Move
		expectedOk             bool
		expectedErr            error
		expectedGameBoardState Board
	}
	testCases := []testCase{
		{
			name:           "brand new game, valid move should leave board with a cell crossed off",
			inputGameBoard: NewGameBoard(),
			inputMove: actions.Move{
				RowColor:   actions.RowColorRed,
				CellNumber: 3,
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
			inputMove: actions.Move{
				RowColor:   actions.RowColorRed,
				CellNumber: 1,
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
			inputMove: actions.Move{
				RowColor:   actions.RowColorGreen,
				CellNumber: 4,
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
			inputMove: actions.Move{
				RowColor:   actions.RowColorGreen,
				CellNumber: 10,
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
			err := tc.inputGameBoard.MakeMove(tc.inputMove)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedGameBoardState, tc.inputGameBoard)
			}

		})
	}
}

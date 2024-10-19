package board

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrintRow(t *testing.T) {
	type testCase struct {
		name                         string
		input                        Row
		expectedStringRepresentation string
	}
	testCases := []testCase{
		{
			name:                         "base case ascending row",
			input:                        newValidatedAscendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			expectedStringRepresentation: `[2| ] [3| ] [4| ] [5| ] [6| ] [7| ] [8| ] [9| ] [10| ] [11| ] [12| ] [L| ]`,
		},
		{
			name:                         "base case descending row",
			input:                        newValidatedDescendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			expectedStringRepresentation: `[12| ] [11| ] [10| ] [9| ] [8| ] [7| ] [6| ] [5| ] [4| ] [3| ] [2| ] [L| ]`,
		},
		{
			name:                         "non-base case ascending row",
			input:                        newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0}, false),
			expectedStringRepresentation: `[2| ] [3|X] [4|X] [5| ] [6|X] [7| ] [8| ] [9|X] [10| ] [11|X] [12| ] [L| ]`,
		},
		{
			name:                         "non-base case descending row",
			input:                        newValidatedDescendingRowFromCells(t, []int{1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0}, false),
			expectedStringRepresentation: `[12|X] [11| ] [10| ] [9|X] [8|X] [7|X] [6|X] [5| ] [4|X] [3| ] [2| ] [L| ]`,
		},
		{
			name:                         "ascending row locked by this player",
			input:                        newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 1}, true),
			expectedStringRepresentation: `[2| ] [3|X] [4|X] [5| ] [6|X] [7| ] [8| ] [9|X] [10| ] [11|X] [12|X] [L|X]`,
		},
		{
			name:                         "descending row locked by this player",
			input:                        newValidatedDescendingRowFromCells(t, []int{1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 1}, true),
			expectedStringRepresentation: `[12|X] [11| ] [10| ] [9|X] [8|X] [7|X] [6|X] [5| ] [4|X] [3| ] [2|X] [L|X]`,
		},
		{
			name:                         "ascending row locked by other player",
			input:                        newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0}, true),
			expectedStringRepresentation: `[2| ] [3|X] [4|X] [5| ] [6|X] [7| ] [8| ] [9|X] [10| ] [11|X] [12| ] [L| ]`,
		},
		{
			name:                         "descending row locked by other player",
			input:                        newValidatedDescendingRowFromCells(t, []int{1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0}, true),
			expectedStringRepresentation: `[12|X] [11| ] [10| ] [9|X] [8|X] [7|X] [6|X] [5| ] [4|X] [3| ] [2| ] [L| ]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedStringRepresentation, tc.input.Print())
		})
	}
}

func TestIsMoveValid(t *testing.T) {
	type testCase struct {
		name            string
		input           Row
		inputCellNumber int
		expectedErr     error
	}
	validMoveCases := []testCase{
		{
			name:            "valid move: empty ascending row, cross off leftmost cell",
			input:           newValidatedAscendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 2,
		},
		{
			name:            "valid move: empty ascending row, cross off any middle cell",
			input:           newValidatedAscendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 7,
		},
		{
			name:            "valid move: ascending row, cross off empty cell with nothing to the right",
			input:           newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 8,
		},
		{
			name:            "valid move: ascending row, can only cross off rightmost cell if five cells are already crossed",
			input:           newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 12,
		},
		{
			name:            "valid move: empty descending row, cross off leftmost cell = true",
			input:           newValidatedDescendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 12,
		},
		{
			name:            "valid move: descending row, can only cross off rightmost cell if five cells are already crossed",
			input:           newValidatedDescendingRowFromCells(t, []int{0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 2,
		},
		{
			name:            "valid move: non-empty descending row, cross off empty cell with nothing to the right",
			input:           newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 6,
		},
	}
	invalidMoveCases := []testCase{
		{
			name:            "invalid move: ascending row, cross off already crossed off cell",
			input:           newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 4,
			expectedErr:     errors.New("cell 4 is already crossed off"),
		},

		{
			name:            "invalid move: non-empty ascending row, cross off empty cell with already crossed off cells to the right",
			input:           newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 6,
			expectedErr:     errors.New("cell 6 is to the left of already crossed off cells"),
		},

		{
			name:            "invalid move: ascending row, can only cross off rightmost cell if five cells are already crossed",
			input:           newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 12,
			expectedErr:     errors.New("cannot cross off rightmost cell of row unless 5 cells have been crossed off in that row"),
		},

		{
			name:            "invalid move: non-empty descending row, cross off already crossed off cell",
			input:           newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 10,
			expectedErr:     errors.New("cell 10 is already crossed off"),
		},

		{
			name:            "invalid move: descending row, can only cross off rightmost cell if five cells are already crossed",
			input:           newValidatedDescendingRowFromCells(t, []int{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 2,
			expectedErr:     errors.New("cannot cross off rightmost cell of row unless 5 cells have been crossed off in that row"),
		},

		{
			name:            "invalid move: non-empty descending row, cross off empty cell with already crossed off cells to the right",
			input:           newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 8,
			expectedErr:     errors.New("cell 8 is to the left of already crossed off cells"),
		},
		{
			name:            "invalid move: locked non-empty row regardless of actual move",
			input:           newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, true),
			inputCellNumber: 8,
			expectedErr:     errors.New("row is locked"),
		},
		{
			name:            "invalid move: locked empty row regardless of actual move",
			input:           newValidatedDescendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, true),
			inputCellNumber: 8,
			expectedErr:     errors.New("row is locked"),
		},
	}
	for _, tc := range validMoveCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := tc.input.IsMoveValid(tc.inputCellNumber)
			require.NoError(t, err)
			require.True(t, ok)
		})
	}
	for _, tc := range invalidMoveCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := tc.input.IsMoveValid(tc.inputCellNumber)
			require.Error(t, err)
			require.False(t, ok)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestMakeMove(t *testing.T) {
	type testCase struct {
		name            string
		inputRow        Row
		inputCellNumber int
		expectedOk      bool
		expectedErr     error
		expectedRow     Row
	}
	testCases := []testCase{
		{
			name:            "empty ascending row, cross off leftmost cell succeeds",
			inputRow:        newValidatedAscendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 2,
			expectedOk:      true,
			expectedErr:     nil,
			expectedRow:     newValidatedAscendingRowFromCells(t, []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
		},
		{
			name:            "non-empty ascending row, cross off non-empty cell fails",
			inputRow:        newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 4,
			expectedOk:      false,
			expectedErr:     errors.New("cell 4 is already crossed off"),
			expectedRow:     newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
		},
		{
			name:            "non-empty ascending row, cross off empty cell with nothing to the right succeeds",
			inputRow:        newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 8,
			expectedOk:      true,
			expectedErr:     nil,
			expectedRow:     newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0}, false),
		},
		{
			name:            "non-empty ascending row, cross off empty cell with non-empty cells to the right fails",
			inputRow:        newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 6,
			expectedOk:      false,
			expectedErr:     errors.New("cell 6 is to the left of already crossed off cells"),
			expectedRow:     newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
		},
		{
			name:            "empty descending row, cross off leftmost cell succeeds",
			inputRow:        newValidatedDescendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 12,
			expectedOk:      true,
			expectedErr:     nil,
			expectedRow:     newValidatedDescendingRowFromCells(t, []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
		},
		{
			name:            "non-empty descending row, cross off non-empty cell fails",
			inputRow:        newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 10,
			expectedOk:      false,
			expectedErr:     errors.New("cell 10 is already crossed off"),
			expectedRow:     newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
		},
		{
			name:            "non-empty descending row, cross off empty cell with nothing to the right succeeds",
			inputRow:        newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 6,
			expectedOk:      true,
			expectedErr:     nil,
			expectedRow:     newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0}, false),
		},
		{
			name:            "non-empty descending row, cross off empty cell with non-empty cells to the right fails",
			inputRow:        newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
			inputCellNumber: 8,
			expectedOk:      false,
			expectedErr:     errors.New("cell 8 is to the left of already crossed off cells"),
			expectedRow:     newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, false),
		},
		{
			name:            "valid move but row is locked, any move fails (ascending)",
			inputRow:        newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, true),
			inputCellNumber: 8,
			expectedOk:      false,
			expectedErr:     errors.New("cannot make move, row is already locked"),
			expectedRow:     newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, true),
		},
		{
			name:            "valid move but row is locked, any move fails (descending)",
			inputRow:        newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, true),
			inputCellNumber: 6,
			expectedOk:      false,
			expectedErr:     errors.New("cannot make move, row is already locked"),
			expectedRow:     newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, true),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := tc.inputRow.MakeMove(tc.inputCellNumber)
			require.Equal(t, tc.expectedOk, ok)
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedRow, tc.inputRow)
		})
	}
}

func newValidatedAscendingRowFromCells(t *testing.T, cells []int, locked bool) Row {
	validateRowInvariants(t, cells)
	return newAscendingRowFromCells(cells, locked)
}

func newValidatedDescendingRowFromCells(t *testing.T, cells []int, locked bool) Row {
	validateRowInvariants(t, cells)
	return newDescendingRowFromCells(cells, locked)
}

// validateRowInvariants ensures that the given cells
// - is of length 11
// - only contains 0s or 1s
func validateRowInvariants(t *testing.T, cells []int) {
	require.Equal(t, len(cells), 11)
	for _, value := range cells {
		if value != 0 && value != 1 {
			require.True(t, value == 0 || value == 1)
		}
	}
}

func TestIndexToCellNumber(t *testing.T) {
	type testCase struct {
		name               string
		rowType            rowType
		inputIndex         int
		expectedCellNumber int
		expectedErr        error
	}
	testCases := []testCase{
		{
			name:        "invalid too low index ascending",
			rowType:     RowTypeAscending,
			inputIndex:  -1,
			expectedErr: errors.New("invalid index: -1. must be between 0 and 10"),
		},
		{
			name:               "leftmost index ascending",
			rowType:            RowTypeAscending,
			inputIndex:         0,
			expectedCellNumber: 2,
		},
		{
			name:               "index in middle ascending",
			rowType:            RowTypeAscending,
			inputIndex:         6,
			expectedCellNumber: 8,
		},
		{
			name:               "rightmost index ascending",
			rowType:            RowTypeAscending,
			inputIndex:         10,
			expectedCellNumber: 12,
		},
		{
			name:        "invalid too high index ascending",
			rowType:     RowTypeAscending,
			inputIndex:  11,
			expectedErr: errors.New("invalid index: 11. must be between 0 and 10"),
		},
		{
			name:        "invalid too low index descending",
			rowType:     RowTypeDescending,
			inputIndex:  -1,
			expectedErr: errors.New("invalid index: -1. must be between 0 and 10"),
		},
		{
			name:               "leftmost index descending",
			rowType:            RowTypeDescending,
			inputIndex:         0,
			expectedCellNumber: 12,
		},
		{
			name:               "index in middle descending",
			rowType:            RowTypeDescending,
			inputIndex:         6,
			expectedCellNumber: 6,
		},
		{
			name:               "rightmost index descending (the lock cell)",
			rowType:            RowTypeDescending,
			inputIndex:         10,
			expectedCellNumber: 2,
		},
		{
			name:        "invalid too high index descending",
			rowType:     RowTypeDescending,
			inputIndex:  11,
			expectedErr: errors.New("invalid index: 11. must be between 0 and 10"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cellNumber, err := indexToCellNumber(tc.rowType, tc.inputIndex)
			if tc.expectedErr != nil {
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.Equal(t, tc.expectedCellNumber, cellNumber)
			}
		})
	}
}

func TestCellNumberToIndex(t *testing.T) {
	type testCase struct {
		name            string
		rowType         rowType
		inputCellNumber int
		expectedIndex   int
		expectedErr     error
	}
	testCases := []testCase{
		{
			name:            "invalid too low cell number ascending",
			rowType:         RowTypeAscending,
			inputCellNumber: 1,
			expectedErr:     errors.New("invalid cell number: 1. must be between 2 and 12"),
		},
		{
			name:            "leftmost cell number ascending",
			rowType:         RowTypeAscending,
			inputCellNumber: 2,
			expectedIndex:   0,
		},
		{
			name:            "cell number im middle ascending",
			rowType:         RowTypeAscending,
			inputCellNumber: 8,
			expectedIndex:   6,
		},
		{
			name:            "rightmost cell number ascending (lock)",
			rowType:         RowTypeAscending,
			inputCellNumber: 12,
			expectedIndex:   10,
		},
		{
			name:            "invalid too high cell number ascending",
			rowType:         RowTypeAscending,
			inputCellNumber: 13,
			expectedErr:     errors.New("invalid cell number: 13. must be between 2 and 12"),
		},
		{
			name:            "leftmost cell number descending",
			rowType:         RowTypeDescending,
			inputCellNumber: 12,
			expectedIndex:   0,
		},
		{
			name:            "cell number in middle descending",
			rowType:         RowTypeDescending,
			inputCellNumber: 6,
			expectedIndex:   6,
		},
		{
			name:            "rightmost cell number descending (lock)",
			rowType:         RowTypeDescending,
			inputCellNumber: 2,
			expectedIndex:   10,
		},
		{
			name:            "invalid too high cell number descending",
			rowType:         RowTypeDescending,
			inputCellNumber: 13,
			expectedErr:     errors.New("invalid cell number: 13. must be between 2 and 12"),
		},
		{
			name:            "invalid too low cell number descending",
			rowType:         RowTypeDescending,
			inputCellNumber: 1,
			expectedErr:     errors.New("invalid cell number: 1. must be between 2 and 12"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			index, err := cellNumberToIndex(tc.rowType, tc.inputCellNumber)
			if tc.expectedErr != nil {
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.Equal(t, tc.expectedIndex, index)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	type testCase struct {
		name          string
		input         Row
		expectedScore int
	}
	// TODO testing every case here is a little extra without iterating in some way, but ill change that later
	testCases := []testCase{
		{
			name:          "ascending row with 0",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			expectedScore: 0,
		},
		{
			name:          "ascending row with 1",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			expectedScore: 1,
		},
		{
			name:          "ascending row with 2 (last cell counts as +1 for lock cell)",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, false),
			expectedScore: 3,
		},
		{
			name:          "ascending row with 3",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0}, false),
			expectedScore: 6,
		},
		{
			name:          "ascending row with 4",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0}, false),
			expectedScore: 10,
		},
		{
			name:          "ascending row with 5",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0}, false),
			expectedScore: 15,
		},
		{
			name:          "ascending row with 6",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0}, false),
			expectedScore: 21,
		},
		{
			name:          "ascending row with 7 (last cell counts as +1 for lock cell)",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 1}, true),
			expectedScore: 28,
		},
		{
			name:          "ascending row with 8",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0}, true),
			expectedScore: 36,
		},
		{
			name:          "ascending row with 9",
			input:         newValidatedAscendingRowFromCells(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0}, false),
			expectedScore: 45,
		},
		{
			name:          "ascending row with 10 (last cell counts as +1 for lock cell)",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1}, true),
			expectedScore: 55,
		},
		{
			name:          "ascending row with 11 (last cell counts as +1 for lock cell)",
			input:         newValidatedAscendingRowFromCells(t, []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, true),
			expectedScore: 66,
		},
		{
			name:          "ascending row with 12 (last cell counts as +1 for lock cell)",
			input:         newValidatedAscendingRowFromCells(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, true),
			expectedScore: 78,
		},
		// scoring doesn't differ between ascending and descending, but testing for safety
		{
			name:          "descending row with 0",
			input:         newValidatedDescendingRowFromCells(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			expectedScore: 0,
		},
		{
			name:          "descending row with 1",
			input:         newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, false),
			expectedScore: 1,
		},
		{
			name:          "descending row with 2",
			input:         newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0}, false),
			expectedScore: 3,
		},
		{
			name:          "descending row with 3",
			input:         newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0}, false),
			expectedScore: 6,
		},
		{
			name:          "descending row with 4",
			input:         newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 0}, false),
			expectedScore: 10,
		},
		{
			name:          "descending row with 5",
			input:         newValidatedDescendingRowFromCells(t, []int{0, 0, 1, 1, 0, 1, 1, 0, 0, 1, 0}, false),
			expectedScore: 15,
		},
		{
			name:          "descending row with 6",
			input:         newValidatedDescendingRowFromCells(t, []int{1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0}, false),
			expectedScore: 21,
		},
		{
			name:          "descending row with 7",
			input:         newValidatedDescendingRowFromCells(t, []int{1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0}, false),
			expectedScore: 28,
		},
		{
			name:          "descending row with 8",
			input:         newValidatedDescendingRowFromCells(t, []int{1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 0}, false),
			expectedScore: 36,
		},
		{
			name:          "descending row with 9",
			input:         newValidatedDescendingRowFromCells(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0}, false),
			expectedScore: 45,
		},
		{
			name:          "descending row with 10",
			input:         newValidatedDescendingRowFromCells(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, false),
			expectedScore: 55,
		},
		{
			name:          "descending row with 11 (last cell counts as +1 for lock cell)",
			input:         newValidatedDescendingRowFromCells(t, []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, true),
			expectedScore: 66,
		},
		{
			name:          "descending row with 12 (last cell counts as +1 for lock cell)",
			input:         newValidatedDescendingRowFromCells(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, true),
			expectedScore: 78,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedScore, tc.input.CalculateScore())
		})
	}
}

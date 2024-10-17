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
			input:                        newRowFromCells(RowTypeAscending, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expectedStringRepresentation: `[2| ] [3| ] [4| ] [5| ] [6| ] [7| ] [8| ] [9| ] [10| ] [11| ] [12| ]`,
		},
		{
			name:                         "base case descending row",
			input:                        newRowFromCells(RowTypeDescending, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expectedStringRepresentation: `[12| ] [11| ] [10| ] [9| ] [8| ] [7| ] [6| ] [5| ] [4| ] [3| ] [2| ]`,
		},
		{
			name:                         "non-base case ascending row", // TODO show locked
			input:                        newRowFromCells(RowTypeAscending, []int{0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0}),
			expectedStringRepresentation: `[2| ] [3|X] [4|X] [5| ] [6|X] [7| ] [8| ] [9|X] [10| ] [11|X] [12| ]`,
		},
		{
			name:                         "non-base case descending row",
			input:                        newRowFromCells(RowTypeDescending, []int{1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0}),
			expectedStringRepresentation: `[12|X] [11| ] [10| ] [9|X] [8|X] [7|X] [6|X] [5| ] [4|X] [3| ] [2| ]`,
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
		inputCells      []int
		inputRowType    rowType
		inputCellNumber int
		inputLocked     bool
		expectedErr     error
	}
	testCases := []testCase{
		{
			name:            "empty ascending row, cross off leftmost cell = true",
			inputCells:      []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeAscending,
			inputCellNumber: 2,
			inputLocked:     false,
			expectedErr:     nil,
		},
		{
			name:            "non-empty ascending row, cross off non-empty cell = false",
			inputCells:      []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeAscending,
			inputCellNumber: 4,
			inputLocked:     false,
			expectedErr:     errors.New("cell 4 is already crossed off"),
		},
		{
			name:            "non-empty ascending row, cross off empty cell with nothing to the right = true",
			inputCells:      []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeAscending,
			inputCellNumber: 8,
			inputLocked:     false,
			expectedErr:     nil,
		},
		{
			name:            "non-empty ascending row, cross off empty cell with non-empty cells to the right = false",
			inputCells:      []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeAscending,
			inputCellNumber: 6,
			inputLocked:     false,
			expectedErr:     errors.New("cell 6 is to the left of already crossed off cells"),
		},
		{
			name:            "empty descending row, cross off leftmost cell = true",
			inputCells:      []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeDescending,
			inputCellNumber: 12,
			inputLocked:     false,
			expectedErr:     nil,
		},
		{
			name:            "non-empty descending row, cross off non-empty cell = false",
			inputCells:      []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeDescending,
			inputCellNumber: 10,
			inputLocked:     false,
			expectedErr:     errors.New("cell 10 is already crossed off"),
		},
		{
			name:            "non-empty descending row, cross off empty cell with nothing to the right = true",
			inputCells:      []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeDescending,
			inputCellNumber: 6,
			inputLocked:     false,
			expectedErr:     nil,
		},
		{
			name:            "non-empty descending row, cross off empty cell with non-empty cells to the right = false",
			inputCells:      []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeDescending,
			inputCellNumber: 8,
			inputLocked:     false,
			expectedErr:     errors.New("cell 8 is to the left of already crossed off cells"),
		},
		{
			name:            "locked non-empty row regardless of actual move = false",
			inputCells:      []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeAscending,
			inputCellNumber: 8,
			inputLocked:     true,
			expectedErr:     errors.New("row is locked"),
		},
		{
			name:            "locked empty row regardless of actual move = false",
			inputCells:      []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			inputRowType:    RowTypeDescending,
			inputCellNumber: 8,
			inputLocked:     true,
			expectedErr:     errors.New("row is locked"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validateRowInvariants(t, tc.inputCells)
			err := isMoveValid(tc.inputCells, tc.inputRowType, tc.inputLocked, tc.inputCellNumber)
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
			inputRow:        newRowFromCells(RowTypeAscending, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			inputCellNumber: 2,
			expectedOk:      true,
			expectedErr:     nil,
			expectedRow:     newRowFromCells(RowTypeAscending, []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		},
		{
			name:            "non-empty ascending row, cross off non-empty cell fails",
			inputRow:        newRowFromCells(RowTypeAscending, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}),
			inputCellNumber: 4,
			expectedOk:      false,
			expectedErr:     errors.New("cell 4 is already crossed off"),
			expectedRow:     newRowFromCells(RowTypeAscending, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}),
		},
		{
			name:            "non-empty ascending row, cross off empty cell with nothing to the right succeeds",
			inputRow:        newRowFromCells(RowTypeAscending, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}),
			inputCellNumber: 8,
			expectedOk:      true,
			expectedErr:     nil,
			expectedRow:     newRowFromCells(RowTypeAscending, []int{0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0}),
		},
		{
			name:            "non-empty ascending row, cross off empty cell with non-empty cells to the right fails",
			inputRow:        newRowFromCells(RowTypeAscending, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}),
			inputCellNumber: 6,
			expectedOk:      false,
			expectedErr:     errors.New("cell 6 is to the left of already crossed off cells"),
			expectedRow:     newRowFromCells(RowTypeAscending, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}),
		},
		{
			name:            "empty descending row, cross off leftmost cell succeeds",
			inputRow:        newRowFromCells(RowTypeDescending, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			inputCellNumber: 12,
			expectedOk:      true,
			expectedErr:     nil,
			expectedRow:     newRowFromCells(RowTypeDescending, []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		},
		{
			name:            "non-empty descending row, cross off non-empty cell fails",
			inputRow:        newRowFromCells(RowTypeDescending, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}),
			inputCellNumber: 10,
			expectedOk:      false,
			expectedErr:     errors.New("cell 10 is already crossed off"),
			expectedRow:     newRowFromCells(RowTypeDescending, []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}),
		},
		{
			name:            "non-empty descending row, cross off empty cell with nothing to the right succeeds",
			inputRow:        newRowFromCells(RowTypeDescending, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}),
			inputCellNumber: 6,
			expectedOk:      true,
			expectedErr:     nil,
			expectedRow:     newRowFromCells(RowTypeDescending, []int{0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0}),
		},
		{
			name:            "non-empty descending row, cross off empty cell with non-empty cells to the right fails",
			inputRow:        newRowFromCells(RowTypeDescending, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}),
			inputCellNumber: 8,
			expectedOk:      false,
			expectedErr:     errors.New("cell 8 is to the left of already crossed off cells"),
			expectedRow:     newRowFromCells(RowTypeDescending, []int{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}),
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

// validateRowInvariants ensures that the given cells is of length 11 and only contains 0s or 1s
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
			expectedErr: errors.New("invalid index -1"),
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
			expectedErr: errors.New("invalid index 11"),
		},
		{
			name:        "invalid too low index descending",
			rowType:     RowTypeDescending,
			inputIndex:  -1,
			expectedErr: errors.New("invalid index -1"),
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
			name:               "rightmost index descending",
			rowType:            RowTypeDescending,
			inputIndex:         10,
			expectedCellNumber: 2,
		},
		{
			name:        "invalid too high index descending",
			rowType:     RowTypeDescending,
			inputIndex:  11,
			expectedErr: errors.New("invalid index 11"),
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
			expectedErr:     errors.New("invalid cell number 1"),
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
			name:            "rightmost cell number ascending",
			rowType:         RowTypeAscending,
			inputCellNumber: 12,
			expectedIndex:   10,
		},
		{
			name:            "invalid too high cell number ascending",
			rowType:         RowTypeAscending,
			inputCellNumber: 13,
			expectedErr:     errors.New("invalid cell number 13"),
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
			name:            "rightmost cell number descending",
			rowType:         RowTypeDescending,
			inputCellNumber: 2,
			expectedIndex:   10,
		},
		{
			name:            "invalid too high cell number descending",
			rowType:         RowTypeDescending,
			inputCellNumber: 13,
			expectedErr:     errors.New("invalid cell number 13"),
		},
		{
			name:            "invalid too low cell number descending",
			rowType:         RowTypeDescending,
			inputCellNumber: 1,
			expectedErr:     errors.New("invalid cell number 1"),
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

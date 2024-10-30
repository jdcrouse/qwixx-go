package rule_checker

import (
	"qwixx/internal/game/actions"
	"qwixx/internal/game/board"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWhiteDiceMoveIsValidForBoard(t *testing.T) {
	type testCase struct {
		name              string
		inputBoard        board.Board
		inputDiceRoll     actions.DiceRoll
		inputProposedMove actions.Move
		expectedValid     bool
	}
	testCases := []testCase{
		{
			name:       "valid move",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 8,
				RowColor:   actions.RowColorBlue,
			},
			expectedValid: true,
		},
		{
			name:       "invalid move: cell number doesn't match white dice sum",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 10,
				RowColor:   actions.RowColorBlue,
			},
			expectedValid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedValid, WhiteDiceMoveIsValidForBoard(tc.inputBoard, tc.inputDiceRoll, tc.inputProposedMove))
		})
	}
}

func TestColorDiceMoveIsValidForBoard(t *testing.T) {
	type testCase struct {
		name              string
		inputBoard        board.Board
		inputDiceRoll     actions.DiceRoll
		inputProposedMove actions.Move
		expectedValid     bool
	}
	testCases := []testCase{
		{
			name:       "valid move with red die sum 1",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 4,
				RowColor:   actions.RowColorRed,
			},
			expectedValid: true,
		},
		{
			name:       "valid move with red die sum 2",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 6,
				RowColor:   actions.RowColorRed,
			},
			expectedValid: true,
		},
		{
			name:       "invalid move: cell number doesn't match either red die sum",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 10,
				RowColor:   actions.RowColorRed,
			},
			expectedValid: false,
		},
		{
			name:       "valid move with yellow die sum 1",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 7,
				RowColor:   actions.RowColorYellow,
			},
			expectedValid: true,
		},
		{
			name:       "valid move with yellow die sum 2",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 9,
				RowColor:   actions.RowColorYellow,
			},
			expectedValid: true,
		},
		{
			name:       "invalid move: cell number doesn't match either yellow die sum",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 10,
				RowColor:   actions.RowColorYellow,
			},
			expectedValid: false,
		},
		{
			name:       "valid move with green die sum 1",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 6,
				RowColor:   actions.RowColorGreen,
			},
			expectedValid: true,
		},
		{
			name:       "valid move with green die sum 2",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 8,
				RowColor:   actions.RowColorGreen,
			},
			expectedValid: true,
		},
		{
			name:       "invalid move: cell number doesn't match either green die sum",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 2,
				RowColor:   actions.RowColorGreen,
			},
			expectedValid: false,
		},
		{
			name:       "valid move with blue die sum 1",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 5,
				RowColor:   actions.RowColorBlue,
			},
			expectedValid: true,
		},
		{
			name:       "valid move with blue die sum 2",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 7,
				RowColor:   actions.RowColorBlue,
			},
			expectedValid: true,
		},
		{
			name:       "invalid move: cell number doesn't match either blue die sum",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 3,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    1,
					Blue:   2,
					Green:  3,
					Yellow: 4,
				},
			},
			inputProposedMove: actions.Move{
				CellNumber: 3,
				RowColor:   actions.RowColorBlue,
			},
			expectedValid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedValid, ColorDiceMoveIsValidForBoard(tc.inputBoard, tc.inputDiceRoll, tc.inputProposedMove))
		})
	}
}

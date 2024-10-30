package qwixx

import (
	"qwixx/actions"
	"qwixx/board"
	"qwixx/player"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunGame(t *testing.T) {
	alice := player.NewComputerPlayer("alice")
	bob := player.NewComputerPlayer("bob")
	charlie := player.NewComputerPlayer("charlie")
	players := []player.Player{alice, bob, charlie}
	runner := NewGameRunner(players)
	runner.RunGame()
}

func TestIsActiveTurnPenalty(t *testing.T) {
	type testCase struct {
		name           string
		input          actions.ActivePlayerTurn
		expectedOutput bool
	}
	testCases := []testCase{
		{
			name:           "both moves are nil means penalty",
			input:          actions.ActivePlayerTurn{},
			expectedOutput: true,
		},
		{
			name: "non-nil white dice move means NOT a penalty",
			input: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{
					RowColor:   actions.RowColorBlue,
					CellNumber: 5,
				},
			},
			expectedOutput: false,
		},
		{
			name: "non-nil color dice move means NOT a penalty",
			input: actions.ActivePlayerTurn{
				ColorDiceMove: &actions.Move{
					RowColor:   actions.RowColorYellow,
					CellNumber: 9,
				},
			},
			expectedOutput: false,
		},
		{
			name: "both non-nil moves means NOT a penalty",
			input: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{
					RowColor:   actions.RowColorRed,
					CellNumber: 8,
				},
				ColorDiceMove: &actions.Move{
					RowColor:   actions.RowColorGreen,
					CellNumber: 3,
				},
			},
			expectedOutput: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedOutput, isActiveTurnPenalty(tc.input))
		})
	}
}

func TestIsActiveTurnValid(t *testing.T) {
	type testCase struct {
		name              string
		inputBoard        board.Board
		inputDiceRoll     actions.DiceRoll
		inputProposedTurn actions.ActivePlayerTurn
		expectedOutput    bool
	}
	testCases := []testCase{
		{
			name:              "explicit penalty is valid",
			inputBoard:        board.NewGameBoard(),
			inputDiceRoll:     actions.RollQwixxDice(),
			inputProposedTurn: actions.ActivePlayerTurn{},
			expectedOutput:    true,
		},
		{
			name:       "turn with one valid white dice move where cell number DOES match dice sum is valid",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 4,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{},
			},
			inputProposedTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{
					RowColor:   actions.RowColorRed,
					CellNumber: 9,
				},
			},
			expectedOutput: true,
		},
		{
			name:       "turn with one valid white dice move where cell number DOES NOT MATCH dice sum is invalid",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 4,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{},
			},
			inputProposedTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{
					RowColor:   actions.RowColorRed,
					CellNumber: 10,
				},
			},
			expectedOutput: false,
		},
		{
			name:       "turn with one valid color dice move where cell number DOES MATCH possible dice sum is valid",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 4,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    4,
					Yellow: 3,
					Green:  6,
					Blue:   2,
				},
			},
			inputProposedTurn: actions.ActivePlayerTurn{
				ColorDiceMove: &actions.Move{
					RowColor:   actions.RowColorYellow,
					CellNumber: 7,
				},
			},
			expectedOutput: true,
		},
		{
			name:       "turn with one valid color dice move where cell number DOES NOT MATCH possible dice sum is invalid",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 4,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    4,
					Yellow: 3,
					Green:  6,
					Blue:   2,
				},
			},
			inputProposedTurn: actions.ActivePlayerTurn{
				ColorDiceMove: &actions.Move{
					RowColor:   actions.RowColorYellow,
					CellNumber: 4,
				},
			},
			expectedOutput: false,
		},
		{
			name:       "turn with invalid white dice move, any color dice move is invalid",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 4,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    4,
					Yellow: 3,
					Green:  6,
					Blue:   2,
				},
			},
			inputProposedTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{
					RowColor:   actions.RowColorBlue,
					CellNumber: 8,
				},
				ColorDiceMove: &actions.Move{
					RowColor:   actions.RowColorYellow,
					CellNumber: 10,
				},
			},
			expectedOutput: false,
		},
		{
			name:       "turn with valid white dice move, invalid color dice move is invalid",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 4,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    4,
					Yellow: 3,
					Green:  6,
					Blue:   2,
				},
			},
			inputProposedTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{
					RowColor:   actions.RowColorBlue,
					CellNumber: 9,
				},
				ColorDiceMove: &actions.Move{
					RowColor:   actions.RowColorYellow,
					CellNumber: 4,
				},
			},
			expectedOutput: false,
		},

		{
			name:       "turn with two valid moves is valid",
			inputBoard: board.NewGameBoard(),
			inputDiceRoll: actions.DiceRoll{
				WhiteDiceRoll: actions.WhiteDiceRoll{
					White1: 4,
					White2: 5,
				},
				ColorDiceRoll: actions.ColorDiceRoll{
					Red:    4,
					Yellow: 3,
					Green:  6,
					Blue:   2,
				},
			},
			inputProposedTurn: actions.ActivePlayerTurn{
				WhiteDiceMove: &actions.Move{
					RowColor:   actions.RowColorBlue,
					CellNumber: 9,
				},
				ColorDiceMove: &actions.Move{
					RowColor:   actions.RowColorBlue,
					CellNumber: 7,
				},
			},
			expectedOutput: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedOutput,
				isActiveTurnValid(
					tc.inputBoard,
					tc.inputDiceRoll,
					tc.inputProposedTurn,
				),
			)
		})
	}
}

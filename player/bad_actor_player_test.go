package player

import (
	"qwixx/actions"
	"qwixx/board"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPromptTurnCantMutateBoardCopy tests that a board COPY given to a bad actor player
// who tries to mutate the board will not affect the original board
func TestPromptTurnCantMutateBoardCopy(t *testing.T) {
	newBoard := board.NewGameBoard()
	pl := BadActorPlayer{}
	diceRoll := actions.DiceRoll{
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
	}
	_ = pl.PromptActivePlayerTurn(newBoard.Copy(), diceRoll)
	require.Equal(t, board.NewGameBoard(), newBoard)
}

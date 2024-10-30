package board

import (
	"qwixx/internal/game/actions"
)

// ApplyActivePlayerTurn applies the given turn to the given board, applying the white dice move and then the color dice move
// a partially applied turn that results in an error will return the error, but the given board is still updated so this should be called with a copy
func ApplyActivePlayerTurn(board Board, turn actions.ActivePlayerTurn) (Board, error) {
	if turn.WhiteDiceMove != nil {
		if err := board.MakeMove(*turn.WhiteDiceMove); err != nil {
			return board, err
		}
	}
	if turn.ColorDiceMove != nil {
		if err := board.MakeMove(*turn.ColorDiceMove); err != nil {
			return board, err
		}
	}
	return board, nil
}

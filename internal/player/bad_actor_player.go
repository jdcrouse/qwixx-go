package player

import (
	"qwixx/internal/actions"
	"qwixx/internal/board"
)

var _ Player = BadActorPlayer{}

// BadActorPlayer is a player that tries to break the game, to be used for testing purposes
type BadActorPlayer struct {
}

func (b BadActorPlayer) GetName() string {
	//TODO implement me
	panic("implement me")
}

func (b BadActorPlayer) InformOfPlayOrder(playerNames []string) {
	//TODO implement me
	panic("implement me")
}

// PromptActivePlayerTurn attempts to fill in the entire red row other than the final locking cell
func (b BadActorPlayer) PromptActivePlayerTurn(playerBoard board.Board, diceRoll actions.DiceRoll) actions.ActivePlayerTurn {
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 2,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 3,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 4,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 5,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 6,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 7,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 8,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 9,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 10,
	})
	_ = playerBoard.MakeMove(actions.Move{
		RowColor:   actions.RowColorRed,
		CellNumber: 11,
	})
	return actions.ActivePlayerTurn{}
}

func (b BadActorPlayer) PromptInactivePlayerTurn(playerBoard board.Board, diceRoll actions.DiceRoll) actions.InactivePlayerTurn {
	//TODO implement me
	panic("implement me")
}

func (b BadActorPlayer) InformSuccessfulTurn(updatedBoard board.Board) {
	//TODO implement me
	panic("implement me")
}

func (b BadActorPlayer) InformOfOpponentMove(playerID PlayerID, move actions.Move) {
	//TODO implement me
	panic("implement me")
}

func (b BadActorPlayer) InformRowLocked(color actions.RowColor) {
	//TODO implement me
	panic("implement me")
}

func (b BadActorPlayer) InformWin() {
	//TODO implement me
	panic("implement me")
}

func (b BadActorPlayer) InformLoss(winnerID PlayerID) {
	//TODO implement me
	panic("implement me")
}

package player

import (
	"net"
	"qwixx/board"
)

var _ Player = &httpPlayer{}

type httpPlayer struct {
	name       string
	connection net.Conn
}

func (hp *httpPlayer) InformOfID(playerID PlayerID) {
	//TODO implement me
	panic("implement me")
}

func (hp *httpPlayer) InformOfPlayOrder(playerIDs []PlayerID) {
	//TODO implement me
	panic("implement me")
}

func (hp *httpPlayer) InformOfOpponentMove(playerID PlayerID, move board.Move) {
	//TODO implement me
	panic("implement me")
}

func (hp *httpPlayer) InformRowLocked(color board.RowColor) {
	//TODO implement me
	panic("implement me")
}

func (hp *httpPlayer) InformGameOver(winnerID PlayerID) {
	//TODO implement me
	panic("implement me")
}

func (hp *httpPlayer) PromptMove(diceRoll board.DiceRoll) board.Move {
	//TODO implement me
	panic("implement me")
}

func NewHttpPlayer(name string, connection net.Conn) Player {
	return &httpPlayer{
		name:       name,
		connection: connection,
	}
}

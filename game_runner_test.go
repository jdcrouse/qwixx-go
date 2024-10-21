package qwixx

import (
	"qwixx/player"
	"testing"
)

func TestRunGame(t *testing.T) {
	alice := player.NewComputerPlayer("alice")
	bob := player.NewComputerPlayer("bob")
	charlie := player.NewComputerPlayer("charlie")
	players := []player.Player{alice, bob, charlie}
	runner := NewGameRunner(players)
	runner.RunGame()
}

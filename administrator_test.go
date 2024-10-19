package qwixx

import (
	"github.com/stretchr/testify/require"
	"qwixx/player"
	"testing"
)

func TestAdministrator_CreateGame(t *testing.T) {
	admin := NewAdministrator()
	require.Empty(t, admin.lobbies)

	player1 := player.NewComputerPlayer()
	game1ID := admin.CreateGame(player1)
	require.Len(t, admin.lobbies, 1)
	require.Len(t, admin.lobbies[game1ID], 1)

	player2 := player.NewComputerPlayer()
	game2ID := admin.CreateGame(player2)
	require.Len(t, admin.lobbies, 2)
	require.Len(t, admin.lobbies[game2ID], 1)

	player3 := player.NewComputerPlayer()
	game3ID := admin.CreateGame(player3)
	require.Len(t, admin.lobbies, 3)
	require.Len(t, admin.lobbies[game3ID], 1)
}

func TestAdministrator_JoinGame(t *testing.T) {
	admin := NewAdministrator()
	require.Empty(t, admin.lobbies)

	player1 := player.NewComputerPlayer()
	game1ID := admin.CreateGame(player1)
	require.Len(t, admin.lobbies, 1)
	require.Len(t, admin.lobbies[game1ID], 1)

	player2 := player.NewComputerPlayer()
	admin.JoinGame(game1ID, player2)
	require.Len(t, admin.lobbies, 1)
	require.Len(t, admin.lobbies[game1ID], 2)

	player3 := player.NewComputerPlayer()
	admin.JoinGame(game1ID, player3)
	require.Len(t, admin.lobbies, 1)
	require.Len(t, admin.lobbies[game1ID], 3)
}

func TestAdministrator_StartGame(t *testing.T) {
	admin := NewAdministrator()
	require.Empty(t, admin.lobbies)

	player1 := player.NewComputerPlayer()
	game1ID := admin.CreateGame(player1)

	player2 := player.NewComputerPlayer()
	admin.JoinGame(game1ID, player2)

	player3 := player.NewComputerPlayer()
	admin.JoinGame(game1ID, player3)

	require.Len(t, admin.lobbies, 1)
	require.Len(t, admin.lobbies[game1ID], 3)

	admin.StartGame(game1ID)
	require.Len(t, admin.lobbies, 0)
}

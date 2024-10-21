package qwixx

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"qwixx/actions"
	"qwixx/board"
	"qwixx/player"
	"qwixx/rule_checker"
	"strings"
)

type playerWithID struct {
	ID     player.PlayerID
	player player.Player
}

type GameRunner interface {
	RunGame()
}

type gameRunnerImpl struct {
	playOrder []playerWithID
	boards    map[player.PlayerID]board.Board
	locks     map[actions.RowColor]bool
	penalties map[player.PlayerID]int
}

func NewGameRunner(players []player.Player) GameRunner {
	playOrder := establishPlayOrder(players)
	return &gameRunnerImpl{
		playOrder: playOrder,
		boards:    initializeBoards(playOrder),
		locks:     make(map[actions.RowColor]bool),
		penalties: make(map[player.PlayerID]int),
	}
}

func (gr *gameRunnerImpl) RunGame() {
	gr.notifyPlayersOfPlayOrder()

	turnCount := 0
	for !gr.isGameWon() {
		if turnCount > 1000 { // guard against eternal games
			break
		}
		currentPlayer := gr.playOrder[turnCount%len(gr.playOrder)]
		gr.runSingleTurn(currentPlayer)
		turnCount++
	}
	gr.endGame()
}

// establishPlayOrder establishes the play order of a game comprised of the given list of players
func establishPlayOrder(players []player.Player) []playerWithID {
	playOrder := make([]playerWithID, 0, len(players))
	for _, p := range players {
		playOrder = append(playOrder, playerWithID{ID: player.PlayerID(uuid.New().String()), player: p})
	}
	rand.Shuffle(len(playOrder), func(i, j int) {
		playOrder[i], playOrder[j] = playOrder[j], playOrder[i]
	})
	orderNames := make([]string, 0, len(playOrder))
	for _, p := range playOrder {
		orderNames = append(orderNames, p.player.GetName())
	}
	fmt.Println(fmt.Sprintf("play order: %v", strings.Join(orderNames, ", ")))
	return playOrder
}

func initializeBoards(playOrder []playerWithID) map[player.PlayerID]board.Board {
	boards := make(map[player.PlayerID]board.Board, len(playOrder))
	for _, p := range playOrder {
		boards[p.ID] = board.NewGameBoard()
	}
	return boards
}

func (gr *gameRunnerImpl) notifyPlayersOfPlayOrder() {
	orderIDs := make([]player.PlayerID, 0, len(gr.playOrder))
	for _, p := range gr.playOrder {
		orderIDs = append(orderIDs, p.ID)
	}
	for _, p := range gr.playOrder {
		p.player.InformOfPlayOrder(orderIDs)
	}
}

// isGameWon determines if the currently running game is won
// a game is won if either
// - two rows are locked
// - a player has taken four penalties
func (gr *gameRunnerImpl) isGameWon() bool {
	anyPlayerHasFourPenalties := func() bool {
		for _, count := range gr.penalties {
			if count >= 4 {
				return true
			}
		}
		return false
	}
	return len(gr.locks) >= 2 || anyPlayerHasFourPenalties()
}

func (gr *gameRunnerImpl) runSingleTurn(currentPlayer playerWithID) {
	// Each turn, there is one active player and the rest of the players are inactive.
	// all six dice are rolled (two white and one of each row color)
	// the active player can cross off a cell in any color row with the sum of the white dice
	// the active player can then cross off a cell in a color row with the sum of that color die and one white die
	// the active player must cross off a cell with the sum of the white die before they cross off the sum of a color die.
	// they can choose to cross of only the white die sum or a color die sum if they desire
	// if the active player cannot cross of any cells or does not want to cross off any cells, they must take a penalty
	// each inactive player can cross off a cell in any color row with the sum of the white dice as well, if they like.
	// they cannot do anything with the color dice when they are not the active player, and they do not need to take a penalty if they do not make a move.

	fmt.Println(fmt.Sprintf("it's player %v's turn", currentPlayer.player.GetName()))
	diceRoll := actions.RollQwixxDice()
	fmt.Println(fmt.Sprintf("the white dice rolled were %v and %v", diceRoll.White1, diceRoll.White2))
	fmt.Println(
		fmt.Sprintf(
			"the color dice rolled were red:%v, yellow:%v, green:%v, and blue:%v",
			diceRoll.Red,
			diceRoll.Yellow,
			diceRoll.Green,
			diceRoll.Blue,
		),
	)
	playerBoard := gr.boards[currentPlayer.ID]

	for try := 0; try < 3; try++ {
		fmt.Println(fmt.Sprintf("%v's board before:", currentPlayer.player.GetName()))
		fmt.Println(playerBoard.Print())
		whiteDiceMove, colorDiceMove, takePenalty := currentPlayer.player.PromptActivePlayerMove(playerBoard, diceRoll)
		if whiteDiceMove == nil && colorDiceMove == nil && !takePenalty {
			fmt.Println("player has to make SOME move or take a penalty")
			continue
		}

		if takePenalty {
			gr.penalties[currentPlayer.ID] += 1
			fmt.Println(fmt.Sprintf("player %v took a penalty (they have %v penalties)", currentPlayer.player.GetName(), gr.penalties[currentPlayer.ID]))
			return
		}

		if whiteDiceMove != nil {
			if !rule_checker.WhiteDiceMoveIsValidForBoard(playerBoard, diceRoll, *whiteDiceMove) {
				fmt.Println(fmt.Sprintf("player %v made an invalid move based on the dice rolled %v", currentPlayer.player.GetName(), whiteDiceMove.String()))
				continue
			} else {
				err := playerBoard.MakeMove(*whiteDiceMove)
				if err != nil {
					fmt.Println(fmt.Sprintf("player %v made an invalid move %v", currentPlayer.player.GetName(), whiteDiceMove.String()))
				} else {
					fmt.Println(fmt.Sprintf("player %v made a valid move %v", currentPlayer.player.GetName(), whiteDiceMove.String()))
					fmt.Println(fmt.Sprintf("%v's board after:", currentPlayer.player.GetName()))
					fmt.Println(playerBoard.Print())
				}
			}
		}

		if colorDiceMove != nil {
			if !rule_checker.ColorDiceMoveIsValidForBoard(playerBoard, diceRoll, *colorDiceMove) {
				fmt.Println(fmt.Sprintf("player %v made an invalid move based on the dice rolled %v", currentPlayer.player.GetName(), colorDiceMove.String()))
				continue
			} else {
				err := playerBoard.MakeMove(*colorDiceMove)
				if err != nil {
					fmt.Println(fmt.Sprintf("player %v made an invalid move %v", currentPlayer.player.GetName(), colorDiceMove.String()))
				} else {
					fmt.Println(fmt.Sprintf("player %v made a valid move %v", currentPlayer.player.GetName(), colorDiceMove.String()))
					fmt.Println(fmt.Sprintf("%v's board after:", currentPlayer.player.GetName()))
					fmt.Println(playerBoard.Print())
				}
			}
		}

	}

	// three invalid attempts in one turn forces a penalty
	gr.penalties[currentPlayer.ID] += 1
	return

	// TODO for all other players, prompt inactive player move (just the white dice)
}

func (gr *gameRunnerImpl) endGame() {
	// TODO broadcast end game
	fmt.Println("GAME OVERRR")
}

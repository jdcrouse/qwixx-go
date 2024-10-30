package game

import (
	"fmt"
	"math/rand"
	"qwixx/internal/game/actions"
	"qwixx/internal/game/board"
	"qwixx/internal/game/player"
	"qwixx/internal/game/rule_checker"
	"strings"

	"github.com/google/uuid"
)

type GameRunner interface {
	RunGame()
}

type gameRunnerImpl struct {
	playersByID map[player.PlayerID]player.Player
	boards      map[player.PlayerID]board.Board
	penalties   map[player.PlayerID]int
	locks       map[actions.RowColor]bool
}

func NewGameRunner(players []player.Player) GameRunner {
	playersByID := makePlayersByID(players)
	return &gameRunnerImpl{
		playersByID: playersByID,
		boards:      initializeBoards(playersByID),
		penalties:   make(map[player.PlayerID]int),
		locks:       make(map[actions.RowColor]bool),
	}
}

func (gr *gameRunnerImpl) RunGame() {
	playOrder := establishPlayOrder(gr.playersByID)
	gr.notifyPlayersOfPlayOrder(playOrder)

	turnCount := 0
	for !gr.isGameWon() {
		if turnCount > 1000 { // guard against eternal games
			break
		}
		currentPlayer := playOrder[turnCount%len(playOrder)]
		err := gr.runSingleTurn(currentPlayer)
		if err != nil {
			// TODO do something better
			fmt.Printf("error: %v", err.Error())
		}
		turnCount++
	}
	gr.endGame()
}

func makePlayersByID(players []player.Player) map[player.PlayerID]player.Player {
	playersByID := make(map[player.PlayerID]player.Player, len(players))
	for _, pl := range players {
		id := player.PlayerID(uuid.New().String())
		playersByID[id] = pl
	}
	return playersByID
}

// establishPlayOrder establishes the play order of a game comprised of the given list of players
func establishPlayOrder(playersByID map[player.PlayerID]player.Player) []player.PlayerID {
	playOrder := make([]player.PlayerID, 0, len(playersByID))
	for playerID, _ := range playersByID {
		playOrder = append(playOrder, playerID)
	}
	rand.Shuffle(len(playOrder), func(i, j int) {
		playOrder[i], playOrder[j] = playOrder[j], playOrder[i]
	})
	return playOrder
}

func initializeBoards(playOrder map[player.PlayerID]player.Player) map[player.PlayerID]board.Board {
	boards := make(map[player.PlayerID]board.Board, len(playOrder))
	for playerID, _ := range playOrder {
		boards[playerID] = board.NewGameBoard()
	}
	return boards
}

func (gr *gameRunnerImpl) notifyPlayersOfPlayOrder(playOrder []player.PlayerID) {
	orderNames := make([]string, 0, len(playOrder))
	for _, playerID := range playOrder {
		orderNames = append(orderNames, gr.playersByID[playerID].GetName())
	}
	fmt.Printf("play order: %v\n", strings.Join(orderNames, ", "))

	for _, p := range gr.playersByID {
		p.InformOfPlayOrder(orderNames)
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

func (gr *gameRunnerImpl) runSingleTurn(currentPlayerID player.PlayerID) error {
	// Each turn, there is one active player and the rest of the players are inactive.
	// all six dice are rolled (two white and one of each row color)
	// the active player can cross off a cell in any color row with the sum of the white dice
	// the active player can then cross off a cell in a color row with the sum of that color die and one white die
	// the active player must cross off a cell with the sum of the white die before they cross off the sum of a color die.
	// they can choose to cross of only the white die sum or a color die sum if they desire
	// if the active player cannot cross of any cells or does not want to cross off any cells, they must take a penalty
	// each inactive player can cross off a cell in any color row with the sum of the white dice as well, if they like.
	// they cannot do anything with the color dice when they are not the active player, and they do not need to take a penalty if they do not make a move.

	currentPlayer := gr.playersByID[currentPlayerID]
	fmt.Printf("it's player %v's turn\n", currentPlayer.GetName())

	diceRoll := actions.RollQwixxDice()
	printDiceRoll(diceRoll)

	currentPlayerBoard := gr.boards[currentPlayerID]

	// pass another copy so any mutations in prompting don't affect the board we're going to apply real changes to
	activePlayerTurn := promptActivePlayerTurn(currentPlayer, currentPlayerBoard.Copy(), diceRoll)

	if isActiveTurnPenalty(activePlayerTurn) {
		gr.penalties[currentPlayerID] += 1
		printPenalty(currentPlayer.GetName(), gr.penalties[currentPlayerID])
	} else {

		updatedBoard, err := board.ApplyActivePlayerTurn(currentPlayerBoard.Copy(), activePlayerTurn)
		if err != nil {
			return err
		}

		gr.boards[currentPlayerID] = updatedBoard
	}

	for playerID, pl := range gr.playersByID {
		if playerID != currentPlayerID {
			inactivePlayerBoard := gr.boards[playerID]
			proposedTurn := promptInactivePlayerTurn(pl, inactivePlayerBoard, diceRoll)
			// player can elect to do nothing without a penalty if they are not the active player
			// so only do something if they provided a move
			if proposedTurn.WhiteDiceMove != nil {
				err := inactivePlayerBoard.MakeMove(*proposedTurn.WhiteDiceMove)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// promptActivePlayerTurn prompts a player three times for their active player turn, where they can:
// 1. make a move with the sum of the two white dice
// 2. make a move with the sum of one white die and one color die
// OR
// take a penalty
//
// the returned turn has been guaranteed to be valid for the copy of the board they were given
func promptActivePlayerTurn(
	currentPlayer player.Player,
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
) actions.ActivePlayerTurn {
	for try := 0; try < 3; try++ {
		printPlayerBoard(currentPlayer.GetName(), playerBoard)

		// copy the board so the player can't manipulate it
		proposedTurn := currentPlayer.PromptActivePlayerTurn(playerBoard.Copy(), diceRoll)

		// copy the board so validity checking does not mutate the original board if something was invalid
		if isActiveTurnValid(playerBoard.Copy(), diceRoll, proposedTurn) {
			printValidTurn(currentPlayer.GetName(), proposedTurn.String())
			printPlayerBoard(currentPlayer.GetName(), playerBoard)
			return proposedTurn
		} else {
			printInvalidTurn(currentPlayer.GetName(), proposedTurn.String())
		}
	}

	// three invalid attempts in one turn forces a penalty
	return actions.ActivePlayerTurn{}

}

// promptInactivePlayerTurn prompts a player for their inactive player turn,
// where they can make a move with the sum of the two white dice
// the returned turn has been guaranteed to be valid for the copy of the board they were given
func promptInactivePlayerTurn(
	currentPlayer player.Player,
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
) actions.InactivePlayerTurn {
	for try := 0; try < 3; try++ {
		proposedTurn := currentPlayer.PromptInactivePlayerTurn(playerBoard, diceRoll)
		if isInactiveTurnValid(playerBoard, diceRoll, proposedTurn) {
			return proposedTurn
		}
	}

	// three invalid attempts in one turn forces a no-op
	return actions.InactivePlayerTurn{}
}

func isActiveTurnValid(
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
	activePlayerTurn actions.ActivePlayerTurn,
) bool {
	if isActiveTurnPenalty(activePlayerTurn) {
		return true
	}
	if activePlayerTurn.WhiteDiceMove == nil {
		return isColorDiceMoveValid(playerBoard, diceRoll, *activePlayerTurn.ColorDiceMove)
	}
	if activePlayerTurn.ColorDiceMove == nil {
		return isWhiteDiceMoveValid(playerBoard, diceRoll, *activePlayerTurn.WhiteDiceMove)
	}

	return isWhiteDiceMoveValid(
		playerBoard, diceRoll, *activePlayerTurn.WhiteDiceMove,
	) && isColorDiceMoveValid(
		playerBoard, diceRoll, *activePlayerTurn.ColorDiceMove,
	)
}

func isInactiveTurnValid(
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
	inactivePlayerTurn actions.InactivePlayerTurn,
) bool {
	if inactivePlayerTurn.WhiteDiceMove == nil {
		return true
	}
	return isWhiteDiceMoveValid(playerBoard, diceRoll, *inactivePlayerTurn.WhiteDiceMove)
}

// isWhiteDiceMoveValid determines if the rulechecker says the given color dice move is valid,
// and the given board will allow the move to be played
func isWhiteDiceMoveValid(
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
	move actions.Move,
) bool {
	if !rule_checker.WhiteDiceMoveIsValidForBoard(playerBoard, diceRoll, move) {
		return false
	}

	return playerBoard.MakeMove(move) == nil
}

// isColorDiceMoveValid determines if the rulechecker says the given color dice move is valid,
// and the given board will allow the move to be played
func isColorDiceMoveValid(
	playerBoard board.Board,
	diceRoll actions.DiceRoll,
	move actions.Move,
) bool {
	if !rule_checker.ColorDiceMoveIsValidForBoard(playerBoard, diceRoll, move) {
		return false
	}
	return playerBoard.MakeMove(move) == nil
}

// isActiveTurnPenalty determines if the given turn represents a penalty,
// which is the case when both moves it contains are nil
func isActiveTurnPenalty(activePlayerTurn actions.ActivePlayerTurn) bool {
	return activePlayerTurn.WhiteDiceMove == nil && activePlayerTurn.ColorDiceMove == nil
}

func printDiceRoll(diceRoll actions.DiceRoll) {
	fmt.Printf("the white dice rolled were %v and %v\n", diceRoll.White1, diceRoll.White2)
	fmt.Printf(
		"the color dice rolled were red:%v, yellow:%v, green:%v, and blue:%v\n",
		diceRoll.Red,
		diceRoll.Yellow,
		diceRoll.Green,
		diceRoll.Blue,
	)
}

func printValidTurn(playerName string, validTurnString string) {
	fmt.Printf("player %v played a valid turn: %v\n", playerName, validTurnString)
}

func printInvalidTurn(playerName string, invalidTurnString string) {
	fmt.Printf("player %v played an invalid turn: %v\n", playerName, invalidTurnString)
}

func printPlayerBoard(playerName string, playerBoard board.Board) {
	fmt.Printf("%v's board:\n", playerName)
	fmt.Println(playerBoard.Print())
}

func printPenalty(playerName string, penaltyCount int) {
	fmt.Printf("player %v took a penalty (they have %v penalties)\n", playerName, penaltyCount)
}

func (gr *gameRunnerImpl) endGame() {
	// TODO handle ties
	var winnerID player.PlayerID
	var highScore int
	for playerID, playerBoard := range gr.boards {
		score := playerBoard.CalculateScore()
		if score > highScore {
			winnerID = playerID
		}
	}
	for playerID, pl := range gr.playersByID {
		if playerID == winnerID {
			pl.InformWin()
		} else {
			pl.InformLoss(winnerID)
		}
	}
	fmt.Println("GAME OVERRR")
}

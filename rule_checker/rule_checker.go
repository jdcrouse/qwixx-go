package rule_checker

import (
	"qwixx/actions"
	"qwixx/board"
	"slices"
)

// a white dice move is valid if:
//  1. the board says the move is valid
//     a. the cell number of the given row is not already crossed off
//     b. there are no crossed off cells to the right of the given cell number in the given row
//     c. the given row is not already locked
//     d. if the given cell number is the rightmost cell, there are already five other cells crossed off in that row
//  2. the proposed cell number of the move matches the sum of the two white dice from the dice roll
func WhiteDiceMoveIsValidForBoard(playerBoard board.Board, diceRoll actions.DiceRoll, proposedMove actions.Move) bool {
	// TODO do other rule checking here rather than in board? or both? idk
	boardSaysValid, _ := playerBoard.IsMoveValid(proposedMove)
	possibleMoves := DeterminePossibleWhiteDiceMoves(diceRoll)
	return boardSaysValid && slices.Contains(possibleMoves, proposedMove)
}

func ColorDiceMoveIsValidForBoard(playerBoard board.Board, diceRoll actions.DiceRoll, proposedMove actions.Move) bool {
	// TODO do other rule checking here rather than in board? or both? idk
	boardSaysValid, _ := playerBoard.IsMoveValid(proposedMove)
	possibleMoves := DeterminePossibleColorDiceMoves(diceRoll)
	return boardSaysValid && slices.Contains(possibleMoves, proposedMove)
}

// DeterminePossibleWhiteDiceMoves determines the possible moves that can be made based on only the white dice from the given dice roll.
// This does not take into account the state of any board, just the moves based on the sums from the rolled dice.
// The sum of the white dice can be played on any row, generally.
func DeterminePossibleWhiteDiceMoves(diceRoll actions.DiceRoll) []actions.Move {
	sums := determineDiceRollSums(diceRoll)
	return []actions.Move{
		actions.NewMove(actions.RowColorRed, sums.White),
		actions.NewMove(actions.RowColorYellow, sums.White),
		actions.NewMove(actions.RowColorGreen, sums.White),
		actions.NewMove(actions.RowColorBlue, sums.White),
	}
}

// DeterminePossibleColorDiceMoves determines the possible moves that can be made based on the color dice from the given dice roll.
// This does not take into account the state of any board, just the moves based on the sums from the rolled dice.
// The sum of either white die and one color die can be played on the row with that die's color,
// so you have two possible combinations of (white, color) you could play on each color row.
func DeterminePossibleColorDiceMoves(diceRoll actions.DiceRoll) []actions.Move {
	sums := determineDiceRollSums(diceRoll)
	return []actions.Move{
		actions.NewMove(actions.RowColorRed, sums.Red1),
		actions.NewMove(actions.RowColorRed, sums.Red2),
		actions.NewMove(actions.RowColorYellow, sums.Yellow1),
		actions.NewMove(actions.RowColorYellow, sums.Yellow2),
		actions.NewMove(actions.RowColorGreen, sums.Green1),
		actions.NewMove(actions.RowColorGreen, sums.Green2),
		actions.NewMove(actions.RowColorBlue, sums.Blue1),
		actions.NewMove(actions.RowColorBlue, sums.Blue2),
	}
}

// diceRollSums represents the sums from a dice roll based on the Qwixx rules of summing the different dice colors
// - The white dice must be summed together to get the white dice sum
// - Each other color of die can be summed with either of the individual white dice to get the two possible sums for that color
// Thus there is one white dice sum and two sums per other color of dice, for a total of 7 numbers
type diceRollSums struct {
	// White is the sum of the two white dice from the roll
	White int
	// Red1 is the sum of the red die and the first white die
	Red1 int
	// Red2 is the sum of the red die and the second white die
	Red2 int
	// Yellow1 is the sum of the yellow die and the first white die
	Yellow1 int
	// Yellow2 is the sum of the yellow die and the second white die
	Yellow2 int
	// Green1 is the sum of the green die and the first white die
	Green1 int
	// Green2 is the sum of the green die and the second white die
	Green2 int
	// Blue1 is the sum of the blue die and the first white die
	Blue1 int
	// Blue2 is the sum of the blue die and the second white die
	Blue2 int
}

func determineDiceRollSums(diceRoll actions.DiceRoll) diceRollSums {
	return diceRollSums{
		White:   diceRoll.White1 + diceRoll.White2,
		Red1:    diceRoll.White1 + diceRoll.Red,
		Red2:    diceRoll.White2 + diceRoll.Red,
		Yellow1: diceRoll.White1 + diceRoll.Yellow,
		Yellow2: diceRoll.White2 + diceRoll.Yellow,
		Green1:  diceRoll.White1 + diceRoll.Green,
		Green2:  diceRoll.White2 + diceRoll.Green,
		Blue1:   diceRoll.White1 + diceRoll.Blue,
		Blue2:   diceRoll.White2 + diceRoll.Blue,
	}
}

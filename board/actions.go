package board

type RowColor int

const (
	RowColorRed RowColor = iota
	RowColorYellow
	RowColorGreen
	RowColorBlue
)

// a Move represents crossing off the square with the given number on the row with the given color
type Move struct {
	rowColor   RowColor
	cellNumber int
}

func NewMove(rowColor RowColor, cellNumber int) Move {
	return Move{rowColor: rowColor, cellNumber: cellNumber}
}

// DiceRoll represents the roll of the six Qwixx dice, where two are white
// and the other four are one of each row color from the Qwixx board (red, yellow, green, blue)
type DiceRoll struct {
	White1 int
	White2 int
	Red    int
	Blue   int
	Green  int
	Yellow int
}

// DeterminePossibleMoves determines the possible moves that can be made just based on the given dice roll.
// This does not take into account the state of any board, just the moves based on the sums from the rolled dice.
func DeterminePossibleMoves(diceRoll DiceRoll) []Move {
	sums := determineDiceRollSums(diceRoll)
	return []Move{
		NewMove(RowColorRed, sums.White),
		NewMove(RowColorYellow, sums.White),
		NewMove(RowColorGreen, sums.White),
		NewMove(RowColorBlue, sums.White),

		NewMove(RowColorRed, sums.Red1),
		NewMove(RowColorRed, sums.Red2),
		NewMove(RowColorYellow, sums.Yellow1),
		NewMove(RowColorYellow, sums.Yellow2),
		NewMove(RowColorGreen, sums.Green1),
		NewMove(RowColorGreen, sums.Green2),
		NewMove(RowColorBlue, sums.Blue1),
		NewMove(RowColorBlue, sums.Blue2),
	}
}

// diceRollSums represents the sums from a dice roll based on the Qwixx rules of summing the different dice colors
// - The white dice must be summed together to get the white dice sum
// - Each other color of die can be summed with either of the individual white dice to get the two possible sums for that color
// Thus there is one white dice sum and two sums per other color of dice, for a total of 7 numbers
type DiceRollSums struct {
	White   int
	Red1    int
	Red2    int
	Yellow1 int
	Yellow2 int
	Green1  int
	Green2  int
	Blue1   int
	Blue2   int
}

// determineDiceRollSums determines the valid cell numbers that can be crossed off just based on the given dice roll
// - This does not take into account the state of any board, just the possible sums from the dice roll
// - The white dice must be summed together to get the only white dice sum
// - Each other color of die can be summed with either of the individual white dice to get the two possible sums for that color
func determineDiceRollSums(diceRoll DiceRoll) DiceRollSums {
	return DiceRollSums{
		White:   diceRoll.White1 + diceRoll.White2,
		Red1:    diceRoll.White1 + diceRoll.Red,
		Red2:    diceRoll.White2 + diceRoll.Red,
		Yellow1: diceRoll.White1 + diceRoll.Red,
		Yellow2: diceRoll.White2 + diceRoll.Red,
		Green1:  diceRoll.White1 + diceRoll.Red,
		Green2:  diceRoll.White2 + diceRoll.Red,
		Blue1:   diceRoll.White1 + diceRoll.Red,
		Blue2:   diceRoll.White2 + diceRoll.Red,
	}
}

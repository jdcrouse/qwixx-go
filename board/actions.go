package board

type RowColor int

const (
	Red RowColor = iota
	Yellow
	Green
	Blue
)

// a Move represents crossing off the square with the given number on the row with the given color
type Move struct {
	rowColor   RowColor
	cellNumber int
}

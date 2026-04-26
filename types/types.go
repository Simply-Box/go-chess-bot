package types

import "fmt"

// GameState is all relevant information about a chessboard
type GameState struct {
	Board       [][]string
	WhiteToMove bool
	Castling    CastlingRights
	EnPassant   Coord
	HalfMove	int
	FullMove	int
	Results     GameResult
}

// Makes GameResult struct able to be printed
func (result GameResult) String() string {
	switch result {
	case Ongoing:
		return "Ongoing"
	case Checkmate:
		return "Checkmate"
	case Draw:
		return "Draw"
	default:
		return "Ongoing"
	}
}

// Castling rights
type CastlingRights struct {
	WhiteKingside  bool
	WhiteQueenside bool
	BlackKingside  bool
	BlackQueenside bool
}

// Coordinate as two integers
type Coord struct {
	Row int
	Col int
}

// GameResult determines if current game is ongoing, checkmate or a draw
type GameResult int

const (
	Ongoing GameResult = iota
	Checkmate
	Draw
)

// Containing all relevant information about a single chess move
type Move struct {
	FromRow, FromCol int
	ToRow, ToCol     int
	Piece            string
	Capture          string
	Promotion        string
	Castle           string
}

// Returns true if two Move structs are the same
func (m Move) Equal(other Move) bool {
	return m.FromRow == other.FromRow &&
		m.FromCol == other.FromCol &&
		m.ToRow == other.ToRow &&
		m.ToCol == other.ToCol
}

// Makes Move struct able to be printed
func (m Move) String() string {
	toSquare := func(row, col int) string {
		file := string(rune('a' + col))
		rank := 8 - row
		return fmt.Sprintf("%s%d", file, rank)
	}

	from := toSquare(m.FromRow, m.FromCol)
	to := toSquare(m.ToRow, m.ToCol)

	if m.Capture != "" {
		return fmt.Sprintf("%s %s -> %s, Captured: %s", m.Piece, from, to, m.Capture)
	}

	return fmt.Sprintf("%s %s -> %s", m.Piece, from, to)
}

// GameData is a FEN with a name and an evaluation
type GameData struct {
	FEN   string
	Name  string
	Score int
}

// Used instead of Coord as null
var invalidCoord = Coord{Row: -1, Col: -1}
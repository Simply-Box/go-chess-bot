package board

import (
	"fmt"
	"strings"
	"github.com/Simply-Box/go-chess-bot/types"
)

// Updates the GameState according to the input move
func ApplyMove(gs *types.GameState, move types.Move) {
	piece := move.Piece

	// Move piece to destination
	if move.Promotion != "" {
		gs.Board[move.ToRow][move.ToCol] = move.Promotion
	} else {
		gs.Board[move.ToRow][move.ToCol] = piece
	}
	// clear the moved from square
	gs.Board[move.FromRow][move.FromCol] = "."

	// Handle en passant capture
	if piece == "P" && move.ToRow == gs.EnPassant.Row && move.ToCol == gs.EnPassant.Col {
		gs.Board[move.ToRow+1][move.ToCol] = "."
	} else if piece == "p" && move.ToRow == gs.EnPassant.Row && move.ToCol == gs.EnPassant.Col {
		gs.Board[move.ToRow-1][move.ToCol] = "."
	}

	// Revaluate enPassant
	gs.EnPassant = types.Coord{Row:-1, Col:-1}
	if piece == "P" && move.FromRow == 6 && move.ToRow == 4 {
		gs.EnPassant = types.Coord{Row: 5, Col: move.FromCol}
	} else if piece == "p" && move.FromRow == 1 && move.ToRow == 3 {
		gs.EnPassant = types.Coord{Row: 2, Col: move.FromCol}
	}

	// Castling
	if move.Castle != "" {
		switch move.Castle {
		case "K":
			gs.Board[7][5] = "R"
			gs.Board[7][7] = "."
		case "Q":
			gs.Board[7][3] = "R"
			gs.Board[7][0] = "."
		case "k":
			gs.Board[0][5] = "r"
			gs.Board[0][7] = "."
		case "q":
			gs.Board[0][3] = "r"
			gs.Board[0][0] = "."
		}
	}
	// TODO: if no castling; skip
	UpdateCastlingRights(&gs.Castling, move)

	// Halfmove clock
	if piece == "P" || piece == "p" || move.Capture != "." {
		gs.HalfMove = 0
	} else {
		gs.HalfMove++
		if gs.HalfMove == 100 { // this number could be wrong, 99 or 101??
			gs.Results = types.Draw
		}
	}

	// Fullmove counter
	if !gs.WhiteToMove {
		gs.FullMove++
	}

	gs.WhiteToMove = !gs.WhiteToMove
}

// Updates castling rights
func UpdateCastlingRights(cr *types.CastlingRights, move types.Move) {
	piece := move.Piece
	fromRow, fromCol := move.FromRow, move.FromCol
	toRow, toCol := move.ToRow, move.ToCol

	// King moves
	if piece == "K" {
		cr.WhiteKingside = false
		cr.WhiteQueenside = false
	}
	if piece == "k" {
		cr.BlackKingside = false
		cr.BlackQueenside = false
	}

	// Rook moves
	if piece == "R" && fromRow == 7 {
		switch fromCol {
		case 0:
			cr.WhiteQueenside = false
		case 7:
			cr.WhiteKingside = false
		}
	}
	if piece == "r" && fromRow == 0 {
		switch fromCol {
		case 0:
			cr.BlackQueenside = false
		case 7:
			cr.BlackKingside = false
		}
	}

	// Rook captures start squares
	if move.Capture == "R" && toRow == 7 {
		switch toCol {
		case 0:
			cr.WhiteQueenside = false
		case 7:
			cr.WhiteKingside = false
		}
	}
	if move.Capture == "r" && toRow == 0 {
		switch toCol {
		case 0:
			cr.BlackQueenside = false
		case 7:
			cr.BlackKingside = false
		}
	}
}

func CheckHumanMove(gs *types.GameState, allMoves *[]types.Move, input string) error {

	// Input string into a Move struct
	move, err := ParseMove(input, *gs)
	if err != nil {
		return fmt.Errorf("invalid move format: %v", err)
	}

	var isLegal bool
	for _, legalMove := range *allMoves {
		if legalMove.Equal(move) {
			isLegal = true
			move = legalMove
			break
		}
	}
	if !isLegal {
		return fmt.Errorf("invalid move: %s", input)
	}

	ApplyMove(gs, move)
	return nil
}

// Returns the input move and an error if invalid move
func ParseMove(algebraicMove string, gs types.GameState) (types.Move, error) {
	isWhite := gs.WhiteToMove

	if len(algebraicMove) != 4 {
		return types.Move{}, fmt.Errorf("invalid move format, expected 4 chars like 'e2e4'")
	}

	from := algebraicMove[:2]
	to := algebraicMove[2:]

	if !ValidSquare(from) {
		return types.Move{}, fmt.Errorf("invalid from square: %s", from)
	}
	if !ValidSquare(to) {
		return types.Move{}, fmt.Errorf("invalid to square: %s", to)
	}

	fromCoord := AlgebraicToCoords(from)
	toCoord := AlgebraicToCoords(to)

	// set the current piece moving
	piece := gs.Board[fromCoord.Row][fromCoord.Col]

	if piece == "." || IsEnemy(piece, isWhite) {
		return types.Move{}, fmt.Errorf("not your piece at %s", from)
	}

	// sets the captured piece
	capture := gs.Board[toCoord.Row][toCoord.Col]

	move := types.Move{
		FromRow: fromCoord.Row,
		FromCol: fromCoord.Col,
		ToRow:   toCoord.Row,
		ToCol:   toCoord.Col,
		Piece:   piece,
		Capture: capture,
	}

	return move, nil
}

// Returns true if the algebraic square is inside the board
func ValidSquare(s string) bool {
	file := s[0]
	rank := s[1]
	return file >= 'a' && file <= 'h' && rank >= '1' && rank <= '8'
}

// Returns a FEN string
func GetStartingFEN() string {
	gameData := []types.GameData{
		{FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", Name: "Start Position", Score: 0},
		{FEN: "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2", Name: "King's Knight Opening", Score: 0},
		{FEN: "rnbqkbnr/pppp1ppp/8/8/8/8/PPPP1PPp/RNBQKB2 b Kkq - 0 10", Name: "Promotion", Score: 0},
	}
	fen := gameData[0].FEN
	return fen
}

// Converts a FEN string into a GameState struct
func FENtoGameState(fen string) types.GameState {
	board := make([][]string, 8)
	parts := strings.Split(fen, " ")
	position := strings.Split(parts[0], "/")

	for i, row := range position {
		boardRow := make([]string, 0, 8)
		for _, char := range row {
			if char >= '1' && char <= '8' {
				// add empty squares "."
				numEmpty := int(char - '0')
				for range numEmpty {
					boardRow = append(boardRow, ".")
				}
			} else {
				boardRow = append(boardRow, string(char))
			}
		}
		board[i] = boardRow
	}

	// parts[1] is "w" otherwise "b" and sets the bool
	isWhiteToMove := parts[1] == "w"

	castling := types.CastlingRights{WhiteKingside: false, WhiteQueenside: false, BlackKingside: false, BlackQueenside: false}

	if parts[2] != "-" {
		for i := range parts[2] {
			letter := parts[2][i]
			switch letter {
			case 'K':
				castling.WhiteKingside = true
			case 'Q':
				castling.WhiteQueenside = true
			case 'k':
				castling.BlackKingside = true
			case 'q':
				castling.BlackQueenside = true
			}
		}
	}

	enPassantSquare := types.Coord{Row:-1, Col:-1}
	if parts[3] != "-" {
		enPassantSquare = AlgebraicToCoords(parts[3])
	}

	gameState := types.GameState{
		Board: board, WhiteToMove: isWhiteToMove,
		Castling: castling, EnPassant: enPassantSquare,
		HalfMove: 0, FullMove: 1, Results: types.Ongoing,
	}

	return gameState
}

// Is this usued??
func NewGameState() types.GameState {
	castling := types.CastlingRights{WhiteKingside: true, WhiteQueenside: true, BlackKingside: true, BlackQueenside: true}

	enPassantSquare := types.Coord{Row:-1, Col:-1}

	gameState := types.GameState{
		WhiteToMove: true,
		Castling: castling, EnPassant: enPassantSquare, HalfMove: 0,
		FullMove: 1, Results: types.Ongoing,
	}

	return gameState
}

func NewSliceBoard() [][]string {
	board := make([][]string, 8)
	startingString := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
	position := strings.Split(startingString, "/")

	for i, row := range position {
		boardRow := make([]string, 0, 8)
		for _, char := range row {
			if char >= '1' && char <= '8' {
				// add empty squares "."
				numEmpty := int(char - '0')
				for range numEmpty {
					boardRow = append(boardRow, ".")
				}
			} else {
				boardRow = append(boardRow, string(char))
			}
		}
		board[i] = boardRow
	}
	return board
}

// WHat is this function for????????????
func NewBitboardBoard() int {
	return 0
}

// Converts 'e4' into a types.Coord struct
func AlgebraicToCoords(algebraic string) types.Coord {
	row := 8 - int(algebraic[1]-'0') // '3',  8 - 3 = 5
	col := int(algebraic[0] - 'a')   // 'd', 'd' - 'a'= 3
	return types.Coord{Row: row, Col: col}
}

// Prints the board to terminal
func PrintBoard(n [][]string) {
	for row := range 8 {
		for col := range 8 {
			fmt.Printf("%s ", n[row][col])
		}
		fmt.Printf("%d\n", 8-row)
	}
	fmt.Println("a b c d e f g h")
}

// Prints the result of the current GameState
func PrintResult(Result types.GameResult, isWhite bool) {
	switch Result {
	case types.Checkmate:
		if isWhite {
			fmt.Println("Black wins")
		} else {
			fmt.Println("White wins")
		}
	case types.Draw:
		fmt.Println("Draw")
	}
}

// Should be rewritten as class thing
func CheckGameEnd(result types.GameResult) bool {
	switch result {
	case types.Ongoing:
		return false
	case types.Checkmate:
		return true
	case types.Draw:
		return true
	default:
		return false
	}
}

// Returns true if it is an enemy piece
func IsEnemy(piece string, isWhite bool) bool {
	if piece == "." {
		return false
	}
	if isWhite {
		return piece >= "a" && piece <= "z"
	}
	return piece >= "A" && piece <= "Z"
}
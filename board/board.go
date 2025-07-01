package board

import (
	"fmt"
	"strings"
)

// Patrik joins the battle!

//TODO: Simulate MakeMove with a SimulateMove function, only used for checks or something
//TODO: Add checkmate, stalemate, 50 move rule, three same boardstate, halfmove
//TODO: If not in check don't put yourself in check (self pinned pieces)
//TODO: Remove Fullmove and therefore MoveCounters

//TODO: change board into a hexadecimal and introduce bitboards
//TODO: Clean code: var and func names are understandable,
//TODO: Pointers for big structs or things
//TODO: unit test
//TODO: gs as parameter instead of board+iswhite
//TODO: Maybe import math instead of abs function?
//TODO: Remove the parameter isWhite bool from all functions and use gs.isWhiteToMove instead
//TODO: Human promotion creates enemy color queen
//TODO: Make GameState methods instead of functions?
//TODO: Add a Standard Algebraic Notation method to Move

// For function is square attacked:
// 	change so that: return true actually increments var k with one and return k in the end.
// 	check if != 0 istället för true
// 	get legal moves från kung kan vara 2> så måste jag flytta kungen

// GameState is all relevant information about a chessboard
type GameState struct {
	Board       [][]string
	WhiteToMove bool
	Castling    CastlingRights
	EnPassant   Coord
	Counters    MoveCounters
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

// Movecounters for halfmove and fullmove for FEN and draw
type MoveCounters struct {
	HalfMove int
	FullMove int
}

// Containing all relevant information about a single chess move
type Move struct {
	FromRow, FromCol int
	ToRow, ToCol     int
	Piece            string
	Capture          string
	Promotion        string
	Castle           string
}

// GameData is a FEN with a name and an evaluation
type GameData struct {
	FEN   string
	Name  string
	Score int
}

// Used instead of Coord as null
var invalidCoord = Coord{Row: -1, Col: -1}

// Directions used for pieces
var (
	knightDirs = [8][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	bishopDirs = [][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	rookDirs   = [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	queenDirs  = append(bishopDirs, rookDirs...)
)

// Returns all legal moves for a GameState
func GenerateAllMoves(gs *GameState) []Move {
	board := gs.Board
	isWhite := gs.WhiteToMove
	allMoves := []Move{}

	if false { // IsInCheck(gs, isWhite) {

		// calc all danger lines

		// calc all pieces that can't move

		// gen all moves that end up on danger line

		// what happens if several danger lines?

		// gen all king

		// if allMoves is empty; checkmate

		// 1. Try capturing checking piece
		// 2. Try blocking (if attacker is sliding)

		//attackSquare := checkingPieces[0]
		//attackLine := []Coord{attackSquare}

		// Check if it's a sliding piece (rook, bishop, queen)
		//if IsSlidingPiece(gs.Board[attackSquare.Row][attackSquare.Col]) {
		// Calculate the line between the attacker and king
		//attackLine = GetAttackLine(attackSquare, Coord{kingRow, kingCol})
		//}

		//if len(checkingPieces) >= 2 {

		//	return GenerateKingMoves(board, row, col, isWhite, gs)
		//}

	} else {

		for row := range 8 {
			for col := range 8 {
				piece := board[row][col]

				if !IsEnemy(piece, isWhite) && piece != "." {
					allMoves = append(allMoves, GeneratePieceMoves(gs, row, col, piece)...)
				}
			}
		}
	}
	return allMoves
}

// Returns moves from the piece
func GeneratePieceMoves(gs *GameState, row, col int, piece string) []Move {
	board := gs.Board
	var moveList []Move

	switch piece {
	case "P":
		moveList = GeneratePawnMoves(board, row, col, true, gs.EnPassant)
	case "p":
		moveList = GeneratePawnMoves(board, row, col, false, gs.EnPassant)
	case "N":
		moveList = GenerateKnightMoves(board, row, col, true)
	case "n":
		moveList = GenerateKnightMoves(board, row, col, false)
	case "B":
		moveList = GenerateBishopMoves(board, row, col, true)
	case "b":
		moveList = GenerateBishopMoves(board, row, col, false)
	case "R":
		moveList = GenerateRookMoves(board, row, col, true)
	case "r":
		moveList = GenerateRookMoves(board, row, col, false)
	case "Q":
		moveList = GenerateQueenMoves(board, row, col, true)
	case "q":
		moveList = GenerateQueenMoves(board, row, col, false)
	case "K":
		moveList = GenerateKingMoves(board, row, col, true, gs)
	case "k":
		moveList = GenerateKingMoves(board, row, col, false, gs)
	}
	return moveList
}

// Returns moves from the pawn
func GeneratePawnMoves(board [][]string, row, col int, isWhite bool, enPassantSquare Coord) []Move {
	var moves []Move

	piece := board[row][col]
	direction := -1
	startRow := 6

	if !isWhite {
		direction = 1
		startRow = 1
	}

	// Forward one square
	newRow := row + direction
	if newRow >= 0 && newRow < 8 && board[newRow][col] == "." {
		if (isWhite && newRow == 0) || (!isWhite && newRow == 7) {
			if isWhite {
				for _, promo := range []string{"Q", "R", "B", "N"} {
					moves = append(moves, Move{
						FromRow: row, FromCol: col,
						ToRow: newRow, ToCol: col,
						Piece:     piece,
						Promotion: promo,
					})
				}
			} else {
				for _, promo := range []string{"q", "r", "b", "n"} {
					moves = append(moves, Move{
						FromRow: row, FromCol: col,
						ToRow: newRow, ToCol: col,
						Piece:     piece,
						Promotion: promo,
					})
				}
			}
		} else {
			moves = append(moves, Move{
				FromRow: row, FromCol: col,
				ToRow: newRow, ToCol: col,
				Piece: piece,
			})

			// Forward two squares
			if row == startRow && board[row+2*direction][col] == "." {
				moves = append(moves, Move{
					FromRow: row, FromCol: col,
					ToRow: row + 2*direction, ToCol: col,
					Piece: piece,
				})
			}
		}
	}

	// Captures
	for _, colOffset := range []int{-1, 1} {
		newCol := col + colOffset
		if IsOnBoard(newRow, newCol) {
			target := board[newRow][newCol]
			if target != "." && IsEnemy(target, isWhite) {
				if (isWhite && newRow == 0) || (!isWhite && newRow == 7) {
					for _, promo := range []string{"Q", "R", "B", "N"} {
						moves = append(moves, Move{
							FromRow: row, FromCol: col,
							ToRow: newRow, ToCol: newCol,
							Piece:     piece,
							Capture:   target,
							Promotion: promo,
						})
					}
				} else {
					moves = append(moves, Move{
						FromRow: row, FromCol: col,
						ToRow: newRow, ToCol: newCol,
						Piece:   piece,
						Capture: target,
					})
				}
			}
		}
	}

	// EnPassant
	if enPassantSquare != invalidCoord {
		epRow := enPassantSquare.Row
		epCol := enPassantSquare.Col

		if isWhite && row == 3 && epRow == 2 && Abs(col-epCol) == 1 {
			moves = append(moves, Move{
				FromRow: row, FromCol: col,
				ToRow: epRow, ToCol: epCol,
				Piece:   "P",
				Capture: board[row][epCol],
			})
		}

		if !isWhite && row == 4 && epRow == 5 && Abs(col-epCol) == 1 {
			moves = append(moves, Move{
				FromRow: row, FromCol: col,
				ToRow: epRow, ToCol: epCol,
				Piece:   "p",
				Capture: board[row][epCol],
			})
		}
	}

	return moves
}

// Returns moves from the knight
func GenerateKnightMoves(board [][]string, row, col int, isWhite bool) []Move {
	var moves []Move
	piece := board[row][col]

	for _, offset := range knightDirs {
		newRow := row + offset[0]
		newCol := col + offset[1]

		if IsOnBoard(newRow, newCol) {
			target := board[newRow][newCol]
			if target == "." || IsEnemy(target, isWhite) {
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					Piece: piece,
				}
				if target != "." {
					move.Capture = target
				}
				moves = append(moves, move)
			}
		}
	}
	return moves
}

// Returns moves from sliding pieces such as bishop, rook and queen
func GenerateSlidingMoves(board [][]string, row, col int, isWhite bool, directions [][2]int) []Move {
	var moves []Move
	piece := board[row][col]

	for _, dir := range directions {
		newRow, newCol := row+dir[0], col+dir[1]
		for IsOnBoard(newRow, newCol) {
			target := board[newRow][newCol]
			if target == "." {
				moves = append(moves, Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					Piece: piece,
				})
			} else {
				if IsEnemy(target, isWhite) {
					moves = append(moves, Move{
						FromRow: row, FromCol: col,
						ToRow: newRow, ToCol: newCol,
						Piece:   piece,
						Capture: target,
					})
				}
				break
			}
			newRow += dir[0]
			newCol += dir[1]
		}
	}
	return moves
}

// Returns moves from the bishop
func GenerateBishopMoves(board [][]string, row, col int, isWhite bool) []Move {
	return GenerateSlidingMoves(board, row, col, isWhite, bishopDirs)
}

// Returns moves from the rook
func GenerateRookMoves(board [][]string, row, col int, isWhite bool) []Move {
	return GenerateSlidingMoves(board, row, col, isWhite, rookDirs)
}

// Returns moves from the queen
func GenerateQueenMoves(board [][]string, row, col int, isWhite bool) []Move {
	return GenerateSlidingMoves(board, row, col, isWhite, queenDirs)
}

// Returns moves from the king
func GenerateKingMoves(board [][]string, row, col int, isWhite bool, gs *GameState) []Move {
	var moves []Move

	piece := board[row][col]

	for _, dir := range queenDirs {
		newRow, newCol := row+dir[0], col+dir[1]

		if IsOnBoard(newRow, newCol) {
			target := board[newRow][newCol]
			if target == "." || IsEnemy(target, isWhite) {
				capture := ""
				if target != "." {
					capture = target
				}
				if !IsSquareAttacked(board, newRow, newCol, !isWhite) {
					moves = append(moves, Move{
						FromRow: row, FromCol: col,
						ToRow: newRow, ToCol: newCol,
						Piece:   piece,
						Capture: capture,
					})
				}
			}
		}
	}

	cr := gs.Castling

	// Castling
	if isWhite {
		if cr.WhiteKingside {
			if board[7][5] == "." && board[7][6] == "." &&
				!IsSquareAttacked(board, 7, 4, false) &&
				!IsSquareAttacked(board, 7, 5, false) &&
				!IsSquareAttacked(board, 7, 6, false) {
				moves = append(moves, Move{
					FromRow: 7, FromCol: 4,
					ToRow: 7, ToCol: 6,
					Piece:  "K",
					Castle: "K",
				})
			}
		}
		if cr.WhiteQueenside {
			if board[7][3] == "." && board[7][2] == "." && board[7][1] == "." &&
				!IsSquareAttacked(board, 7, 4, false) &&
				!IsSquareAttacked(board, 7, 3, false) &&
				!IsSquareAttacked(board, 7, 2, false) {
				moves = append(moves, Move{
					FromRow: 7, FromCol: 4,
					ToRow: 7, ToCol: 2,
					Piece:  "K",
					Castle: "Q",
				})
			}
		}
	} else {
		if cr.BlackKingside {
			if board[0][5] == "." && board[0][6] == "." &&
				!IsSquareAttacked(board, 0, 4, true) &&
				!IsSquareAttacked(board, 0, 5, true) &&
				!IsSquareAttacked(board, 0, 6, true) {
				moves = append(moves, Move{
					FromRow: 0, FromCol: 4,
					ToRow: 0, ToCol: 6,
					Piece:  "k",
					Castle: "k",
				})
			}
		}
		if cr.BlackQueenside {
			if board[0][3] == "." && board[0][2] == "." && board[0][1] == "." &&
				!IsSquareAttacked(board, 0, 4, true) &&
				!IsSquareAttacked(board, 0, 3, true) &&
				!IsSquareAttacked(board, 0, 2, true) {
				moves = append(moves, Move{
					FromRow: 0, FromCol: 4,
					ToRow: 0, ToCol: 2,
					Piece:  "k",
					Castle: "q",
				})
			}
		}
	}

	return moves
}

// Not done: Returns the number of pieces checking the king
func GetCheckingPieces(gs *GameState) []Coord {
	isWhite := gs.WhiteToMove
	var checkingPieces []Coord
	board := gs.Board

	kingRow, kingCol := GetKing(board, isWhite)

	// Loop over board to find attackers
	for row := range 8 {
		for col := range 8 {
			piece := board[row][col]
			if !IsEnemy(piece, isWhite) {
				continue
			}
			// Generate pseudo-legal moves for this piece only
			moves := []Move{}
			for _, move := range moves {
				if move.ToRow == kingRow && move.ToCol == kingCol {
					checkingPieces = append(checkingPieces, Coord{Row: row, Col: col})
				}
			}
		}
	}

	return checkingPieces
}

// Returns the row and col for your king
func GetKing(board [][]string, isWhite bool) (int, int) {

	king := "K"
	if !isWhite {
		king = "k"
	}

	for row := range 8 {
		for col := range 8 {
			if board[row][col] == king {
				return row, col
			}
		}
	}
	return 0, 0
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

// Returns true if king is in check
func IsInCheck(gs *GameState, isWhite bool) bool {
	kingRow, kingCol := GetKing(gs.Board, isWhite)
	return IsSquareAttacked(gs.Board, kingRow, kingCol, !isWhite)
}

// Returns true if the square is threatened
func IsSquareAttacked(board [][]string, row, col int, isWhite bool) bool {

	pawnDir := 1
	pawn := "p"
	if isWhite {
		pawnDir = -1
		pawn = "P"
	}
	for _, negPos := range []int{-1, 1} {
		newRow, newCol := row+pawnDir, col+negPos
		if IsOnBoard(newRow, newCol) && board[newRow][newCol] == pawn {
			return true
		}
	}

	knight := "N"
	if !isWhite {
		knight = "n"
	}
	for _, offset := range knightDirs {
		newRow, newCol := row+offset[0], col+offset[1]
		if IsOnBoard(newRow, newCol) && board[newRow][newCol] == knight {
			return true
		}
	}

	rook, queen := "R", "Q"
	if !isWhite {
		rook, queen = "r", "q"
	}
	for _, dir := range rookDirs {
		for i := 1; i < 8; i++ {
			newRow, newCol := row+i*dir[0], col+i*dir[1]
			if !IsOnBoard(newRow, newCol) {
				break
			}
			piece := board[newRow][newCol]
			if piece == "." {
				continue
			}
			if piece == rook || piece == queen {
				return true
			}
			break
		}
	}

	bishop := "B"
	if !isWhite {
		bishop = "b"
	}
	for _, dir := range bishopDirs {
		for i := 1; i < 8; i++ {
			newRow, newCol := row+i*dir[0], col+i*dir[1]
			if !IsOnBoard(newRow, newCol) {
				break
			}
			piece := board[newRow][newCol]
			if piece == "." {
				continue
			}
			if piece == bishop || piece == queen {
				return true
			}
			break
		}
	}

	king := "K"
	if !isWhite {
		king = "k"
	}
	for _, dir := range queenDirs {
		newRow, newCol := row+dir[0], col+dir[1]
		if IsOnBoard(newRow, newCol) && board[newRow][newCol] == king {
			return true
		}
	}

	return false
}

// Returns true if row and col values are within 0 to 8
func IsOnBoard(row, col int) bool {
	if row >= 0 && row < 8 && col >= 0 && col < 8 {
		return true
	} else {
		return false
	}
}

// Updates the GameState according to the input move
func ApplyMove(gs *GameState, move Move) {
	piece := move.Piece

	// Move piece to destination
	if move.Promotion != "" {
		gs.Board[move.ToRow][move.ToCol] = move.Promotion
	} else {
		gs.Board[move.ToRow][move.ToCol] = piece
	}
	gs.Board[move.FromRow][move.FromCol] = "."

	// Handle en passant capture
	if piece == "P" && move.ToRow == gs.EnPassant.Row && move.ToCol == gs.EnPassant.Col {
		gs.Board[move.ToRow+1][move.ToCol] = "."
	} else if piece == "p" && move.ToRow == gs.EnPassant.Row && move.ToCol == gs.EnPassant.Col {
		gs.Board[move.ToRow-1][move.ToCol] = "."
	}

	// Revaluate enPassant
	gs.EnPassant = Coord{-1, -1}
	if piece == "P" && move.FromRow == 6 && move.ToRow == 4 {
		gs.EnPassant = Coord{Row: 5, Col: move.FromCol}
	} else if piece == "p" && move.FromRow == 1 && move.ToRow == 3 {
		gs.EnPassant = Coord{Row: 2, Col: move.FromCol}
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
	if piece == "P" || move.Capture != "." {
		gs.Counters.HalfMove = 0
	} else {
		gs.Counters.HalfMove++
	}

	// Fullmove number
	if !gs.WhiteToMove {
		gs.Counters.FullMove++
	}

	gs.WhiteToMove = !gs.WhiteToMove
}

// Updates castling rights
func UpdateCastlingRights(cr *CastlingRights, move Move) {
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

func MakeMove(gs *GameState, allMoves *[]Move, input string) error {

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
func ParseMove(algebraicMove string, gs GameState) (Move, error) {
	isWhite := gs.WhiteToMove

	if len(algebraicMove) != 4 {
		return Move{}, fmt.Errorf("invalid move format, expected 4 chars like 'e2e4'")
	}

	from := algebraicMove[:2]
	to := algebraicMove[2:]

	if !ValidSquare(from) {
		return Move{}, fmt.Errorf("invalid from square: %s", from)
	}
	if !ValidSquare(to) {
		return Move{}, fmt.Errorf("invalid to square: %s", to)
	}

	fromCoord := AlgebraicToCoords(from)
	toCoord := AlgebraicToCoords(to)

	piece := gs.Board[fromCoord.Row][fromCoord.Col]
	if piece == "." || IsEnemy(piece, isWhite) {
		return Move{}, fmt.Errorf("not your piece at %s", from)
	}

	capture := gs.Board[toCoord.Row][toCoord.Col]

	move := Move{
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
	gameData := []GameData{
		{FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", Name: "Start Position", Score: 0},
		{FEN: "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2", Name: "King's Knight Opening", Score: 0},
		{FEN: "rnbqkbnr/pppp1ppp/8/8/8/8/PPPP1PPp/RNBQKB2 b Kkq - 0 10", Name: "Promotion", Score: 0},
	}
	fen := gameData[0].FEN
	return fen
}

// Converts a FEN string into a GameState struct
func FENtoGameState(fen string) GameState {
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

	castling := CastlingRights{WhiteKingside: false, WhiteQueenside: false, BlackKingside: false, BlackQueenside: false}

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

	enPassantSquare := invalidCoord
	if parts[3] != "-" {
		enPassantSquare = AlgebraicToCoords(parts[3])
	}

	counters := MoveCounters{0, 1}

	gameState := GameState{Board: board, WhiteToMove: isWhiteToMove,
		Castling: castling, EnPassant: enPassantSquare, Counters: counters}

	return gameState
}

// Converts 'e4' into a Coord struct
func AlgebraicToCoords(algebraic string) Coord {
	row := 8 - int(algebraic[1]-'0') // '3',  8 - 3 = 5
	col := int(algebraic[0] - 'a')   // 'd', 'd' - 'a'= 3
	return Coord{Row: row, Col: col}
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

// Returns the absolute value
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

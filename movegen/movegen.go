package movegen

import (
	"github.com/Simply-Box/go-chess-bot/types"
	"github.com/Simply-Box/go-chess-bot/board"
)

// Directions used for pieces
var (
	knightDirs = [8][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	bishopDirs = [][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	rookDirs   = [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	queenDirs  = append(bishopDirs, rookDirs...)
)

// Returns all legal moves for a GameState
func GenerateAllMoves(gs *types.GameState) []types.Move {
	current_board := gs.Board
	isWhite := gs.WhiteToMove
	allMoves := []types.Move{}

	// Fancier implementation, maybe faster (unfinished)
	//
	// kRow, kCol := GetKing(gs)
	// if SquareAttackedAmount(gs, kRow, kCol) > 0 {

	// 	if SquareAttackedAmount(gs, kRow, kCol) > 2 {
	// 		allMoves = append(allMoves, GenerateKingMoves(gs, kRow, kCol)...)
	// 		return allMoves
	// 	}

	// 	// Shadow knight on king position
	// 	for _, move := range GenerateKnightMoves(gs, kRow, kCol) {
	// 		if strings.Contains("Nn", move.Capture) {

	// 			if !IsSquareAttacked(gs, move.ToRow, move.ToCol) {
	// 				allMoves = append(allMoves, GenerateKingMoves(gs, kRow, kCol)...)
	// 				return allMoves
	// 			}

	// 		}
	// 	}

	// 	// Shadow Bishop on king position
	// 	for _, move := range GenerateBishopMoves(gs, kRow, kCol) {
	// 		if strings.Contains("Nn", move.Capture) {

	// 		}
	// 	}

	// 	// Shadow Rook on king position
	// 	for _, move := range GenerateBishopMoves(gs, kRow, kCol) {
	// 		if strings.Contains("Nn", move.Capture) {

	// 		}
	// 	}

	// 	// IsInCheck(gs, isWhite) {

	// calc all danger lines

	// calc all pieces that can't move

	// gen all moves that end up on danger line

	// what happens if several danger lines?

	// gen all king

	// if allMoves is empty; checkmate

	// 1. Try capturing checking piece
	// 2. Try blocking (if attacker is sliding)

	//attackSquare := checkingPieces[0]
	//attackLine := []types.Coord{attackSquare}

	// Check if it's a sliding piece (rook, bishop, queen)
	//if IsSlidingPiece(gs.Board[attackSquare.Row][attackSquare.Col]) {
	// Calculate the line between the attacker and king
	//attackLine = GetAttackLine(attackSquare, types.Coord{kingRow, kingCol})
	//}

	//if len(checkingPieces) >= 2 {

	//	return GenerateKingMoves(board, row, col, isWhite, gs)
	//}

	kRow, kCol := GetKing(gs)
	if kRow == -1 && kCol == -1 {
		gs.Results = types.Checkmate
	}
	for row := range 8 {
		for col := range 8 {
			piece := current_board[row][col]

			if IsEnemy(piece, isWhite) || piece == "." {
				continue
			}

			for _, move := range GeneratePieceMoves(gs, row, col, piece) {
				if piece == "K" || piece == "k" {
					kRow, kCol = move.ToRow, move.ToCol
				}

				// Manual Deep copy of GameState.
				// Bad and why??
				newGs := *gs
				ptrNewGs := &newGs

				slice := make([][]string, 8)
				ptrNewGs.Board = slice

				for i := range 8 {
					slice[i] = make([]string, 8)
					for j := range 8 {
						slice[i][j] = gs.Board[i][j]
					}
				}

				// If this move is made, Am I in check? Fixes both Checks and pinned pieces. ??
				board.ApplyMove(ptrNewGs, move)
				ptrNewGs.WhiteToMove = !newGs.WhiteToMove

				if SquareAttackedAmount(ptrNewGs, kRow, kCol) == 0 {
					allMoves = append(allMoves, move)
				}
			}
		}
	}

	// If no legal moves then checkmate or draw
	if len(allMoves) == 0  {
		if IsInCheck(gs) {
			gs.Results = types.Checkmate
		} else {
			gs.Results = types.Draw
		}
	}
	
	return allMoves
}

// Returns moves from the piece
func GeneratePieceMoves(gs *types.GameState, row, col int, piece string) []types.Move {
	var moveList []types.Move

	switch piece {
	case "P":
		moveList = GeneratePawnMoves(gs, row, col)
	case "p":
		moveList = GeneratePawnMoves(gs, row, col)
	case "N":
		moveList = GenerateKnightMoves(gs, row, col)
	case "n":
		moveList = GenerateKnightMoves(gs, row, col)
	case "B":
		moveList = GenerateBishopMoves(gs, row, col)
	case "b":
		moveList = GenerateBishopMoves(gs, row, col)
	case "R":
		moveList = GenerateRookMoves(gs, row, col)
	case "r":
		moveList = GenerateRookMoves(gs, row, col)
	case "Q":
		moveList = GenerateQueenMoves(gs, row, col)
	case "q":
		moveList = GenerateQueenMoves(gs, row, col)
	case "K":
		moveList = GenerateKingMoves(gs, row, col)
	case "k":
		moveList = GenerateKingMoves(gs, row, col)
	}
	return moveList
}

// Returns moves from the pawn
func GeneratePawnMoves(gs *types.GameState, row, col int) []types.Move {
	var moves []types.Move
	current_board := gs.Board
	piece := current_board[row][col]
	isWhite := gs.WhiteToMove
	enPassantSquare := gs.EnPassant

	direction := -1
	startRow := 6

	if !isWhite {
		direction = 1
		startRow = 1
	}

	// Forward one square
	newRow := row + direction
	if newRow >= 0 && newRow < 8 && current_board[newRow][col] == "." {
		if (isWhite && newRow == 0) || (!isWhite && newRow == 7) {
			if isWhite {
				for _, promo := range []string{"Q", "N", "R", "B"} {
					moves = append(moves, types.Move{
						FromRow: row, FromCol: col,
						ToRow: newRow, ToCol: col,
						Piece:     piece,
						Promotion: promo,
					})
				}
			} else {
				for _, promo := range []string{"q", "n", "r", "b"} {
					moves = append(moves, types.Move{
						FromRow: row, FromCol: col,
						ToRow: newRow, ToCol: col,
						Piece:     piece,
						Promotion: promo,
					})
				}
			}
		} else {
			moves = append(moves, types.Move{
				FromRow: row, FromCol: col,
				ToRow: newRow, ToCol: col,
				Piece: piece,
			})

			// Forward two squares
			if row == startRow && current_board[row+2*direction][col] == "." {
				moves = append(moves, types.Move{
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
			target := current_board[newRow][newCol]
			if target != "." && IsEnemy(target, isWhite) {
				if (isWhite && newRow == 0) || (!isWhite && newRow == 7) {
					if isWhite {
						for _, promo := range []string{"Q", "R", "B", "N"} {
							moves = append(moves, types.Move{
								FromRow: row, FromCol: col,
								ToRow: newRow, ToCol: newCol,
								Piece:     piece,
								Capture:   target,
								Promotion: promo,
							})
						}
					} else {
						for _, promo := range []string{"q", "r", "b", "n"} {
							moves = append(moves, types.Move{
								FromRow: row, FromCol: col,
								ToRow: newRow, ToCol: newCol,
								Piece:     piece,
								Capture:   target,
								Promotion: promo,
							})
						}
					}
				} else {
					moves = append(moves, types.Move{
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
	if enPassantSquare != (types.Coord{Row:-1, Col:-1}) {
		epRow := enPassantSquare.Row
		epCol := enPassantSquare.Col

		if isWhite && row == 3 && epRow == 2 && Abs(col-epCol) == 1 {
			moves = append(moves, types.Move{
				FromRow: row, FromCol: col,
				ToRow: epRow, ToCol: epCol,
				Piece:   "P",
				Capture: current_board[row][epCol],
			})
		}

		if !isWhite && row == 4 && epRow == 5 && Abs(col-epCol) == 1 {
			moves = append(moves, types.Move{
				FromRow: row, FromCol: col,
				ToRow: epRow, ToCol: epCol,
				Piece:   "p",
				Capture: current_board[row][epCol],
			})
		}
	}

	return moves
}

// Returns moves from the knight
func GenerateKnightMoves(gs *types.GameState, row, col int) []types.Move {
	var moves []types.Move
	current_board := gs.Board
	piece := current_board[row][col]
	isWhite := gs.WhiteToMove

	for _, offset := range knightDirs {
		newRow := row + offset[0]
		newCol := col + offset[1]

		if IsOnBoard(newRow, newCol) {
			target := current_board[newRow][newCol]
			if target == "." || IsEnemy(target, isWhite) {
				move := types.Move{
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
func GenerateSlidingMoves(gs *types.GameState, row, col int, directions [][2]int) []types.Move {
	var moves []types.Move
	current_board := gs.Board
	isWhite := gs.WhiteToMove
	piece := current_board[row][col]

	for _, dir := range directions {
		newRow, newCol := row+dir[0], col+dir[1]
		for IsOnBoard(newRow, newCol) {
			target := current_board[newRow][newCol]
			if target == "." {
				moves = append(moves, types.Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					Piece: piece,
				})
			} else {
				if IsEnemy(target, isWhite) {
					moves = append(moves, types.Move{
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
func GenerateBishopMoves(gs *types.GameState, row, col int) []types.Move {
	return GenerateSlidingMoves(gs, row, col, bishopDirs)
}

// Returns moves from the rook
func GenerateRookMoves(gs *types.GameState, row, col int) []types.Move {
	return GenerateSlidingMoves(gs, row, col, rookDirs)
}

// Returns moves from the queen
func GenerateQueenMoves(gs *types.GameState, row, col int) []types.Move {
	return GenerateSlidingMoves(gs, row, col, queenDirs)
}

// Returns moves from the king
func GenerateKingMoves(gs *types.GameState, row, col int) []types.Move {
	isWhite := gs.WhiteToMove
	current_board := gs.Board
	var moves []types.Move

	piece := current_board[row][col]

	for _, dir := range queenDirs {
		newRow, newCol := row+dir[0], col+dir[1]

		if IsOnBoard(newRow, newCol) {
			target := current_board[newRow][newCol]
			if target == "." || IsEnemy(target, isWhite) {
				capture := ""
				if target != "." {
					capture = target
				}
				if !IsSquareAttacked(gs, newRow, newCol) {
					moves = append(moves, types.Move{
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
			if current_board[7][5] == "." && current_board[7][6] == "." &&
				!IsSquareAttacked(gs, 7, 4) &&
				!IsSquareAttacked(gs, 7, 5) &&
				!IsSquareAttacked(gs, 7, 6) {
				moves = append(moves, types.Move{
					FromRow: 7, FromCol: 4,
					ToRow: 7, ToCol: 6,
					Piece:  "K",
					Castle: "K",
				})
			}
		}
		if cr.WhiteQueenside {
			if current_board[7][3] == "." && current_board[7][2] == "." && current_board[7][1] == "." &&
				!IsSquareAttacked(gs, 7, 4) &&
				!IsSquareAttacked(gs, 7, 3) &&
				!IsSquareAttacked(gs, 7, 2) {
				moves = append(moves, types.Move{
					FromRow: 7, FromCol: 4,
					ToRow: 7, ToCol: 2,
					Piece:  "K",
					Castle: "Q",
				})
			}
		}
	} else {
		if cr.BlackKingside {
			if current_board[0][5] == "." && current_board[0][6] == "." &&
				!IsSquareAttacked(gs, 0, 4) &&
				!IsSquareAttacked(gs, 0, 5) &&
				!IsSquareAttacked(gs, 0, 6) {
				moves = append(moves, types.Move{
					FromRow: 0, FromCol: 4,
					ToRow: 0, ToCol: 6,
					Piece:  "k",
					Castle: "k",
				})
			}
		}
		if cr.BlackQueenside {
			if current_board[0][3] == "." && current_board[0][2] == "." && current_board[0][1] == "." &&
				!IsSquareAttacked(gs, 0, 4) &&
				!IsSquareAttacked(gs, 0, 3) &&
				!IsSquareAttacked(gs, 0, 2) {
				moves = append(moves, types.Move{
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
func GetCheckingPieces(gs *types.GameState) []types.Coord {
	isWhite := gs.WhiteToMove
	current_board := gs.Board
	var checkingPieces []types.Coord

	kingRow, kingCol := GetKing(gs)

	// Loop over board to find attackers
	for row := range 8 {
		for col := range 8 {
			piece := current_board[row][col]
			if !IsEnemy(piece, isWhite) {
				continue
			}
			// Generate pseudo-legal moves for this piece only
			moves := []types.Move{}
			for _, move := range moves {
				if move.ToRow == kingRow && move.ToCol == kingCol {
					checkingPieces = append(checkingPieces, types.Coord{Row: row, Col: col})
				}
			}
		}
	}

	return checkingPieces
}

// Returns the row and col for your king
func GetKing(gs *types.GameState) (int, int) {
	current_board := gs.Board
	king := "K"
	if !gs.WhiteToMove {
		king = "k"
	}

	for row := range 8 {
		for col := range 8 {
			if current_board[row][col] == king {
				return row, col
			}
		}
	}
	return -1, -1
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
func IsInCheck(gs *types.GameState) bool {
	kingRow, kingCol := GetKing(gs)
	return IsSquareAttacked(gs, kingRow, kingCol)
}

func SquareAttackedAmount(gs *types.GameState, row, col int) int {
	attackers := 0
	attackerIsWhite := !gs.WhiteToMove
	current_board := gs.Board

	enemyPawnDir := 1
	pawn, knight, bishop, rook, queen, king := "p", "n", "b", "r", "q", "k"

	if attackerIsWhite {
		enemyPawnDir = -1
		pawn, knight, bishop, rook, queen, king = "P", "N", "B", "R", "Q", "K"
	}

	for _, negPos := range []int{-1, 1} {
		newRow, newCol := row+enemyPawnDir, col+negPos
		if IsOnBoard(newRow, newCol) && current_board[newRow][newCol] == pawn {
			attackers = attackers + 1
		}
	}

	for _, offset := range knightDirs {
		newRow, newCol := row+offset[0], col+offset[1]
		if IsOnBoard(newRow, newCol) && current_board[newRow][newCol] == knight {
			attackers = attackers + 1
		}
	}

	for _, dir := range rookDirs {
		for i := 1; i < 8; i++ {
			newRow, newCol := row+i*dir[0], col+i*dir[1]
			if !IsOnBoard(newRow, newCol) {
				break
			}
			piece := current_board[newRow][newCol]
			if piece == "." {
				continue
			}
			if piece == rook || piece == queen {
				attackers = attackers + 1
			}
			break
		}
	}

	for _, dir := range bishopDirs {
		for i := 1; i < 8; i++ {
			newRow, newCol := row+i*dir[0], col+i*dir[1]
			if !IsOnBoard(newRow, newCol) {
				break
			}
			piece := current_board[newRow][newCol]
			if piece == "." {
				continue
			}
			if piece == bishop || piece == queen {
				attackers = attackers + 1
			}
			break
		}
	}

	for _, dir := range queenDirs {
		newRow, newCol := row+dir[0], col+dir[1]
		if IsOnBoard(newRow, newCol) && current_board[newRow][newCol] == king {
			attackers = attackers + 1
		}
	}

	return attackers
}

// Returns true if the square is threatened
func IsSquareAttacked(gs *types.GameState, row, col int) bool {
	attackerIsWhite := !gs.WhiteToMove
	current_board := gs.Board

	enemyPawnDir := 1
	pawn, knight, bishop, rook, queen, king := "p", "n", "b", "r", "q", "k"

	if attackerIsWhite {
		enemyPawnDir = -1
		pawn, knight, bishop, rook, queen, king = "P", "N", "B", "R", "Q", "K"
	}

	for _, negPos := range []int{-1, 1} {
		newRow, newCol := row+enemyPawnDir, col+negPos
		if IsOnBoard(newRow, newCol) && current_board[newRow][newCol] == pawn {
			return true
		}
	}

	for _, offset := range knightDirs {
		newRow, newCol := row+offset[0], col+offset[1]
		if IsOnBoard(newRow, newCol) && current_board[newRow][newCol] == knight {
			return true
		}
	}

	for _, dir := range rookDirs {
		for i := 1; i < 8; i++ {
			newRow, newCol := row+i*dir[0], col+i*dir[1]
			if !IsOnBoard(newRow, newCol) {
				break
			}
			piece := current_board[newRow][newCol]
			if piece == "." {
				continue
			}
			if piece == rook || piece == queen {
				return true
			}
			break
		}
	}

	for _, dir := range bishopDirs {
		for i := 1; i < 8; i++ {
			newRow, newCol := row+i*dir[0], col+i*dir[1]
			if !IsOnBoard(newRow, newCol) {
				break
			}
			piece := current_board[newRow][newCol]
			if piece == "." {
				continue
			}
			if piece == bishop || piece == queen {
				return true
			}
			break
		}
	}

	for _, dir := range queenDirs {
		newRow, newCol := row+dir[0], col+dir[1]
		if IsOnBoard(newRow, newCol) && current_board[newRow][newCol] == king {
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

// Returns the absolute value
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
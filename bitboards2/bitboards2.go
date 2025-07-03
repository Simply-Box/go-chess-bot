package bitboard2
// Working: Implementing bitbords
// 	Changing everything into hexadecimal representation + bitboards
// 	Change enPassant to a hexadecimal (bitboard)
//	Add a performance check through ppef in seperate terminal
//	Add test case to check performance diff

import 	"math/bits"

// Bitboard helper layout
//
//	A1 = bit 0, H1 = bit 7, A8 = bit 56, H8 = bit 63
//
// Stor/liten bokstav?
const (
	A1 uint = iota
	B1;C1;D1;E1;F1;G1;H1
	A2;B2;C2;D2;E2;F2;G2;H2
	A3;B3;C3;D3;E3;F3;G3;H3
	A4;B4;C4;D4;E4;F4;G4;H4
	A5;B5;C5;D5;E5;F5;G5;H5
	A6;B6;C6;D6;E6;F6;G6;H6
	A7;B7;C7;D7;E7;F7;G7;H7
	A8;B8;C8;D8;E8;F8;G8;H8
)

type Bitboard uint64

// Containing all relevant information about a single chess move
type Move struct {
	FromRow, FromCol int
	ToRow, ToCol     int
	Piece            string
	Capture          string
	Promotion        string
	Castle           string
}

// Helper functions
func (b Bitboard) IsSet(sq uint) bool { return (b>>sq)&1 != 0 }
func (b *Bitboard) Set(sq uint)       { *b |= 1 << sq }
func (b *Bitboard) Clear(sq uint)     { *b &^= 1 << sq }

// iterate & pop least‑significant‑bit (LSB)
func PopLSB(b *Bitboard) uint {
	lsb := bits.TrailingZeros64(uint64(*b))
	*b &^= 1 << uint(lsb)
	return uint(lsb)
}

type AllBitboards struct {
	WhitePawns, BlackPawns     Bitboard
	WhiteKnights, BlackKnights Bitboard
	WhiteBishops, BlackBishops Bitboard
	WhiteRooks, BlackRooks     Bitboard
	WhiteQueens, BlackQueens   Bitboard
	WhiteKing, BlackKing       Bitboard
	WhitePieces, BlackPieces   Bitboard
	AllPieces                  Bitboard
}

// Load the starting position
func NewStartPos() *AllBitboards {
	allB := &AllBitboards{}

	for i := range uint(8) {
		allB.WhitePawns.Set(A2 + i)
		allB.BlackPawns.Set(A7 + i)
	}

	allB.WhiteRooks.Set(A1)
	allB.WhiteRooks.Set(H1)
	allB.WhiteKnights.Set(B1)
	allB.WhiteKnights.Set(G1)
	allB.WhiteBishops.Set(C1)
	allB.WhiteBishops.Set(F1)
	allB.WhiteQueens.Set(D1)
	allB.WhiteKing.Set(E1)

	allB.BlackRooks.Set(A8)
	allB.BlackRooks.Set(H8)
	allB.BlackKnights.Set(B8)
	allB.BlackKnights.Set(G8)
	allB.BlackBishops.Set(C8)
	allB.BlackBishops.Set(F8)
	allB.BlackQueens.Set(D8)
	allB.BlackKing.Set(E8)

	allB.updateAllPieces()
	return allB
}

func (allB *AllBitboards) updateAllPieces() {
	allB.WhitePieces = allB.WhitePawns | allB.WhiteKnights | allB.WhiteBishops |
		allB.WhiteRooks | allB.WhiteQueens | allB.WhiteKing
	allB.BlackPieces = allB.BlackPawns | allB.BlackKnights | allB.BlackBishops |
		allB.BlackRooks | allB.BlackQueens | allB.BlackKing
	allB.AllPieces = allB.WhitePieces | allB.BlackPieces
}

// Returns true if row and col values are within 0 to 8
func IsOnBoard(row, col int) bool {
	if row >= 0 && row < 8 && col >= 0 && col < 8 {
		return true
	} else {
		return false
	}
}

func (b *AllBitboards) KnightMoves(isWhite bool) []Move {
	var knights, ownPieces Bitboard
	var pieceChar string
	if isWhite {
		knights, ownPieces = b.WhiteKnights, b.WhitePieces
	} else {
		knights, ownPieces = b.BlackKnights, b.BlackPieces
	}

	var moves []Move
	for knights != 0 {
		from := PopLSB(&knights)
		targets := knightMask[from] & ^ownPieces // blocked by own pieces
		for targets != 0 {
			to := PopLSB(&targets)
			move := Move{
				FromRow: int(from / 8), FromCol: int(from % 8),
				ToRow: int(to / 8), ToCol: int(to % 8),
				Piece: pieceChar,
			}
			moves = append(moves, move)
		}
	}
	return moves
}

var knightMask [64]Bitboard

// Precompute knight moves for every square
func preKnightMoves() {
	for square := range uint(64) {
		row, col := square/8, square%8
		knightDirs := [8][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
			{1, -2}, {1, 2}, {2, -1}, {2, 1}}
		for _, dir := range knightDirs {
			newRow, newCol := int(row)+dir[0], int(col)+dir[1]
			if IsOnBoard(newRow, newCol) {
				knightMask[square] |= 1 << uint(newRow*8+newCol)
			}
		}
	}
}
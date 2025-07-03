package main

// hej

import (
	"fmt"
	"strconv"
)

type Bitboard struct {
	board uint64
}

var EmptyBoard = Bitboard{uint64(0)}

func (bboard Bitboard) Add(bboard2 Bitboard) Bitboard {
	return Bitboard{bboard.board ^ bboard2.board}
}

func (bboard Bitboard) Sub(bboard2 Bitboard) Bitboard {
	return Bitboard{bboard.board &^ bboard2.board}
}

func (bboard Bitboard) And(bboard2 Bitboard) Bitboard {
	return Bitboard{bboard.board & bboard2.board}
}

func (bboard Bitboard) Or(bboard2 Bitboard) Bitboard {
	return Bitboard{bboard.board | bboard2.board}
}

// Directions used for pieces
var (
	knightDirs = [8][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	bishopDirs = [][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	rookDirs   = [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	queenDirs  = append(bishopDirs, rookDirs...)
)

func IsOnBoard(row, col int) bool {
	if row >= 0 && row < 8 && col >= 0 && col < 8 {
		return true
	} else {
		return false
	}
}

func CreatePrimBoards() map[uint8]Bitboard {
	primBoards := make(map[uint8]Bitboard)
	for i := range 64 {
		i := uint8(i)
		primBoards[i] = Bitboard{uint64(1 << i)}
	}
	return primBoards
}

func Visualize(bboard Bitboard) {
	string := ""
	unused := bboard.board
	for j := range 64 {
		if (unused % 2) == 0 {
			string = string + "0"
		} else {
			string = string + "1"
		}
		unused = unused >> 1
		if (j % 8) == 7 {
			string = string + "\n"
		}
	}

	fmt.Println(string)
}

func KnightBoards(mp map[string]Bitboard, prims map[uint8]Bitboard) {
	Name := "n"
	for row := range 8 {
		for col := range 8 {
			Name = Name + strconv.Itoa(row)
			Name = Name + strconv.Itoa(col)

			bboard := EmptyBoard

			for _, move := range knightDirs {
				nRow, nCol := row+move[0], col+move[1]
				if IsOnBoard(nRow, nCol) {
					pos := uint8(nCol + 8*nRow)
					bboard = bboard.Or(prims[pos])
				}
			}

			fmt.Println(Name)
			mp[Name] = bboard
			Name = "n"
		}
	}
}

func RookBoards(mp map[string]Bitboard, prims map[uint8]Bitboard) {
	Name := "r"
	for row := range 8 {
		for col := range 8 {
			Name = Name + strconv.Itoa(row)
			Name = Name + strconv.Itoa(col)

			bboard := EmptyBoard

			for _, move := range rookDirs {
				nRow, nCol := row+move[0], col+move[1]
				for IsOnBoard(nRow, nCol) {
					pos := uint8(nCol + 8*nRow)
					bboard = bboard.Or(prims[pos])
					nRow, nCol = nRow+move[0], nCol+move[1]
				}

			}

			fmt.Println(Name)
			mp[Name] = bboard
			Name = "r"
		}
	}
}

func BishopBoards(mp map[string]Bitboard, prims map[uint8]Bitboard) {
	Name := "b"
	for row := range 8 {
		for col := range 8 {
			Name = Name + strconv.Itoa(row)
			Name = Name + strconv.Itoa(col)

			bboard := EmptyBoard

			for _, move := range bishopDirs {
				nRow, nCol := row+move[0], col+move[1]
				for IsOnBoard(nRow, nCol) {
					pos := uint8(nCol + 8*nRow)
					bboard = bboard.Or(prims[pos])
					nRow, nCol = nRow+move[0], nCol+move[1]
				}

			}

			fmt.Println(Name)
			mp[Name] = bboard
			Name = "b"
		}
	}
}

// Premakes all static boards (piece movement)
func CreateBoards(prims map[uint8]Bitboard) map[string]Bitboard {
	mp := make(map[string]Bitboard)

	KnightBoards(mp, prims)
	RookBoards(mp, prims)
	BishopBoards(mp, prims)

	return mp
}

func main() {
	prims := CreatePrimBoards()

	PieceBoards := CreateBoards(prims)
	Visualize(PieceBoards["b12"])
}

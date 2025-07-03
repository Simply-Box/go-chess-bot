package main

import (
	"fmt"
)

type Bitboard struct {
	board uint64
}

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

func CreatePrimBoards() map[uint8]Bitboard {
	primBoards := make(map[uint8]Bitboard)
	for i := range 64 {
		i := uint8(i)
		primBoards[i] = Bitboard{uint64(1 << (i))}
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

func main() {
	mp := CreatePrimBoards()
	for i := range 64 {
		Visualize(mp[uint8(i)])
	}
	Visualize(mp[uint8(2)].Add(mp[uint8(2)]))
}

// func CreateBoards () map[string]Bitboard {
// 	continue
// }

// go test ./... runs tests in all sub packages
package movegen

import (
	"testing"
	"github.com/Simply-Box/go-chess-bot/board"
)

func BenchmarkMoveGenerationStringSlice(b *testing.B) {
	gs := board.NewGameState()
	gs.Board = board.NewSliceBoard()

	// This loop will run until the benchmark is satisfied
	for b.Loop() {
		_ = GenerateAllMoves(&gs)
	}
}

func BenchmarkMoveGenerationBitboard(b *testing.B) {
	gs := board.NewGameState()
	gs.Board = board.NewSliceBoard() // func NewBitboardBoard and then test it against SliceBoard

	for b.Loop() {
		_ = GenerateAllMoves(&gs)
	}
}

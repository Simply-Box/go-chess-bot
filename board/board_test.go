// go test ./... runs tests in all sub packages
package board

import "testing"

func TestAlgebraicToCoords(t *testing.T) {
	got := AlgebraicToCoords("e4")
	want := Coord{Row: 4, Col: 4}
	if got != want {
		t.Fatalf("e4: %+v, want: %+v", got, want)
	}
}

func BenchmarkMoveGenerationStringSlice(b *testing.B) {
	gs := NewGameState()
	gs.Board = NewSliceBoard()

	// This loop will run until the benchmark is satisfied
	for b.Loop() {
		_ = GenerateAllMoves(&gs)
	}
}

func BenchmarkMoveGenerationBitboard(b *testing.B) {
	gs := NewGameState()
	gs.Board = NewSliceBoard() // func NewBitboardBoard and then test it against SliceBoard

	for b.Loop() {
		_ = GenerateAllMoves(&gs)
	}
}

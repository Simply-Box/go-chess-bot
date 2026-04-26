// go test ./... runs tests in all sub packages
package board

import (
	"testing"

	model "github.com/Simply-Box/go-chess-bot/types"
)

func TestAlgebraicToCoords(t *testing.T) {
	got := AlgebraicToCoords("e4")
	want := model.Coord{Row: 4, Col: 4}
	if got != want {
		t.Fatalf("e4: %+v, want: %+v", got, want)
	}
}

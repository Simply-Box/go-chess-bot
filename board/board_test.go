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

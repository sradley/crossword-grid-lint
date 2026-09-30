package main

import (
	"reflect"
	"testing"
)

func gridFrom(rows ...string) *Grid {
	return &Grid{Cells: rows, Width: len(rows[0]), Height: len(rows)}
}

func TestUncheckedNone(t *testing.T) {
	g := gridFrom(
		"..#..",
		".....",
		"#...#",
		".....",
		"..#..",
	)
	if got := g.Unchecked(); len(got) != 0 {
		t.Errorf("expected no unchecked squares, got %v", got)
	}
}

func TestUncheckedCorner(t *testing.T) {
	// The top-left square is the end of a 3-long across word, but the
	// black square below it leaves no down word.
	g := gridFrom(
		"...",
		"#..",
		"...",
	)
	want := []Cell{{Row: 0, Col: 0}, {Row: 2, Col: 0}}
	if got := g.Unchecked(); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUncheckedIgnoresIsolatedSquare(t *testing.T) {
	g := gridFrom(
		"#.#",
		"###",
		"#.#",
	)
	if got := g.Unchecked(); len(got) != 0 {
		t.Errorf("isolated squares are in no word and shouldn't be listed, got %v", got)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

// Grid is the block/white layout of a crossword, before any letters are
// filled in. '#' marks a black square, '.' marks a white (fillable) square.
type Grid struct {
	Cells  []string
	Width  int
	Height int
}

// LoadGrid reads a grid from a plain text file: one row per line, using
// '#' for black squares and '.' for white squares. Blank lines are skipped
// so a trailing newline at end of file doesn't turn into a phantom row.
func LoadGrid(path string) (*Grid, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rows []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		rows = append(rows, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("%s: empty grid", path)
	}

	width := len(rows[0])
	for i, row := range rows {
		if len(row) != width {
			return nil, fmt.Errorf("%s: row %d has length %d, expected %d (rows must line up)", path, i+1, len(row), width)
		}
		for _, ch := range row {
			if ch != '#' && ch != '.' {
				return nil, fmt.Errorf("%s: row %d contains %q, only '#' and '.' are allowed", path, i+1, ch)
			}
		}
	}

	return &Grid{Cells: rows, Width: width, Height: len(rows)}, nil
}

// isBlack treats anything outside the grid as black too, so a word run
// stops cleanly at the edge without a separate bounds check everywhere.
func (g *Grid) isBlack(r, c int) bool {
	if r < 0 || r >= g.Height || c < 0 || c >= g.Width {
		return true
	}
	return g.Cells[r][c] == '#'
}

// isSymmetric checks standard 180-degree rotational symmetry: every black
// square has a black square in the mirrored position on the other side of
// the grid's center.
func (g *Grid) isSymmetric() bool {
	for r := 0; r < g.Height; r++ {
		for c := 0; c < g.Width; c++ {
			if g.isBlack(r, c) != g.isBlack(g.Height-1-r, g.Width-1-c) {
				return false
			}
		}
	}
	return true
}

// Word is a single across or down entry, identified the way solvers see it:
// a clue number plus a direction. Row and Col are the 0-indexed position of
// the word's first cell.
type Word struct {
	Number    int    `json:"number"`
	Direction string `json:"direction"`
	Row       int    `json:"row"`
	Col       int    `json:"col"`
	Length    int    `json:"length"`
}

// Words numbers the grid the way a crossword puzzle is numbered: scan
// left-to-right, top-to-bottom, and give a cell a number if it starts an
// across word, a down word, or both (sharing the one number).
func (g *Grid) Words() []Word {
	var words []Word
	number := 0
	for r := 0; r < g.Height; r++ {
		for c := 0; c < g.Width; c++ {
			if g.isBlack(r, c) {
				continue
			}
			startsAcross := g.isBlack(r, c-1) && !g.isBlack(r, c+1)
			startsDown := g.isBlack(r-1, c) && !g.isBlack(r+1, c)
			if !startsAcross && !startsDown {
				continue
			}
			number++
			if startsAcross {
				words = append(words, Word{Number: number, Direction: "across", Row: r, Col: c, Length: runLength(g, r, c, 0, 1)})
			}
			if startsDown {
				words = append(words, Word{Number: number, Direction: "down", Row: r, Col: c, Length: runLength(g, r, c, 1, 0)})
			}
		}
	}
	return words
}

// runLength counts consecutive white cells starting at (r, c) and moving
// by (dr, dc) each step, until a black square or the edge is hit.
func runLength(g *Grid, r, c, dr, dc int) int {
	length := 0
	for !g.isBlack(r, c) {
		length++
		r += dr
		c += dc
	}
	return length
}

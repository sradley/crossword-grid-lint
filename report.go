package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// Report is the result of checking one grid. It's the same data whether
// it ends up printed as text or as JSON.
type Report struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	BlackSquares int    `json:"black_squares"`
	WhiteSquares int    `json:"white_squares"`
	Symmetric    bool   `json:"symmetric"`
	MinWordLen   int    `json:"min_word_length"`
	Words        []Word `json:"words"`
	ShortWords   []Word `json:"short_words"`
}

// BuildReport runs every check against g and collects the results. minWordLen
// is the shortest a word is allowed to be before it's flagged; American
// grids conventionally reject anything under 3 letters.
func BuildReport(g *Grid, minWordLen int) *Report {
	black := 0
	for r := 0; r < g.Height; r++ {
		for c := 0; c < g.Width; c++ {
			if g.isBlack(r, c) {
				black++
			}
		}
	}

	words := g.Words()
	var short []Word
	for _, w := range words {
		if w.Length < minWordLen {
			short = append(short, w)
		}
	}

	total := g.Width * g.Height
	return &Report{
		Width:        g.Width,
		Height:       g.Height,
		BlackSquares: black,
		WhiteSquares: total - black,
		Symmetric:    g.isSymmetric(),
		MinWordLen:   minWordLen,
		Words:        words,
		ShortWords:   short,
	}
}

func (rep *Report) PrintHuman(w io.Writer) {
	fmt.Fprintf(w, "grid: %dx%d (%d black, %d white)\n", rep.Width, rep.Height, rep.BlackSquares, rep.WhiteSquares)

	symText := "yes"
	if !rep.Symmetric {
		symText = "no"
	}
	fmt.Fprintf(w, "180-degree rotational symmetry: %s\n", symText)

	across, down := 0, 0
	for _, wd := range rep.Words {
		if wd.Direction == "across" {
			across++
		} else {
			down++
		}
	}
	fmt.Fprintf(w, "words: %d (%d across, %d down)\n", len(rep.Words), across, down)

	if len(rep.ShortWords) == 0 {
		fmt.Fprintf(w, "no words shorter than %d cells\n", rep.MinWordLen)
		return
	}
	fmt.Fprintf(w, "words shorter than %d cells:\n", rep.MinWordLen)
	for _, wd := range rep.ShortWords {
		fmt.Fprintf(w, "  %d-%s at row %d, col %d (length %d)\n", wd.Number, wd.Direction, wd.Row+1, wd.Col+1, wd.Length)
	}
}

func (rep *Report) PrintJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}

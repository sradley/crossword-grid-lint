package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	minWordLen := flag.Int("min-word-length", 3, "shortest word length that doesn't get flagged")
	jsonOut := flag.Bool("json", false, "print the report as JSON instead of plain text")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] <grid-file>\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Checks a crossword grid's block layout: rotational symmetry, word")
		fmt.Fprintln(os.Stderr, "numbering, and word lengths. The grid file is plain text, '#' for")
		fmt.Fprintln(os.Stderr, "black squares and '.' for white squares, one row per line.")
		fmt.Fprintln(os.Stderr, "\nflags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	grid, err := LoadGrid(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "gridlint:", err)
		os.Exit(1)
	}

	report := BuildReport(grid, *minWordLen)

	if *jsonOut {
		if err := report.PrintJSON(os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "gridlint:", err)
			os.Exit(1)
		}
		return
	}
	report.PrintHuman(os.Stdout)
}

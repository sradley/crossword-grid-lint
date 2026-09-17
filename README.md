# gridlint

A command-line checker for crossword grid layouts.

Constructors usually block out the black-and-white pattern of a grid before
writing a single clue. That layout has to follow a few rules that are easy
to get subtly wrong by hand: the black squares should be symmetric under a
180-degree rotation, and every word should be at least three letters long
(most American outlets reject shorter entries). `gridlint` reads a grid
from a plain text file and reports whether it holds up, either as text you
can read or as JSON you can feed into something else.

It only looks at grid structure right now - black/white layout, symmetry,
numbering, word lengths. It doesn't know about letters, clues, or fill.

## Grid file format

One line per row. `#` is a black square, `.` is a white square. Rows must
all be the same length.

```
..#..
.....
#...#
.....
..#..
```

Save that as `sample.txt` (it's also checked in at `testdata/sample.txt`)
and run:

```
$ go run . testdata/sample.txt
grid: 5x5 (4 black, 21 white)
180-degree rotational symmetry: yes
words: 14 (7 across, 7 down)
words shorter than 3 cells:
  1-across at row 1, col 1 (length 2)
  3-across at row 1, col 4 (length 2)
  10-across at row 5, col 1 (length 2)
  11-across at row 5, col 4 (length 2)
  1-down at row 1, col 1 (length 2)
  4-down at row 1, col 5 (length 2)
  8-down at row 4, col 1 (length 2)
  9-down at row 4, col 5 (length 2)
```

## JSON output

Add `--json` to get the same report as structured data instead:

```
$ go run . --json testdata/sample.txt
{
  "width": 5,
  "height": 5,
  "black_squares": 4,
  "white_squares": 21,
  "symmetric": true,
  "min_word_length": 3,
  "words": [
    { "number": 1, "direction": "across", "row": 0, "col": 0, "length": 2 },
    { "number": 1, "direction": "down", "row": 0, "col": 0, "length": 2 },
    { "number": 2, "direction": "down", "row": 0, "col": 1, "length": 5 },
    ...
  ],
  "short_words": [
    { "number": 1, "direction": "across", "row": 0, "col": 0, "length": 2 },
    ...
  ]
}
```

Rows and columns are 0-indexed in the JSON output, matching normal array
indexing, and 1-indexed in the text output, matching how a solver would
point at a square on paper.

## Flags

```
--json                 print the report as JSON instead of plain text
--min-word-length int  shortest word length that doesn't get flagged (default 3)
```

## Building

Requires Go 1.22 or later. No third-party dependencies.

```
go build .
```

## Status

Early. See the issue tracker for what's planned next - letter fill,
unchecked-square detection, and reading grids from stdin are the near-term
targets.

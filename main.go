package main

import (
	"fmt"
	"math"
	"os"

	"github.com/mattn/go-runewidth"
)

type ScoreFuncs interface {
	MatchOrMissMatch(a, b rune) float64
	Gap() float64
}

var DefaultScoreFunc = ScoreFuncs(defaultScoreFunc{})

type defaultScoreFunc struct{}

func (_ defaultScoreFunc) MatchOrMissMatch(a, b rune) float64 {
	if a == b {
		return 1.0
	} else {
		return -0.5
	}
}

func (_ defaultScoreFunc) Gap() float64 {
	return -0.5
}

type IntPair struct {
	X, Y int
}

type SmithWatermanResult struct {
	Scores  [][]ResultCell
	BestPos IntPair
	Input1  []rune
	Input2  []rune
}

type ResultCell struct {
	Score float64
	prev  []IntPair
}

func SmithWaterman(sf ScoreFuncs, r1, r2 []rune) SmithWatermanResult {
	scores := make([][]ResultCell, len(r1)+1)
	for i := range scores {
		scores[i] = make([]ResultCell, len(r2)+1)
	}

	bestPos := IntPair{0, 0}
	var bestScore float64
	for i, ri := range r1 {
		for j, rj := range r2 {
			mmm := scores[i][j].Score + sf.MatchOrMissMatch(ri, rj)
			skip1 := scores[i][j+1].Score + sf.Gap()
			skip2 := scores[i+1][j].Score + sf.Gap()
			nextScore := math.Max(0, math.Max(mmm, math.Max(skip1, skip2)))
			if nextScore > bestScore {
				bestPos.X = i + 1
				bestPos.Y = j + 1
				bestScore = nextScore
			}
			nextCell := &scores[i+1][j+1]
			nextCell.Score = nextScore
			if nextScore == 0 {
				continue
			}
			if nextScore == mmm {
				nextCell.prev = append(nextCell.prev, IntPair{-1, -1})
			}
			if nextScore == skip1 {
				nextCell.prev = append(nextCell.prev, IntPair{-1, 0})
			}
			if nextScore == skip2 {
				nextCell.prev = append(nextCell.prev, IntPair{0, -1})
			}
		}
	}
	return SmithWatermanResult{scores, bestPos, r1, r2}
}

func SmithWatermanString(sf ScoreFuncs, s1, s2 string) SmithWatermanResult {
	return SmithWaterman(sf, []rune(s1), []rune(s2))
}

func reverse(r []rune) {
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
}

const (
	fullHyphen       = rune(0xFF0D)
	fullSpace        = rune(0x3000)
	fullVerticalLine = rune(0xFF5C)
)

func backTrack(r SmithWatermanResult) ([]rune, []rune, []rune) {
	x := r.BestPos.X
	y := r.BestPos.Y
	var xr, yr []rune
	var mr []rune
	for {
		if len(r.Scores[x][y].prev) == 0 {
			break
		}
		p := r.Scores[x][y].prev[0]
		if p.X == 0 {
			if runewidth.RuneWidth(r.Input2[y-1]) == 1 {
				xr = append(xr, rune('-'))
				mr = append(mr, rune(' '))
			} else {
				xr = append(xr, fullHyphen)
				mr = append(mr, fullSpace)
			}
			yr = append(yr, r.Input2[y-1])
		} else if p.Y == 0 {
			if runewidth.RuneWidth(r.Input1[x-1]) == 1 {
				yr = append(yr, rune('-'))
				mr = append(mr, rune(' '))
			} else {
				yr = append(yr, fullHyphen)
				mr = append(mr, fullSpace)
			}
			xr = append(xr, r.Input1[x-1])
		} else {
			if runewidth.RuneWidth(r.Input1[x-1]) == 1 && runewidth.RuneWidth(r.Input2[y-1]) != 1 {
				xr = append(xr, rune(' '))
			}
			if runewidth.RuneWidth(r.Input1[x-1]) != 1 && runewidth.RuneWidth(r.Input2[y-1]) == 1 {
				yr = append(yr, rune(' '))
			}
			xr = append(xr, r.Input1[x-1])
			yr = append(yr, r.Input2[y-1])

			full := runewidth.RuneWidth(r.Input1[x-1]) != 1 || runewidth.RuneWidth(r.Input2[y-1]) != 1
			if r.Input1[x-1] == r.Input2[y-1] {
				if full {
					mr = append(mr, fullVerticalLine)
				} else {
					mr = append(mr, rune('|'))
				}
			} else {
				if full {
					mr = append(mr, fullSpace)
				} else {
					mr = append(mr, rune(' '))
				}
			}
		}
		x += p.X
		y += p.Y
	}
	reverse(xr)
	reverse(yr)
	reverse(mr)
	return xr, yr, mr
}

func (r SmithWatermanResult) MaxScore() float64 {
	p := r.BestPos
	return r.Scores[p.X][p.Y].Score
}

func _main() int {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: ./smith-waterman string1 string2")
		return 1
	}
	res := SmithWatermanString(DefaultScoreFunc, os.Args[1], os.Args[2])
	a, b, m := backTrack(res)
	fmt.Println(string(a))
	fmt.Println(string(m))
	fmt.Println(string(b))
	return 0
}

func main() {
	os.Exit(_main())
}

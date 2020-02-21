package datamatcher

import (
	"testing"

	"github.com/gnames/bhlnames/refs"
)

func TestYearNear(t *testing.T) {
	years := [][]int{{2001, 2000}, {2000, 2001}, {2000, 2000}, {2000, 2002}, {-1, -1}, {3000, 3001}}
	scores := []float32{0.5, 0.5, 1, 0, 0, 0}
	for i, v := range years {
		score := YearNear(v[0], v[1])
		if score != scores[i] {
			t.Errorf("Wrong score for YearNear(%d, %d): %f", v[0], v[1], score)
		}
	}

}

func TestYearBetween(t *testing.T) {
	type data struct {
		values []int
		score  float32
	}

	dataArray := []data{
		{[]int{0, 0, 0}, 0},
		{[]int{0, 2000, 2001}, 0},
		{[]int{0, 2000, 0}, 0},
		{[]int{0, 0, 2000}, 0},
		{[]int{2000, 0, 0}, 1},
		{[]int{2000, 2000, 2001}, 1},
		{[]int{1999, 2000, 2001}, 0},
		{[]int{1998, 2000, 2001}, 0},
		{[]int{2002, 2000, 2001}, 0},
		{[]int{2003, 2000, 2001}, 0},
		{[]int{2001, 2001, 0}, 1},
		{[]int{2001, 2002, 0}, 0.5},
		{[]int{2001, 2003, 0}, 0},
		{[]int{2003, 2002, 0}, 0.5},
		{[]int{3000, 3000, 3000}, 0},
		{[]int{0, 3000, 3000}, 0},
		{[]int{3000, 0, 0}, 0},
		{[]int{0, 0, 3000}, 0},
		{[]int{0, 3000, 0}, 0},
	}

	for _, d := range dataArray {
		score := YearBetween(d.values[0], d.values[1], d.values[2])
		if score != d.score {
			t.Errorf("Wrong score for YearsBetween(%d, %d, %d): %f %f",
				d.values[0], d.values[1], d.values[2], d.score, score)
		}
	}
}

func TestYearScore(t *testing.T) {

	type data struct {
		refType  string
		refYears []int
		year     int
		score    float32
	}

	dataArray := []data{
		// 0 YearAggr
		// 1 ItemYearStart
		// 2 ItemYearEnd
		// 3 TitleYearStart
		// 4 TitleYearEnd
		{"Part", []int{0, 0, 0, 0, 0}, 0, 0},
		{"Part", []int{0, 2000, 2001, 0, 0}, 0, 0},
		{"Part", []int{2000, 2000, 2000, 0, 0}, 2000, 1},
		{"Part", []int{3000, 2000, 3000, 0, 0}, 3000, 0},
		{"Part", []int{3000, 3000, 2001, 0, 0}, 3000, 0},
		{"Title", []int{2000, 0, 0, 1990, 1991}, 2000, 0},
		{"Title", []int{2000, 0, 0, 2000, 0}, 2000, 1},
		{"Part", []int{0, 2000, 2001, 0, 0}, 0, 0},
		{"Part", []int{0, 2000, 2001, 0, 0}, 0, 0},
		{"Part", []int{0, 2000, 2001, 0, 0}, 0, 0},
	}
	_ = dataArray

	for _, d := range dataArray {
		testRef := refs.Reference{
			YearType:       d.refType,
			YearAggr:       d.refYears[0],
			ItemYearStart:  d.refYears[1],
			ItemYearEnd:    d.refYears[2],
			TitleYearStart: d.refYears[3],
			TitleYearEnd:   d.refYears[4],
		}

		result := YearScore(d.year, &testRef)

		if result != d.score {
			t.Errorf("Wrong score for YearScore(%d, %#v) %f %f", d.year, testRef, result, d.score)
		}
	}

}

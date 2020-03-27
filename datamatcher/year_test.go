package datamatcher

import (
	"math"
	"testing"

	"github.com/gnames/bhlnames/refs"
)

func TestYearNear(t *testing.T) {
	years := [][]int{{2001, 2000}, {2000, 2001}, {2000, 2000}, {2000, 2002}, {2000, 2003}, {-1, -1}, {3000, 3001}}
	scores := []float32{0.7, 0.7, 1, 0.49, 0.343, 0, 0}
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
		{[]int{2000, 0, 0}, 0.01},
		{[]int{2000, 2000, 2001}, 0.7},
		{[]int{1999, 2000, 2001}, 0},
		{[]int{2002, 2000, 2001}, 0},
		{[]int{2001, 2001, 0}, 1},
		{[]int{2001, 2002, 0}, 0.7},
		{[]int{2001, 2003, 0}, 0.49},
		{[]int{2003, 2002, 0}, 0.7},
		{[]int{2003, 2003, 2003}, 1},
		{[]int{2002, 1993, 2003}, 0.7},
		{[]int{1993, 1993, 2003}, 0.028248},
		{[]int{1981, 1980, 2003}, 0},
		{[]int{3000, 3000, 3000}, 0},
		{[]int{0, 3000, 3000}, 0},
		{[]int{3000, 0, 0}, 0},
		{[]int{0, 0, 3000}, 0},
		{[]int{0, 3000, 0}, 0},
	}

	for _, d := range dataArray {
		score := YearBetween(d.values[0], d.values[1], d.values[2])
		if math.Abs(float64(score)-float64(d.score)) > 0.0001 {
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
		{"Part", []int{0, 2000, 2001, 0, 0}, 3000, 0},
		{"Part", []int{3000, 3000, 2001, 0, 0}, 3000, 0},
		{"Part", []int{2000, 2000, 2000, 0, 0}, 2000, 1},
		{"Part", []int{2000, 0, 0, 1990, 2001}, 2000, 1},
		{"Title", []int{2000, 0, 0, 2000, 0}, 2000, 1},
		{"Part", []int{0, 2000, 2001, 0, 0}, 2000, 0.7},
		{"Title", []int{1837, 0, 0, 1837, 1858}, 1849, 0.040354},
		{"Title", []int{1837, 0, 0, 1837, 1858}, 1838, 0},
		{"Item", []int{1837, 1837, 1858, 0, 0}, 1849, 0.040354},
		{"Item", []int{1837, 1837, 1858, 0, 0}, 1838, 0},
		{"Item", []int{1837, 1837, 1839, 1837, 1890}, 1838, 0.7},
		{"Item", []int{1837, 1837, 1849, 1837, 1838}, 1838, 1},
		{"Title", []int{0, 0, 0, 0, 0}, 1849, 0.01},
	}

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

		if math.Abs(float64(result)-float64(d.score)) > 0.0001 {
			t.Errorf("Wrong score for YearScore(%d, %#v) %f %f", d.year, testRef, result, d.score)
		}
	}
}

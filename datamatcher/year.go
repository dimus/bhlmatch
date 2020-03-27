package datamatcher

import (
	"math"
	"time"

	"github.com/gnames/bhlnames/refs"
)

func InvalidYear(year int) bool {
	return year < 1740 || year > (time.Now().Year()+2)
}

func YearNear(year1, year2 int) float32 {
	if InvalidYear(year1) {
		return 0
	}
	coef := 0.7
	dif := math.Abs(float64(year1) - float64(year2))
	if dif > 10 {
		return 0
	}
	return float32(math.Pow(float64(coef), dif))
}

func YearBetween(year, yearMin, yearMax int) float32 {
	if InvalidYear(year) {
		return 0
	}
	if yearMin == 0 && yearMax == 0 {
		return 0.01
	}

	if yearMax < yearMin && yearMax != 0 {
		return 0.01
	}

	if yearMax == 0 {
		return YearNear(year, yearMin)
	}

	if !(year <= yearMax && year >= yearMin) {
		return 0
	}

	return YearNear(year, yearMax)
}

func YearScore(year int, ref *refs.Reference) float32 {
	var score float32 = 1
	YearPart, ItemYearStart, ItemYearEnd, TitleYearStart, TitleYearEnd := getRefYears(ref)

	if YearPart > 0 {
		return score * YearNear(year, YearPart)
	}
	var score1, score2 float32
	item := int(ItemYearStart+ItemYearEnd) > 0
	title := int(TitleYearStart+TitleYearEnd) > 0
	if item || (!item && !title) {
		score1 = YearBetween(year, ItemYearStart, ItemYearEnd)
	}
	if title {
		score2 = YearBetween(year, TitleYearStart, TitleYearEnd)
	}

	if score1 > score2 {
		return score1
	}
	return score2
}

func getRefYears(ref *refs.Reference) (int, int, int, int, int) {
	var yearPart int
	if ref.YearType == "Part" {
		yearPart = ref.YearAggr
	}
	return yearPart, ref.ItemYearStart, ref.ItemYearEnd, ref.TitleYearStart, ref.TitleYearEnd
}

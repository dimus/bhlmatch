package datamatcher

import (
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
	switch year1 - year2 {
	case 0:
		return 1
	case 1, -1:
		return 0.5
	default:
		return 0
	}
}

func YearBetween(year, yearMin, yearMax int) float32 {
	if InvalidYear(year) {
		return 0
	}
	if yearMin == 0 && yearMax == 0 {
		return 1
	}

	if yearMax < yearMin && yearMax != 0 {
		return 1
	}

	if year <= yearMax {
		if year >= yearMin {
			return 1
		} else {
			return 0
		}
	}

	if year == yearMin {
		return 1
	}

	if yearMax == 0 {
		return YearNear(year, yearMin)
	}

	if year > yearMin && year <= yearMax {
		return 1
	}
	return 0
}

func YearScore(year int, ref *refs.Reference) float32 {
	var score float32 = 1
	YearPart, ItemYearStart, ItemYearEnd, TitleYearStart, TitleYearEnd := getRefYears(ref)

	if YearPart > 0 {
		return score * YearNear(year, YearPart)
	}

	score *= YearBetween(year, ItemYearStart, ItemYearEnd)
	score *= YearBetween(year, TitleYearStart, TitleYearEnd)
	return score
}

func getRefYears(ref *refs.Reference) (int, int, int, int, int) {
	var yearPart int
	if ref.YearType == "Part" {
		yearPart = ref.YearAggr
	}
	return yearPart, ref.ItemYearStart, ref.ItemYearEnd, ref.TitleYearStart, ref.TitleYearEnd
}

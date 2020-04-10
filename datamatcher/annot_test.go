package datamatcher

import (
	"math"
	"testing"

	"github.com/gnames/bhlnames/refs"
)

func TestAnnotPresent(t *testing.T) {
	if 1 != 1 {
		t.Error("Something is wrong with your 1!")
	}
}

func TestAnnotScore(t *testing.T) {

	type data struct {
		name       string
		matchName  string
		annotation Annot
		score      float32
	}

	dataArray := []data{
		//  Name
		// 1 MatchName
		// 2 Annotation
		// 3
		{"Aus bus", "Aus bus", SpNov, 1.0},
		{"Aus bus cus", "Aus bus cus", SpNov, 0.0},
		{"Aus bus cus", "Aus bus cus", SubsNov, 1.0},
		{"Aus bus", "Aus bus cus", SpNov, 0.1},
		{"Aus bus cus", "Aus bus", SpNov, 0.6},
		{"Aus bus Ower", "Bus cus Mozzherin", SpNov, 1.0},
		{"Aus (Bus) cus", "Aus cus", SpNov, 1.0},
		{"Aus bus", "Aus bus", SubsNov, 0},
		{"Aus bus", "Aus bus cus", SubsNov, 0.4},
		{"Aus bus cus", "Aus bus", SubsNov, 0.1},
		{"Aus bus", "Aus bus", CombNov, 1.0},
		{"Aus bus cus", "Aus bus cus", CombNov, 1.0},
		{"Aus bus", "Aus bus cus", CombNov, 0.6},
		{"Aus", "Aus", CombNov, 0},
		{"Aus bus cus", "Aus", NoAnnot, 1.0},
		{"Aus bus", "Aus bus", NoAnnot, 1.0},
		{"Aus virus", "Bus cus", SpNov, 0.0},
	}

	for _, d := range dataArray {
		testRef := refs.Reference{
			Name:       d.name,
			MatchName:  d.matchName,
			AnnotNomen: d.annotation.String(),
		}
		result := AnnotScore(&testRef)
		if math.Abs(float64(result)-float64(d.score)) > 0.0001 {
			t.Errorf("Wrong score for AnnotScore(%#v) %f, %f", testRef, result, d.score)
		}

	}
}

package datamatcher

import (
	"github.com/gnames/bhlnames/refs"
	"gitlab.com/gogna/gnparser"
)

type Annot int

const (
	NoAnnot Annot = iota
	SpNov
	SubsNov
	CombNov
)

func NewAnnot(annot string) Annot {
	annotations := map[string]Annot{
		"NO_ANNOT":  NoAnnot,
		"SP_NOV":    SpNov,
		"SUBSP_NOV": SubsNov,
		"COMB_NOV":  CombNov,
	}
	if a, ok := annotations[annot]; ok {
		return a
	}
	return NoAnnot
}

func (a Annot) String() string {
	switch int(a) {
	case 1:
		return "SP_NOV"
	case 2:
		return "SUBSP_NOV"
	case 3:
		return "COMB_NOV"
	default:
		return "NO_ANNOT"
	}

}

// NO_ANNOT = 1
// SP_NOV
// f:sp v:sp = 1
// f:sp v:gen = 0
// f:sp v:ssp = 0.1
// f:ssp v:sp = 0.6
// f:ssp v:gen = 0
// f:gen v:gen = 0
// SUBSP_NOV
// f:ssp v:ssp = 1
// f:ssp v:sp = 0
// f:ssp v:gen = 0
// f:sp v:ssp = 0.6
// f:gen v:gen = 0
// COMB_NOV
// f:sp v:sp = 1
// f:ssp v:ssp = 1
// f:ssp v:sp = 0.4
// f:sp v:ssp = 0.1
// f:gen v:gen = 0

func cardinality(ref *refs.Reference) (int32, int32) {
	gnp := gnparser.NewGNparser()
	n := gnp.ParseToObject(ref.Name)
	mn := gnp.ParseToObject(ref.MatchName)
	return n.Cardinality, mn.Cardinality
}

func calculateScore(ref *refs.Reference) float32 {
	annot := NewAnnot(ref.AnnotNomen)
	cardName, cardMatchName := cardinality(ref)
	if cardName == 0 || cardMatchName == 0 {
		return 0.0
	}
	switch annot {
	case SpNov:
		switch {
		case cardName == 2 && cardMatchName == 2:
			return 1.0
		case cardName == 2 && cardMatchName == 3:
			return 0.1
		case cardName == 3 && cardMatchName == 2:
			return 0.6
		default:
			return 0
		}
	case SubsNov:
		switch {
		case cardName == 3 && cardMatchName == 3:
			return 1.0
		case cardName == 3 && cardMatchName == 2:
			return 0.1
		case cardName == 2 && cardMatchName == 3:
			return 0.4
		default:
			return 0
		}
	case CombNov:
		switch {
		case cardName == 2 && cardMatchName == 2:
			return 1.0
		case cardName == 3 && cardMatchName == 3:
			return 1.0
		case cardName == 2 && cardMatchName == 3:
			return 0.6
		}
	case NoAnnot:
		return 1.0
	}
	return 0
}

func AnnotScore(ref *refs.Reference) float32 {
	annot := NewAnnot(ref.AnnotNomen)
	switch annot {
	case NoAnnot:
		return 1
	default:
		return calculateScore(ref)
	}
}

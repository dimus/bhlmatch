package bhlmatch

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/gdower/bhlmatch/datamatcher"
	"github.com/gnames/bhlnames"
	rfs "github.com/gnames/bhlnames/refs"
)

const (
	nameCodeF   = 1
	genusF      = 3
	subGenusF   = 4
	speciesF    = 5
	infraSpF    = 7
	rankMarkerF = 8
	authF       = 27
	journalF    = 29
	pageF       = 31
	volF        = 32
	yearF       = 33
)

type COLref struct {
	nameCode string
	author   string
	year     string
	pages    string
	vol      string
	journal  string
}

func (bhlm BHLmatch) Run() {
	csvFile, err := os.Open(bhlm.InputFile)
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	nameRef := prepareData(data)

	// Setup stream
	in := make(chan string)
	out := make(chan *rfs.RefsResult)
	var wg sync.WaitGroup
	wg.Add(1)

	bn := bhlm.BHLnames
	go bhlnames.RefsStream(bn, in, out)
	go bhlm.processResults(nameRef, bn.Format, out, &wg)
	for k := range nameRef {
		in <- k
	}
	close(in)

}

func matchYear(refYear string, refs []*rfs.Reference) (*rfs.Reference, float32) {
	var bestRef *rfs.Reference
	var bestScore float32
	var score float32
	bestScore = 0.0
	for _, r := range refs {
		yr, err := strconv.Atoi(refYear)
		if err != nil {
			fmt.Printf("Weird year: %s", refYear)
			continue
		}
		score = datamatcher.YearScore(yr, r)
		if score > bestScore {
			bestRef = r
			bestScore = score
		}
	}
	return bestRef, bestScore
}

func matchAnnot(refYear string, refs []*rfs.Reference) (*rfs.Reference, float32) {
	var bestRef *rfs.Reference
	var bestScore float32
	var score float32
	bestScore = 0.0
	for _, r := range refs {
		score = datamatcher.AnnotScore(r)
		if score > bestScore {
			bestRef = r
			bestScore = score
		}
	}
	if bestScore > 0 {
		yr, err := strconv.Atoi(refYear)
		if err == nil {
			yearScore := datamatcher.YearScore(yr, bestRef)
			bestScore += yearScore
		}
	}
	return bestRef, bestScore
}

func (bhlm BHLmatch) processResults(namRef map[string]COLref, format string,
	out <-chan *rfs.RefsResult, wg *sync.WaitGroup) {

	csvOutputFile, err := os.Create(bhlm.OutputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer csvOutputFile.Close()

	csvWriter := csv.NewWriter(csvOutputFile)
	defer csvWriter.Flush()

	defer wg.Done()
	for r := range out {
		if r.Error != nil {
			log.Println(r.Error)
		}
		// fmt.Println("out", r.Output.NameString, len(r.Output.References))
		if col, ok := namRef[r.Output.NameString]; ok {
			refYear, scoreYear := matchYear(col.year, r.Output.References)
			refAnnot, scoreAnnot := matchAnnot(col.year, r.Output.References)
			if scoreYear+scoreAnnot == 0 {
				continue
			}
			bestRef := refYear
			if scoreAnnot > 0 && bestRef.PageID != refAnnot.PageID {
				if (scoreYear - scoreAnnot) < 0 {
					bestRef = refAnnot
				}
			}

			fmt.Printf("Best match: %v, %f, %s\n\n", bestRef, scoreAnnot+scoreYear, col.year)
			output := []string{
				col.nameCode,
				col.author,
				col.journal,
				col.vol,
				col.pages,
				col.year,
				bestRef.Name,
				bestRef.MatchName,
				strconv.Itoa(bestRef.EditDistance),
				bestRef.AnnotNomen,
				bestRef.URL,
			}
			err := csvWriter.Write(output)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("COULD NOT FIND REFERENCE")
		}

		//  r.Output.References[0].URL
		// fmt.Println(bhlnames.FormatOutput(r.Output, format))
	}
}

func prepareData(data [][]string) map[string]COLref {
	refs := make(map[string]COLref)
	var scientificName string
	for i, r := range data {
		if i == 0 {
			continue
		}

		if r[subGenusF] != "\\N" { // names with subgenera
			scientificName = fmt.Sprintf("%s (%s) %s", r[genusF], r[subGenusF], r[speciesF])
		} else { // names without subgenera
			scientificName = fmt.Sprintf("%s %s", r[genusF], r[speciesF])
		}

		if r[infraSpF] != "\\N" { // names with infraspecies
			scientificName = fmt.Sprintf("%s %s %s", scientificName, r[rankMarkerF], r[infraSpF])
		}

		ref := COLref{
			nameCode: r[nameCodeF],
			author:   r[authF],
			journal:  r[journalF],
			vol:      r[volF],
			pages:    r[pageF],
			year:     r[yearF],
		}
		log.Println(scientificName)
		refs[scientificName] = ref
	}
	return refs
}

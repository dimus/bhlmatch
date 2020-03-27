package bhlmatch

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gnames/bhlnames"
	rfs "github.com/gnames/bhlnames/refs"
)

const (
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
	author  string
	year    string
	pages   string
	vol     string
	journal string
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

	// file, err := os.Create(bhlm.OutputFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// csvWriter := csv.NewWriter(file)
	// _ = csvWriter
	//
	// Setup stream
	in := make(chan string)
	out := make(chan *rfs.RefsResult)
	var wg sync.WaitGroup
	wg.Add(1)

	bn := bhlm.BHLnames
	go bhlnames.RefsStream(bn, in, out)
	go processResults(nameRef, bn.Format, out, &wg)
	for k := range nameRef {
		in <- k
	}
	close(in)
}

func processResults(namRef map[string]COLref, format string,
	out <-chan *rfs.RefsResult, wg *sync.WaitGroup) {
	fmt.Println("PROC RES")
	defer wg.Done()
	for r := range out {
		if r.Error != nil {
			log.Println(r.Error)
		}
		fmt.Println("out", r.Output.NameString, len(r.Output.References))
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
			author:  r[authF],
			journal: r[journalF],
			vol:     r[volF],
			pages:   r[pageF],
			year:    r[yearF],
		}
		log.Println(scientificName)
		refs[scientificName] = ref
	}
	return refs
}

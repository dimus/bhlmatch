package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/gnames/bhlnames"
	rfs "github.com/gnames/bhlnames/refs"
	"log"
	"os"
	"sync"
)

func processResults(format string, out <-chan *rfs.RefsResult,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for r := range out {
		if r.Error != nil {
			log.Println(r.Error)
		}
		//  r.Output.References[0].URL
		fmt.Println(bhlnames.FormatOutput(r.Output, format))
	}
}

func main() {
	bhln := bhlnames.BHLnames{}
	bhln.Name = "bhlnames"
	bhln.Host = "localhost"
	bhln.User = "postgres"
	bhln.Pass = "postgres"
	bhln.BHLindexHost = "bhlrpc.globalnames.org:80"
	bhln.InputDir = "/home/gdo/bhlnames"
	bhln.DumpURL = "https://www.biodiversitylibrary.org/data/data.zip"
	bhln.JobsNum = 12
	bhln.MetaData.Configure(bhln.DbOpts)


	csvFile, err := os.Open("tmp/scientific_names_mod3.csv")
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("tmp/scientific_names_with_bhl_urls.csv")
	if err != nil {
		log.Fatal(err)
	}

	csvWriter := csv.NewWriter(file)


	// Setup stream
	in := make(chan string)
	out := make(chan *rfs.RefsResult)
	var wg sync.WaitGroup
	wg.Add(1)

	log.Println(bhln)
	go bhlnames.RefsStream(bhln, in, out)
	//go processResults(bhln.Format, out, &wg)

	var scientificName string
	for i, r := range data {

		if i != 0 { // skip the header

			if r[4] != "\\N" { // names with subgenera
				scientificName = fmt.Sprintf("%s (%s) %s", r[3], r[4], r[5])
			} else { // names without subgenera
				scientificName = fmt.Sprintf("%s %s", r[3], r[5])
			}

			if r[7] != "\\N" { // names with infraspecies
				scientificName = fmt.Sprintf("%s %s %s", scientificName, r[8], r[7])
			}
			log.Println(scientificName)

			bhlurl := ""
			in <- scientificName
			results := <-out
			if len(results.Output.References) > 0 {

				// get the oldest reference with a non-zero year aggregation
				for _, result := range results.Output.References {
					if result.YearAggr != 0 {
						bhlurl = result.URL
						break
					}
				}
			}

			if bhlurl != "" {
				r[2] = bhlurl // replace GSD URL with BHL URL and write output
				err = csvWriter.Write(r)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

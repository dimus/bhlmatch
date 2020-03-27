package bhlmatch

import (
	"testing"
)

func TestBHLMatch(t *testing.T) {
	bhlm := NewBHLmatch()
	if bhlm.InputFile != "" {
		t.Errorf("InputFile should be empty: %s", bhlm.InputFile)
	}

	if bhlm.BHLnames.Host != "" {
		t.Errorf("Host should be empty: %s", bhlm.BHLnames.Host)
	}

}

func makeOpts() []Option {
	opts := []Option{
		OptDbHost("localhost"),
		OptDbUser("postgres"),
		OptDbPass("postgres"),
		OptDbName("bhlnames"),
		OptInputFile("test.csv"),
	}
	return opts
}

func TestCustomOptions(t *testing.T) {
	bhlm := NewBHLmatch(makeOpts()...)
	if bhlm.BHLnames.Host != "localhost" {
		t.Errorf("Incorrect host: %s", bhlm.BHLnames.Host)
	}
	if bhlm.BHLnames.User != "postgres" {
		t.Errorf("Incorrect user: %s", bhlm.BHLnames.User)
	}
	if bhlm.BHLnames.Pass != "postgres" {
		t.Errorf("Incorrect pass: %s", bhlm.BHLnames.Pass)
	}
	if bhlm.BHLnames.Name != "bhlnames" {
		t.Errorf("Incorrect database name: %s", bhlm.BHLnames.Name)
	}
	if bhlm.InputFile != "test.csv" {
		t.Errorf("Incorrect InputFile: %s", bhlm.InputFile)
	}
	if bhlm.BHLnames.JobsNum != 0 {
		t.Errorf("Incorrect jobs number: %d", bhlm.BHLnames.JobsNum)
	}
}

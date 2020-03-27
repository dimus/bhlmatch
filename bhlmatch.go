package bhlmatch

import (
	"github.com/gnames/bhlnames"
)

type BHLmatch struct {
	InputFile  string
	OutputFile string
	BHLnames   bhlnames.BHLnames
}

// Option type for changing GNfinder settings.
type Option func(*BHLmatch)

func OptDbHost(d string) Option {
	return func(bhlm *BHLmatch) {
		bhlm.BHLnames.DbOpts.Host = d
	}
}

func OptDbUser(u string) Option {
	return func(bhlm *BHLmatch) {
		bhlm.BHLnames.DbOpts.User = u
	}
}

func OptDbPass(p string) Option {
	return func(bhlm *BHLmatch) {
		bhlm.BHLnames.DbOpts.Pass = p
	}
}

func OptDbName(n string) Option {
	return func(bhlm *BHLmatch) {
		bhlm.BHLnames.DbOpts.Name = n
	}
}

func OptInputFile(i string) Option {
	return func(bhlm *BHLmatch) {
		bhlm.InputFile = i
	}
}

func OptOutputFile(o string) Option {
	return func(bhlm *BHLmatch) {
		bhlm.OutputFile = o
	}
}

func OptBHLnamesDir(d string) Option {
	return func(bhlm *BHLmatch) {
		bhlm.BHLnames.MetaData.InputDir = d
	}
}

func OptJobsNum(n int) Option {
	return func(bhlm *BHLmatch) {
		bhlm.BHLnames.JobsNum = n
	}
}

func OptTaxonomicMatch(t bool) Option {
	return func(bhlm *BHLmatch) {
		bhlm.BHLnames.NoSynonyms = t
	}
}

func NewBHLmatch(opts ...Option) BHLmatch {
	bhlm := BHLmatch{}
	for _, opt := range opts {
		opt(&bhlm)
	}
	bn := &bhlm.BHLnames
	bn.MetaData.Configure(bn.DbOpts)
	return bhlm
}

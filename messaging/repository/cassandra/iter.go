package cassandra

import "github.com/gocql/gocql"

type IterInterface interface {
	Scanner() ScannerInterface
}

type Iter struct {
	iter *gocql.Iter
}


func NewIter(iter *gocql.Iter) IterInterface {
	return &Iter{
		iter,
	}
}

func (i *Iter) Scanner() ScannerInterface {
	return NewScanner(i.iter.Scanner())
}
